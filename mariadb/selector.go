package mariadb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pedidopago/go-common/mariadb/errors"
	"github.com/pedidopago/go-common/util"
)

type Selector interface {
	// ApplyToSquirrel applies the selector to the squirrel builder.
	ApplyToSquirrel(b squirrel.SelectBuilder) squirrel.SelectBuilder
	// ValidateFields checks if there is enough data to query the table
	ValidateFields() error
	TableName() string
	SelectColumns() []string
}

func SelectMany[T any](ctx context.Context, db sqlx.ExtContext, input Selector) ([]T, error) {
	db = unsafedb(db) // avoid error on select "*"
	if err := input.ValidateFields(); err != nil {
		return nil, err
	}
	b := squirrel.Select(input.SelectColumns()...).From(input.TableName())
	b = input.ApplyToSquirrel(b)
	if limit, ok := withLimit(ctx); ok {
		b = b.Limit(limit)
	}
	q, args, err := b.ToSql()
	if err != nil {
		return nil, errors.InvalidQuery(err)
	}
	items := make([]T, 0)
	if err := sqlx.SelectContext(ctx, db, &items, q, args...); err != nil {
		return nil, errors.WrapSQLX(err)
	}
	return items, nil
}

func SelectOne[T any](ctx context.Context, db sqlx.ExtContext, input Selector) (T, error) {
	var zv T
	ctx = LimitOne(ctx)
	items, err := SelectMany[T](ctx, db, input)
	if err != nil {
		return zv, err
	}
	if len(items) == 0 {
		return zv, errors.ErrNoResults
	}
	p0 := items[0]
	return p0, nil
}

func unsafedb(db sqlx.ExtContext) sqlx.ExtContext {
	if rdb, ok := db.(*sqlx.DB); ok {
		return rdb.Unsafe()
	}
	return db
}

func withLimit(ctx context.Context) (uint64, bool) {
	if ctx == nil {
		return 0, false
	}
	limit := ctx.Value(CtxLimit)
	if limit == nil {
		return 0, false
	}
	return limit.(uint64), true
}

type ContextVar string

const (
	CtxLimit ContextVar = "model_limit"
)

func LimitOne(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxLimit, uint64(1))
}

func LimitX(ctx context.Context, v uint64) context.Context {
	return context.WithValue(ctx, CtxLimit, v)
}

func SWhereOneOrMany[T comparable](b squirrel.SelectBuilder, col string, vals []T) squirrel.SelectBuilder {
	if len(vals) == 0 {
		return b
	}
	if len(vals) > 1 {
		return b.Where(fmt.Sprintf("%s IN (%s)", col, squirrel.Placeholders(len(vals))), util.ToInterfaces(vals)...)
	}
	return b.Where(fmt.Sprintf("%s = ?", col), vals[0])
}

func SWhereOneOrManyLike(b squirrel.SelectBuilder, col string, vals []string) squirrel.SelectBuilder {
	if len(vals) == 0 {
		return b
	}
	if len(vals) > 1 {
		return b.Where(fmt.Sprintf("%s IN (%s)", col, squirrel.Placeholders(len(vals))), util.ToInterfaces(vals)...)
	}
	if strings.HasPrefix(vals[0], "%") || strings.HasSuffix(vals[0], "%") {
		return b.Where(fmt.Sprintf("%s LIKE ?", col), vals[0])
	}
	return b.Where(fmt.Sprintf("%s = ?", col), vals[0])
}
