package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pedidopago/go-common/mariadb/errors"
	"github.com/pedidopago/go-common/slice"
	"golang.org/x/exp/constraints"
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
		if page, ok := withPage(ctx); ok {
			b = b.Offset((page - 1) * limit)
		}
	}
	if oby, ok := withOrderBy(ctx); ok {
		b = b.OrderBy(oby...)
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

func SelectCount[T constraints.Ordered](ctx context.Context, db sqlx.QueryerContext, input Selector) (T, error) {
	var zv T
	if err := input.ValidateFields(); err != nil {
		return zv, err
	}
	b := squirrel.Select("COUNT(*)").From(input.TableName())
	b = input.ApplyToSquirrel(b)
	q, args, err := b.ToSql()
	if err != nil {
		return zv, errors.InvalidQuery(err)
	}
	var count T
	if err := db.QueryRowxContext(ctx, q, args...).Scan(&count); err != nil {
		return zv, errors.WrapSQLX(err)
	}
	return count, nil
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

func withPage(ctx context.Context) (uint64, bool) {
	if ctx == nil {
		return 0, false
	}
	page := ctx.Value(CtxPage)
	if page == nil {
		return 0, false
	}
	return page.(uint64), true
}

func withOrderBy(ctx context.Context) ([]string, bool) {
	if ctx == nil {
		return nil, false
	}
	page := ctx.Value(CtxOrderBy)
	if page == nil {
		return nil, false
	}
	return page.([]string), true
}

type ContextVar string

const (
	CtxLimit   ContextVar = "model_limit"
	CtxPage    ContextVar = "model_page"
	CtxOrderBy ContextVar = "model_order_by"
)

func LimitOne(ctx context.Context) context.Context {
	return context.WithValue(ctx, CtxLimit, uint64(1))
}

func LimitX(ctx context.Context, v uint64) context.Context {
	return context.WithValue(ctx, CtxLimit, v)
}

func PageX(ctx context.Context, v uint64) context.Context {
	return context.WithValue(ctx, CtxPage, v)
}

func OrderByX(ctx context.Context, items ...string) context.Context {
	return context.WithValue(ctx, CtxOrderBy, items)
}

func SWhereOneOrMany[T comparable](b squirrel.SelectBuilder, col string, vals []T) squirrel.SelectBuilder {
	if len(vals) == 0 {
		return b
	}
	if len(vals) > 1 {
		return b.Where(fmt.Sprintf("%s IN (%s)", col, squirrel.Placeholders(len(vals))), slice.ToInterfaces(vals)...)
	}
	return b.Where(fmt.Sprintf("%s = ?", col), vals[0])
}

func SWhereOneOrManyLike(b squirrel.SelectBuilder, col string, vals []string) squirrel.SelectBuilder {
	if len(vals) == 0 {
		return b
	}
	if len(vals) > 1 {
		return b.Where(fmt.Sprintf("%s IN (%s)", col, squirrel.Placeholders(len(vals))), slice.ToInterfaces(vals)...)
	}
	if strings.HasPrefix(vals[0], "%") || strings.HasSuffix(vals[0], "%") {
		return b.Where(fmt.Sprintf("%s LIKE ?", col), vals[0])
	}
	return b.Where(fmt.Sprintf("%s = ?", col), vals[0])
}

func GetWithBuilder(ctx context.Context, dst interface{}, db sqlx.QueryerContext, q squirrel.SelectBuilder) error {
	sq, args, err := q.ToSql()
	if err != nil {
		return errors.InvalidQuery(err)
	}
	return sqlx.GetContext(ctx, db, dst, sq, args...)
}

func SelectWithBuilder(ctx context.Context, dst interface{}, db sqlx.QueryerContext, q squirrel.SelectBuilder) error {
	sq, args, err := q.ToSql()
	if err != nil {
		return errors.InvalidQuery(err)
	}
	return sqlx.SelectContext(ctx, db, dst, sq, args...)
}

func InsertWithBuilder(ctx context.Context, db sqlx.ExecerContext, q squirrel.InsertBuilder) (sql.Result, error) {
	sq, args, err := q.ToSql()
	if err != nil {
		return nil, errors.InvalidQuery(err)
	}
	return db.ExecContext(ctx, sq, args...)
}
