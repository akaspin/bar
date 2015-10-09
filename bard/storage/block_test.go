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

	s := storage.NewBlockStorage(sPath, 2)

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
		"testdata/blobs/26", m.ID))
	assert.NoError(t, err)

	assert.Equal(t, m.String(), m2.String())
}

//func Test_BlockDriver_Exists(t *testing.T) {
//	bn := fixtures.MakeBLOB(t, 10)
//	defer fixtures.KillBLOB(bn)
//
//	// take manifest
//	m, err := fixtures.NewShadowFromFile(bn)
//	assert.NoError(t, err)
//
//	// Try to store file
//	r, err := os.Open(bn)
//	assert.NoError(t, err)
//	defer r.Close()
//
//	s := storage.NewBlockStorage(sPath, 2)
//
//	err = s.WriteBLOB(m.ID, m.Size, r)
//	assert.NoError(t, err)
//	defer cleanup()
//
//	ok, err := s.IsExists(m.ID)
//	assert.NoError(t, err)
//	assert.True(t, ok)
//}
//
//func Test_BlockDriver_NotExists(t *testing.T) {
//
//	s := storage.NewBlockStorage(sPath, 2)
//	defer cleanup()
//
//	ok, err := s.IsExists(
//		"1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253")
//	assert.NoError(t, err)
//	assert.False(t, ok)
//}
//
//func Test_BlockDriver_ReadBLOB(t *testing.T) {
//	bn := fixtures.MakeBLOB(t, 1024 * 1024 * 5 + 456)
//	defer fixtures.KillBLOB(bn)
//
//	// take manifest
//	m, err := fixtures.NewShadowFromFile(bn)
//	assert.NoError(t, err)
//
//	// Try to store file
//	r, err := os.Open(bn)
//	assert.NoError(t, err)
//	defer r.Close()
//
//	s := storage.NewBlockStorage(sPath, 2)
//
//	err = s.WriteBLOB(m.ID, m.Size, r)
//	assert.NoError(t, err)
//	defer cleanup()
//
//	sr, err := s.ReadBLOB(m.ID)
//	assert.NoError(t, err)
//	defer sr.Close()
//
//	m2, err := manifest.NewFromAny(sr, manifest.CHUNK_SIZE)
//	assert.NoError(t, err)
//	assert.Equal(t, m.ID, m2.ID)
//}

func cleanup() {
	os.RemoveAll(sPath)
}
