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
				}
			]
		}
	}
}`

func TestQuery(t *testing.T) {
	s := Search{}
	BoolFilterTerm(s.Q(), "account_id", 1)
	jb, err := json.MarshalIndent(s, "", "	")
	assert.NoError(t, err)
	assert.Equal(t, qcomparison, string(jb))
}
