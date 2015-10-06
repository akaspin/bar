package proto_test
import (
	"testing"
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/proto"
	"sort"
	"golang.org/x/crypto/sha3"
	"encoding/hex"
)

func Test_Spec1(t *testing.T) {
	m1, err := manifest.NewFromManifest(fixtures.CleanInput(`BAR:MANIFEST

		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234


		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		offset 0
	`))
	assert.NoError(t, err)
	m2, err := manifest.NewFromManifest(fixtures.CleanInput(`BAR:MANIFEST

		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234


		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		offset 0
	`))
	assert.NoError(t, err)

	spec, err := proto.NewSpec(map[string]string{
		"file/1": m1.ID,
		"file/2": m2.ID,
	})

	// hand-made fixture
	var sorted sort.StringSlice
	sorted = append(sorted, "file/2")
	sorted = append(sorted, "file/1")
	sorted.Sort()

	hasher := sha3.New256()
	id, err := hex.DecodeString(m1.ID)
	assert.NoError(t, err)

	_, err = hasher.Write([]byte(sorted[0]))
	_, err = hasher.Write(id)

	_, err = hasher.Write([]byte(sorted[1]))
	_, err = hasher.Write(id)

	assert.Equal(t, spec.ID, hex.EncodeToString(hasher.Sum(nil)))
}