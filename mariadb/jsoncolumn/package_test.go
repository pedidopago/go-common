package jsoncolumn

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestAlpha struct {
	Name  string           `json:"name" db:"name"`
	Betas Text[[]TestBeta] `json:"betas" db:"betas"`
}

type TestBeta struct {
	Name  string `json:"name" db:"name"`
	Score int    `json:"score" db:"score"`
}

func TestMarshalUnmarshalJSON(t *testing.T) {
	mike := TestAlpha{
		Name: "Mike",
		Betas: Text[[]TestBeta]{
			Data: &[]TestBeta{},
		},
	}
	*mike.Betas.Data = append(*mike.Betas.Data, TestBeta{
		Name:  "Beta 1",
		Score: 100,
	})
	*mike.Betas.Data = append(*mike.Betas.Data, TestBeta{
		Name:  "Beta 2",
		Score: 200,
	})
	mshed, err := json.MarshalIndent(mike, "", "  ")
	assert.NoError(t, err)

	expected := `{
  "name": "Mike",
  "betas": [
    {
      "name": "Beta 1",
      "score": 100
    },
    {
      "name": "Beta 2",
      "score": 200
    }
  ]
}`

	assert.Equal(t, expected, string(mshed))

	var mike2 TestAlpha
	err = json.Unmarshal(mshed, &mike2)

	assert.NoError(t, err)
	assert.Equal(t, (*mike.Betas.Data)[0].Name, (*mike2.Betas.Data)[0].Name)
}

func TestScanner(t *testing.T) {
	item := TestAlpha{}

	jdata := `[{"name": "Alpha", "score": -1},{"name": "Bravo", "score": -2}]`
	assert.NoError(t, item.Betas.Scan([]byte(jdata)))
	assert.NotNil(t, item.Betas.Data)
	assert.Equal(t, 2, len(*item.Betas.Data))
	assert.NoError(t, item.Betas.Scan(nil))
}
