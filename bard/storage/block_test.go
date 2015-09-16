package storage_test
import (
	"testing"
	"github.com/akaspin/bar/bard/storage"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/akaspin/bar/shadow"
	"encoding/hex"
)

const sPath  = "test_storage"

func Test_BlockDriver1(t *testing.T) {
	bn, err := fixtures.MakeBLOB(10)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	// take manifest
	r, err := os.Open(bn)
	assert.NoError(t, err)
	m := &shadow.Shadow{}
	assert.NoError(t, m.FromAny(r))
	r.Close()

	// Try to store file
	r, err = os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	s := storage.NewBlockStorage(sPath, 2)

	err = s.StoreBLOB(hex.EncodeToString(m.ID), m.Size, r)
	assert.NoError(t, err)
	defer s.Destroy(hex.EncodeToString(m.ID))
}
