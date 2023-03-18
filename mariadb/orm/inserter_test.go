package orm_test

import (
	"testing"
	"time"

	"github.com/pedidopago/go-common/mariadb/orm"
	"golang.org/x/exp/slices"
)

type Client struct {
	ID        int64            `db:"id" insert:"-"`
	Name      string           `db:"name" insert:"n_a_m_e,omitempty"`
	LastOrder *ClientLastOrder `db:"last_order" insert:",inline"`
	HasBacon  bool
}

type ClientLastOrder struct {
	ID        int64     `db:"id" insert:"last_order_id"`
	OrderDate time.Time `db:"order_date" insert:"last_order_date"`
}

func TestExtractInsertColumnsOfStruct(t *testing.T) {
	c := Client{
		ID:   1,
		Name: "John Doe",
		LastOrder: &ClientLastOrder{
			ID:        2,
			OrderDate: time.Now(),
		},
	}
	keys, values := orm.ExtractInsertColumnsOfStruct(c, "insert")
	t.Logf("keys: %v", keys)
	t.Logf("values: %v", values)
	if !slices.Contains(keys, "n_a_m_e") {
		t.Errorf("expected keys to contain n_a_m_e")
	}
	if !slices.Contains(keys, "hasbacon") {
		t.Errorf("expected keys to contain hasbacon")
	}
}
