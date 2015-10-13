package manifest_test
import (
	"testing"
	"github.com/akaspin/bar/manifest"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

//
//func Test_ID_String(t *testing.T) {
//	var fixt manifest.ID
//	fixt = []byte{0x00, 0xFF}
//
//	assert.Equal(t, "00ff", fixt.String())
//}

func Test_ID_MarshallJson(t *testing.T)  {
	var res1 struct{
		ID manifest.ID
	}
	res, err := json.Marshal(struct{
		Id manifest.ID
	}{manifest.ID{0x00, 0xFF}})
	assert.NoError(t, err)
	t.Log(string(res))

	err = json.Unmarshal(res, &res1)
	assert.NoError(t, err)
	t.Log(res1)
}