package model_test
import (
	"testing"
	"github.com/akaspin/bar/barc/model"
	"os"
	"path/filepath"
	"github.com/akaspin/bar/manifest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
	"encoding/hex"
	"bytes"
	"io/ioutil"
	"github.com/akaspin/bar/fixtures"
	"github.com/akaspin/bar/barc/lists"
)

func Test_Assembler_StoreChunk(t *testing.T) {
	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "testdata", "assembler-StoreChunk")
	m, err := model.New(wd, false, manifest.CHUNK_SIZE, 128)
	assert.NoError(t, err)

	data := []byte("mama myla ramu")
	hasher := sha3.New256()
	_, err = hasher.Write([]byte(data))
	id := manifest.ID(hex.EncodeToString(hasher.Sum(nil)))

	a, err := model.NewAssembler(m)
	assert.NoError(t, err)
	defer a.Close()

	err = a.StoreChunk(bytes.NewReader(data), id)
	assert.NoError(t, err)

	// check stored chunk
	f, err := os.Open(filepath.Join(a.Where, id.String()))
	assert.NoError(t, err)
	defer f.Close()
	defer os.Remove(filepath.Join(a.Where, id.String()))

	r2, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	assert.Equal(t, data, r2)
}

func Test_Assembler_Assemble(t *testing.T) {
	tree := fixtures.NewTree("Assembler", "")
	assert.NoError(t, tree.Populate())
	defer tree.Squash()

	ml, err := model.New(tree.CWD, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	names := []string{
		"file-two.bin", "one/file-two.bin", "one/file-three.bin",
	}

	// get manifests
	mx, err := ml.FeedManifests(true, true, true,
		lists.NewFileList(names...).ListDir(tree.CWD)...)
	assert.NoError(t, err)

	a, err := model.NewAssembler(ml)
	assert.NoError(t, err)

	for name, man := range mx {
		f, err := os.Open(filepath.Join(tree.CWD, name))
		assert.NoError(t, err)
		for _, chunk := range man.Chunks {
			buf := make([]byte, chunk.Size)
			_, err = f.Read(buf)
			assert.NoError(t, err)

			err = a.StoreChunk(bytes.NewReader(buf), chunk.ID)
			assert.NoError(t, err)
		}
	}

	// Kill some blobs
	tree.KillBLOB("file-two.bin")
	tree.KillBLOB("one/file-two.bin")


	err = a.Done(mx)
	assert.NoError(t, err)

	mx1, err := ml.FeedManifests(true, true, true,
		lists.NewFileList(names...).ListDir(tree.CWD)...)
	assert.NoError(t, err)

	assert.Equal(t, mx, mx1)
}