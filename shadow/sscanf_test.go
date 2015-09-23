package shadow_test
import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"encoding/hex"
)

func Test_FieldScan(t *testing.T) {
	var str string
	in := "version 0.1.0"
	_, err := fmt.Sscanf(in, "version %s", &str)
	assert.NoError(t, err)
	assert.Equal(t, "0.1.0", str)
}

func Test_FieldScanFail(t *testing.T) {
	str := "test"
	in := "versio 0.1.0"
	_, err := fmt.Sscanf(in, "version %s", &str)
	assert.Error(t, err, "input does not match format")
	assert.Equal(t, "test", str)
}

func Test_FieldScanHex(t *testing.T) {
	in := "id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c"
	var res []byte
	_, err := fmt.Sscanf(in, "id %x", &res)
	assert.NoError(t, err)
	assert.Equal(t, "ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c", hex.EncodeToString(res))
}

func Test_ScanNoSpace(t *testing.T) {
	in := "id..ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c"
	var res []byte
	_, err := fmt.Sscanf(in, "id..%x", &res)
	assert.NoError(t, err)
	assert.Equal(t, "ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c", hex.EncodeToString(res))
}
