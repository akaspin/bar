package storage_test
import (
	"testing"
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/akaspin/bar/shadow"
)

const sPath  = "test_storage"

func Test_BlockDriver_StoreBLOB(t *testing.T) {
	bn, err := fixtures.MakeBLOB(10)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	// take manifest
	m, err := shadow.NewShadowFromFile(bn, true, shadow.CHUNK_SIZE)
	assert.NoError(t, err)

	// Try to store file
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	s := storage.NewBlockStorage(sPath, 2)

	err = s.StoreBLOB(m.ID, m.Size, r)
	assert.NoError(t, err)
	defer cleanup()

	// check stored file manifest
	m2, err := shadow.NewShadowFromFile(fixtures.StoredName(sPath, m.ID),
		true, shadow.CHUNK_SIZE)
	assert.NoError(t, err)

	assert.Equal(t, m.String(), m2.String())
}

func Test_BlockDriver_Exists(t *testing.T) {
	bn, err := fixtures.MakeBLOB(10)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	// take manifest
	m, err := shadow.NewShadowFromFile(bn, true, shadow.CHUNK_SIZE)
	assert.NoError(t, err)

	// Try to store file
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	s := storage.NewBlockStorage(sPath, 2)

	err = s.StoreBLOB(m.ID, m.Size, r)
	assert.NoError(t, err)
	defer cleanup()

	ok, err := s.IsExists(m.ID)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func Test_BlockDriver_NotExists(t *testing.T) {

	s := storage.NewBlockStorage(sPath, 2)
	defer cleanup()

	ok, err := s.IsExists(
		"1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253")
	assert.NoError(t, err)
	assert.False(t, ok)
}


func cleanup() {
	os.RemoveAll(sPath)
}
