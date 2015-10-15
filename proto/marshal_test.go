package proto_test
import (
	"testing"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
)

func Test_Manifest_MarshalThrift(t *testing.T) {
	in := `BAR:MANIFEST

		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234


		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		offset 0
		`
	m, err := proto.NewFromManifest(fixtures.CleanInput(in))
	assert.NoError(t, err)

	_, err = (*m).MarshalThrift()
	assert.NoError(t, err)
}
