package storage_test
import (
	"testing"
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"os"
	"bytes"
	"path/filepath"
)

const sPath  = "testdata"

func newStorage() *storage.BlockStorage {
	return storage.NewBlockStorage(&storage.BlockStorageOptions{
		sPath, 2, 32, 32,
	})
}

func Test_BlockDriver_WriteBLOB(t *testing.T) {

	cleanup()

	bn := fixtures.MakeBLOB(t, 1024 * 1024 * 5 + 435)
	defer fixtures.KillBLOB(bn)

	// take manifest
	m, err := fixtures.NewShadowFromFile(bn)
	assert.NoError(t, err)

	// Try to store file
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	s := newStorage()

	// declare upload
	err = s.DeclareUpload(*m)
	assert.NoError(t, err)

	// Write chunks
	var buf []byte
	for _, c := range m.Chunks {
		buf = make([]byte, c.Size)
		_, err = r.Read(buf)
		assert.NoError(t, err)
		err = s.WriteChunk(m.ID, c.ID, c.Size, bytes.NewReader(buf))
		assert.NoError(t, err)
	}

	// finish upload
	err = s.FinishUpload(m.ID)
	assert.NoError(t, err)

	// check stored file manifest
	m2, err := fixtures.NewShadowFromFile(filepath.Join(
		"testdata/blobs/26", m.ID.String()))
	assert.NoError(t, err)

	assert.Equal(t, m.String(), m2.String())
}

func cleanup() {
	os.RemoveAll(sPath)
}
