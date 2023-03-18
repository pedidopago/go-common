package orm

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pedidopago/go-common/mariadb/errors"
	"github.com/pedidopago/go-common/util"
)

var (
	ErrColumnValueMismatch = errors.New("column and value length mismatch")
)

type Model interface {
	TableName() string
}

type ModelInserter interface {
	Model
	InsertList() (keys []string, values []any)
}

type ModelUpdater interface {
	Model
	UpdateList() (keys []string, values []any, whereKeys []string, whereValues []any)
}

type ModelUpserter interface {
	Model
	UpsertList() (insertKeys []string, intertValues []any, updateKeys []string, updateValues []any, whereKeys []string, whereValues []any)
}

type ModelDeleter interface {
	Model
	DeleteList() (whereKeys []string, whereValues []any)
}

type Table[T Model] struct {
	DB        sqlx.ExtContext
	ReplicaDB sqlx.QueryerContext // optional read-only replica database
	SelectTag string              // optional tag to use for select columns (default: db)
}

type TableInserter[T ModelInserter] struct {
	Table[T]
}

type TableUpdater[T ModelUpdater] struct {
	Table[T]
}

func (t Table[T]) readableDB() sqlx.QueryerContext {
	if t.ReplicaDB != nil {
		return t.ReplicaDB
	}
	return t.DB
}

type SelectWhere func() (string, []interface{})

type SelectModifier[T Model] func(sb squirrel.SelectBuilder) squirrel.SelectBuilder

func (t Table[T]) Select(ctx context.Context, where SelectWhere, modifiers ...SelectModifier[T]) ([]T, error) {
	var zv T
	sq := squirrel.Select(ExtractSelectColumnsOfStruct(zv, util.Default(t.SelectTag, "db"))...)
	sq = sq.From(zv.TableName())
	if where != nil {
		whereClause, whereArgs := where()
		sq = sq.Where(whereClause, whereArgs...)
	}
	for _, modifier := range modifiers {
		sq = modifier(sq)
	}
	q, args, err := sq.ToSql()
	if err != nil {
		return nil, errors.InvalidQuery(err)
	}
	results := make([]T, 0)
	if err := sqlx.SelectContext(ctx, t.readableDB(), &results, q, args...); err != nil {
		return nil, errors.WrapSQLX(err)
	}
	return results, nil
}

func (t Table[T]) SelectCount(where SelectWhere, modifiers ...SelectModifier[T]) (int64, error) {
	var zv T
	sq := squirrel.Select("COUNT(*)")
	sq = sq.From(zv.TableName())
	if where != nil {
		whereClause, whereArgs := where()
		sq = sq.Where(whereClause, whereArgs...)
	}
	for _, modifier := range modifiers {
		sq = modifier(sq)
	}
	q, args, err := sq.ToSql()
	if err != nil {
		return 0, errors.InvalidQuery(err)
	}
	var count int64
	if err := sqlx.GetContext(context.Background(), t.readableDB(), &count, q, args...); err != nil {
		return 0, errors.WrapSQLX(err)
	}
	return count, nil
}

func (t TableInserter[T]) Insert(ctx context.Context, model T) error {
	sq := squirrel.Insert(model.TableName())
	cols, vals := model.InsertList()
	sq = sq.Columns(cols...).Values(vals...)
	q, args, err := sq.ToSql()
	if err != nil {
		return errors.InvalidQuery(err)
	}
	if _, err := t.DB.ExecContext(ctx, q, args...); err != nil {
		return errors.WrapSQLX(err)
	}
	return nil
}

func (t TableInserter[T]) InsertAndRetrieveLastID(ctx context.Context, model T) (int64, error) {
	sq := squirrel.Insert(model.TableName())
	cols, vals := model.InsertList()
	sq = sq.Columns(cols...).Values(vals...)
	q, args, err := sq.ToSql()
	if err != nil {
		return 0, errors.InvalidQuery(err)
	}
	res, err := t.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, errors.WrapSQLX(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.WrapSQLX(err)
	}
	return id, nil
}

func (t TableUpdater[T]) Update(ctx context.Context, model T) error {
	sq := squirrel.Update(model.TableName())
	cols, vals, whereCols, whereVals := model.UpdateList()
	if len(cols) != len(vals) || len(whereCols) != len(whereVals) {
		return ErrColumnValueMismatch
	}
	for i := range cols {
		sq = sq.Set(cols[i], vals[i])
	}
	for i := range whereCols {
		sq = sq.Where(squirrel.Eq{whereCols[i]: whereVals[i]})
	}
	q, args, err := sq.ToSql()
	if err != nil {
		return errors.InvalidQuery(err)
	}
	if _, err := t.DB.ExecContext(ctx, q, args...); err != nil {
		return errors.WrapSQLX(err)
	}
	return nil
}
