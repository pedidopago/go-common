package mariadb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	ID    int64  `db:"id"`
	Age   string `db:"age"`
	Props struct {
		First string `db:"first"`
		Last  string `db:""`
	} `db:"props"`
	Helium *TestHelium `db:"helium"`
}

type TestHelium struct {
	CreatedAt time.Time `db:"created_at"`
}

func TestExtractTest(t *testing.T) {
	elems := ExtractColumnsOfStruct("db", TestStruct{}, WithBackticksColumns("age"))
	require.Equal(t, 5, len(elems))
	assert.Equal(t, "id", elems[0])
	assert.Equal(t, "`age`", elems[1])
	assert.Equal(t, "props.first", elems[2])
	assert.Equal(t, "props.last", elems[3])
	assert.Equal(t, "helium.created_at", elems[4])
}
