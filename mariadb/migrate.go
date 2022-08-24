package mariadb

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pedidopago/go-common/discord"
	"github.com/pedidopago/go-common/util"
)

type MigrationCommand string

const (
	MigrateUp    MigrationCommand = "up"
	MigrateNew   MigrationCommand = "new"
	MigrateCheck MigrationCommand = "check"
	MigrateDown  MigrationCommand = "down"
	MigrateForce MigrationCommand = "force"
	MigrateStep  MigrationCommand = "step"
)

type MigrateInput struct {
	DatabaseURL    string
	Command        MigrationCommand
	MigrationsPath string
	Args           []string
	DiscordWebhook string
	ServiceName    string
}

func Migrate(input MigrateInput) error {
	if strings.HasPrefix(input.DatabaseURL, "_ENV_") {
		input.DatabaseURL = os.Getenv(strings.TrimPrefix(input.DatabaseURL, "_ENV_"))
	}

	if !strings.Contains(input.DatabaseURL, "multiStatements") {
		if !strings.Contains(input.DatabaseURL, "?") {
			input.DatabaseURL += "?multiStatements=true"
		} else {
			input.DatabaseURL += "&multiStatements=true"
		}
	}

	db, err := sql.Open("mysql", input.DatabaseURL)

	if err != nil && (input.Command != MigrateNew && input.Command != MigrateCheck && input.Command != "c") {
		return fmt.Errorf("could not open database: %w", err)
	}
	if err == nil {
		defer db.Close()
	}
	var driver database.Driver
	var m *migrate.Migrate
	if db != nil {
		driver, err = mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			discord.NewWebhookMessage(discord.NewWebhookMessageInput{
				Username: util.Default(input.ServiceName, "Some Service"),
				URL:      input.DiscordWebhook,
				Content:  "Could not create the migration driver",
				Logs:     err.Error(),
			})
			return err
		}
		m, err = migrate.NewWithDatabaseInstance(input.MigrationsPath, "mysql", driver)
		if err != nil {
			discord.NewWebhookMessage(discord.NewWebhookMessageInput{
				Username: util.Default(input.ServiceName, "Some Service"),
				URL:      input.DiscordWebhook,
				Content:  "Could not create the migration driver (2)",
				Logs:     err.Error(),
			})
			return err
		}
	}

	switch input.Command {
	case "up", "u":
		err = m.Up()
	case "down", "d":
		err = m.Down()
	case "force", "f":
		if len(input.Args) < 1 {
			return fmt.Errorf("please specify a version to force the migration to")
		}
		var version int
		version, err = strconv.Atoi(input.Args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(3)
		}
		err = m.Force(version)
	case "step", "steps", "s":
		if len(input.Args) < 1 {
			fmt.Fprintln(os.Stderr, "Please specify a step count")
			os.Exit(4)
		}
		var steps int
		steps, err = strconv.Atoi(input.Args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(5)
		}
		err = m.Steps(steps)
	case "sync":
		mainDir := strings.TrimPrefix(input.MigrationsPath, "file://")
		maxv := maxMigrationVersion(mainDir)
		err = m.Migrate(uint(maxv))
	case "new":
		if len(input.Args) < 1 {
			fmt.Fprintln(os.Stderr, "Please specify a migration name")
			os.Exit(6)
		}
		mainDir := strings.TrimPrefix(input.MigrationsPath, "file://")
		maxn := maxMigrationVersion(mainDir)
		updata := "-- write your UP migration here\n"
		downdata := "-- write your DOWN migration here\n"
		fnup := filepath.Join(mainDir, fmt.Sprintf("%05d_%s.up.sql", maxn+1, input.Args[0]))
		fndown := filepath.Join(mainDir, fmt.Sprintf("%05d_%s.down.sql", maxn+1, input.Args[0]))
		if err := os.WriteFile(fnup, []byte(updata), 0644); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(7)
		}
		if err := os.WriteFile(fndown, []byte(downdata), 0644); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(8)
		}
		fmt.Printf("Created migration files:\n%s\n%s\n", fnup, fndown)
		os.Exit(0)
	case "check", "c":
		upfiles := make(map[string]bool)
		downfiles := make(map[string]bool)
		mainDir := strings.TrimPrefix(input.MigrationsPath, "file://")
		filepath.Walk(mainDir, func(path string, info fs.FileInfo, err error) error {
			if strings.HasSuffix(path, ".up.sql") {
				_, fname := filepath.Split(path)
				num := strings.SplitN(fname, "_", 2)[0]
				if upfiles[num] {
					fmt.Fprintf(os.Stderr, "Duplicate migration version %s (%s)\n", num, path)
					os.Exit(9)
				}
				upfiles[num] = true
			}
			if strings.HasSuffix(path, ".down.sql") {
				_, fname := filepath.Split(path)
				num := strings.SplitN(fname, "_", 2)[0]
				if downfiles[num] {
					fmt.Fprintf(os.Stderr, "Duplicate migration version %s (%s)\n", num, path)
					os.Exit(9)
				}
				downfiles[num] = true
			}
			return nil
		})
		fmt.Fprintln(os.Stdout, "All migration files are unique")
		os.Exit(0)
	default:
		var version uint64
		version, err = strconv.ParseUint(string(input.Command), 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(10)
		}
		err = m.Migrate(uint(version))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if err.Error() == "no change" {
			return nil
		}
		_ = discord.NewWebhookMessage(discord.NewWebhookMessageInput{
			Username: util.Default(input.ServiceName, "Some Service"),
			URL:      input.DiscordWebhook,
			Content:  "Could not create the migration driver (2)",
			Logs:     err.Error(),
		})
		return err
	}
	return nil
}

// var (
// 	databaseURL = flag.String("database-url", "", "Database URL")
// 	command     = flag.String("command", "up", "Migration command")
// 	migrations  = flag.String("migrations", "file://database/migrations", "Migrations directory")
// )

// func main() {
// 	flag.Parse()
// 	args := append([]string{""}, flag.Args()...)

func maxMigrationVersion(migrationsPath string) uint {
	var maxn uint = 0
	filepath.Walk(migrationsPath, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sql") {
			_, fname := filepath.Split(path)
			num := strings.SplitN(fname, "_", 2)[0]
			uintNum, err := strconv.ParseUint(num, 10, 64)
			if err == nil {
				if uint(uintNum) > maxn {
					maxn = uint(uintNum)
				}
			}
		}
		return nil
	})
	return maxn
}
