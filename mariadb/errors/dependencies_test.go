package errors_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestDatabaseErrorsDoNotImportEcho(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "github.com/pedidopago/go-common/mariadb/errors")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("listar dependências: %v\n%s", err, output)
	}

	for dependency := range strings.Lines(string(output)) {
		if strings.HasPrefix(strings.TrimSpace(dependency), "github.com/labstack/echo/") {
			t.Fatalf("erros de banco não devem importar servidor HTTP: %s", dependency)
		}
	}
}
