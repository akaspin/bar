package transport_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/barc/transport"
	"os"
	"path/filepath"
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/lists"
)

func Test_Ping(t *testing.T) {
//	logx.SetLevel(logx.DEBUG)
	root := "fix-up-ping"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tr := transport.NewTransport("", endpoint.String(), 16)
	defer tr.Close()

	res, err := tr.Ping()
	assert.NoError(t, err)
	assert.Equal(t, int64(1024*1024*2), res.ChunkSize)
}

func Test_DeclareUpload(t *testing.T) {
	root := "fix-up-ping"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)
	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "testdata")

	tr := transport.NewTransport(wd, endpoint.String(), 16)
	defer tr.Close()

	tree := fixtures.NewTree(wd)
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	ml, err := model.New(wd, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(wd)...)
	assert.NoError(t, err)

	toUp, err := tr.NewUpload(mx)
	assert.NoError(t, err)

	assert.Len(t, toUp, 3)
}

func Test_Upload(t *testing.T) {
	root := "fix-up-ping"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "testdata")
	tree := fixtures.NewTree(wd)
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(wd, endpoint.String(), 16)
	defer tr.Close()


	ml, err := model.New(wd, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(wd)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)
}

func Test_GetFetch(t *testing.T) {
	root := "fix-up-ping"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "testdata")
	tree := fixtures.NewTree(wd)
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(wd, endpoint.String(), 16)
	defer tr.Close()

	ml, err := model.New(wd, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(wd)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	check, err := tr.GetFetch(mx.IDMap().IDs())
	assert.NoError(t, err)
	assert.Len(t, check, 3)
}

func Test_Download(t *testing.T) {
	root := "fix-download"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "testdata")
	tree := fixtures.NewTree(wd)
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(wd, endpoint.String(), 16)
	defer tr.Close()

	ml, err := model.New(wd, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(wd)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	// Kill some blobs
	tree.KillBLOB("file-two.bin")
	tree.KillBLOB("one/file-two.bin")
	tree.KillBLOB("one/file-three.bin")

	err = tr.Download(model.Links{
		"file-two.bin": mx["file-two.bin"],
		"one/file-two.bin": mx["one/file-two.bin"],
		"one/file-three.bin": mx["one/file-three.bin"],
	})
	assert.NoError(t, err)
}