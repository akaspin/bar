package storage_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/server/storage"
	"github.com/akaspin/bar/bar/lists"
	"time"
	"os"
	"github.com/tamtam-im/logx"
	"bytes"
	"github.com/nu7hatch/gouuid"
)


func Test_Storage_Upload_CreateUpload(t *testing.T)  {
	logx.SetLevel(logx.DEBUG)

	tree := fixtures.NewTree("create-upload", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	m, err :=  model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	os.RemoveAll("testdata/create-upload-storage")
	stor := storage.NewBlockStorage(&storage.BlockStorageOptions{
		"testdata/create-upload-storage", 2, 16, 32,
	})
	defer os.RemoveAll("testdata/create-upload-storage")

	names := lists.NewFileList().ListDir(tree.CWD)
	mans, err := m.FeedManifests(true, false, true, names...)

	uID, _ := uuid.NewV4()
	missing, err := stor.CreateUploadSession(*uID, mans.GetManifestSlice(), time.Hour)
	assert.NoError(t, err)
	assert.Len(t, missing, 4)
}

func Test_Storage_Upload_UploadChunk(t *testing.T)  {
	logx.SetLevel(logx.DEBUG)

	tree := fixtures.NewTree("upload-chunk", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	m, err :=  model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	os.RemoveAll("testdata/upload-chunk-storage")
	stor := storage.NewBlockStorage(&storage.BlockStorageOptions{
		"testdata/upload-chunk-storage", 2, 16, 32,
	})
	defer os.RemoveAll("testdata/upload-chunk-storage")

	names := lists.NewFileList().ListDir(tree.CWD)
	mans, err := m.FeedManifests(true, false, true, names...)

	uID, _ := uuid.NewV4()
	missing, err := stor.CreateUploadSession(*uID, mans.GetManifestSlice(), time.Hour)
	assert.NoError(t, err)

	toUpload := mans.GetChunkLinkSlice(missing)

	for _, v := range toUpload {
		r, err := os.Open(tree.BlobFilename(v.Name))
		assert.NoError(t, err)
		defer r.Close()

		buf := make([]byte, v.Size)
		_, err = r.ReadAt(buf, v.Offset)

		err = stor.UploadChunk(*uID, v.Chunk.ID, bytes.NewReader(buf))
		assert.NoError(t, err)
	}
}

func Test_Storage_Upload_FinishUpload(t *testing.T)  {
	logx.SetLevel(logx.DEBUG)

	tree := fixtures.NewTree("finish-upload", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	m, err :=  model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	os.RemoveAll("testdata/finish-upload-storage")
	stor := storage.NewBlockStorage(&storage.BlockStorageOptions{
		"testdata/finish-upload-storage", 2, 16, 32,
	})
	defer os.RemoveAll("testdata/finish-upload-storage")

	names := lists.NewFileList().ListDir(tree.CWD)
	mans, err := m.FeedManifests(true, false, true, names...)

	uID, _ := uuid.NewV4()
	missing, err := stor.CreateUploadSession(*uID, mans.GetManifestSlice(), time.Hour)
	assert.NoError(t, err)

	toUpload := mans.GetChunkLinkSlice(missing)

	for _, v := range toUpload {
		r, err := os.Open(tree.BlobFilename(v.Name))
		assert.NoError(t, err)
		defer r.Close()

		buf := make([]byte, v.Size)
		_, err = r.ReadAt(buf, v.Offset)

		err = stor.UploadChunk(*uID, v.Chunk.ID, bytes.NewReader(buf))
		assert.NoError(t, err)
	}

	err = stor.FinishUploadSession(*uID)
	assert.NoError(t, err)
}