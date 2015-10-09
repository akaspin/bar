package transport_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/barc/transport"
	"os"
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/proto"
)

func Test_Ping(t *testing.T) {
//	logx.SetLevel(logx.DEBUG)
	root := "testdata-Ping"
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
	root := "testdata-DeclareUpload"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("declare-upload", "")
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(tree.CWD, endpoint.String(), 16)
	defer tr.Close()


	ml, err := model.New(tree.CWD, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	toUp, err := tr.NewUpload(mx)
	assert.NoError(t, err)

	assert.Len(t, toUp, 3)
}

func Test_Upload(t *testing.T) {
	root := "testdata-Upload"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("testdata-Upload", "")
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(tree.CWD, endpoint.String(), 16)
	defer tr.Close()


	ml, err := model.New(tree.CWD, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)
}

func Test_GetFetch(t *testing.T) {
	root := "testdata-GetFetch"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("GetFetch", "")
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(tree.CWD, endpoint.String(), 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	check, err := tr.GetFetch(mx.IDMap().IDs())
	assert.NoError(t, err)
	assert.Len(t, check, 3)
}

func Test_Download(t *testing.T) {
	root := "testdata-Download"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("testdata-Download", "")
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(tree.CWD, endpoint.String(), 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	// Kill some blobs
	tree.KillBLOB("file-two.bin")
	tree.KillBLOB("one/file-two.bin")
	tree.KillBLOB("one/file-three.bin")

	err = tr.Download(lists.Links{
		"file-two.bin": mx["file-two.bin"],
		"one/file-two.bin": mx["one/file-two.bin"],
		"one/file-three.bin": mx["one/file-three.bin"],
	})
	assert.NoError(t, err)
}

func Test_Check(t *testing.T) {
	root := "testdata-Check"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("testdata-Check", "")
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(tree.CWD, endpoint.String(), 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	res, err := tr.Check([]string{
		"eebd7b0c388d7f4d4ede4681b472969d5f09228c0473010d670a6918a3c05e79",
		"eebd7b0c388d7f4d4ede4681b472969d5f09228c0473010d670a6918a3c05e7a",
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"eebd7b0c388d7f4d4ede4681b472969d5f09228c0473010d670a6918a3c05e79",
	}, res)
}

func Test_Spec(t *testing.T) {
	root := "testdata-Spec"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree(root, "")
	assert.NoError(t, tree.Populate())

	tr := transport.NewTransport(tree.CWD, endpoint.String(), 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, manifest.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.CollectManifests(true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	// make spec
	nameMap := map[string]string{}
	for name, m := range mx {
		nameMap[name] = m.ID
	}

	spec1, err := proto.NewSpec(nameMap, []string{})
	assert.NoError(t, err)

	err = tr.UploadSpec(spec1)
	assert.NoError(t, err)

	_, err = tr.GetSpec(spec1.ID)
	assert.NoError(t, err)
}


