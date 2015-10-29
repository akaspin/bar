package transport_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JSONUnmarshallTypeless(t *testing.T) {
	data := []string{
		"one", "two",
	}
	raw, err := json.Marshal(data)
	assert.NoError(t, err)

	var res []string
	err = json.Unmarshal(raw, &res)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
}
