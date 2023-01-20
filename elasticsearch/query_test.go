package elasticsearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const qcomparison = `{
	"query": {
		"bool": {
			"filter": [
				{
					"term": {
						"account_id": 1
					}
				},
				{
					"term": {
						"store_id": "WABA"
					}
				}
			]
		}
	}
}`

func TestQuery(t *testing.T) {
	s := &Search{}
	BoolFilterTerm(s.Q(), "account_id", 1)
	BoolFilterTerm(s.Q(), "store_id", "WABA")
	jb, err := json.MarshalIndent(s, "", "	")
	assert.NoError(t, err)
	assert.Equal(t, qcomparison, string(jb))
}
