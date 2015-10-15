package transport_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/barc/transport"
	"os"
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/barc/lists"
	"time"
	"fmt"
)

func Test_Ping(t *testing.T) {
	root := "testdata-Ping"
	endpoint, _, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	mod, err := model.New("", false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, endpoint, endpoint, 16)
	defer tr.Close()

	res, err := tr.Ping()
	assert.NoError(t, err)
	assert.Equal(t, int64(1024*1024*2), res.ChunkSize)
}

func Test_DeclareUpload(t *testing.T) {
	root := "testdata-DeclareUpload"
	endpoint, _, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("declare-upload", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	mod, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, endpoint, endpoint, 16)
	defer tr.Close()


	ml, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	toUp, err := tr.NewUpload(mx)
	assert.NoError(t, err)

	assert.Len(t, toUp, 3)
}

func Test_Upload(t *testing.T) {
	root := "testdata-Upload"
	endpoint, _, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("Upload", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	mod, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, endpoint, endpoint, 16)
	defer tr.Close()


	ml, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)
}

func Test_GetFetch(t *testing.T) {
	root := "testdata-GetFetch"
	endpoint, tEP, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("GetFetch", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	mod, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, endpoint, tEP, 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	check, err := tr.GetManifests(mx.IDMap().IDs())
	assert.NoError(t, err)
	assert.Len(t, check, 3)
}

func Test_Download(t *testing.T) {
	root := "Download"
	endpoint, rpcEP, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("Download", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	mod, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, endpoint, rpcEP, 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
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

func Test_Download_Many(t *testing.T) {
	root := "Download"
	endpoint, rpcEP, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("Download", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())
	assert.NoError(t, tree.PopulateN(10, 300))

	mod, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, endpoint, rpcEP, 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	// Kill some blobs
	tree.KillBLOB("file-two.bin")
	tree.KillBLOB("one/file-two.bin")
	tree.KillBLOB("one/file-three.bin")
	req := lists.Links{
		"file-two.bin": mx["file-two.bin"],
		"one/file-two.bin": mx["one/file-two.bin"],
		"one/file-three.bin": mx["one/file-three.bin"],
	}
	for i := 0; i < 256; i++ {
		nm := fmt.Sprintf("big/file-big-%d.bin", i)
		tree.KillBLOB(nm)
		req[nm] = mx[nm]
	}

	err = tr.Download(req)
	assert.NoError(t, err)
}

func Test_Check(t *testing.T) {
	root := "Check"
	endpoint, _, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree("Check", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	mod, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, endpoint, endpoint, 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	res, err := tr.Check([]proto.ID{
		"eebd7b0c388d7f4d4ede4681b472969d5f09228c0473010d670a6918a3c05e79",
		"eebd7b0c388d7f4d4ede4681b472969d5f09228c0473010d670a6918a3c05e7a",
	})
	assert.NoError(t, err)
	assert.Equal(t, []proto.ID{
		"eebd7b0c388d7f4d4ede4681b472969d5f09228c0473010d670a6918a3c05e79",
	}, res)
}

func Test_Spec(t *testing.T) {
	root := "Spec"
	endpoint, _, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tree := fixtures.NewTree(root, "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())

	mod, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, endpoint, endpoint, 16)
	defer tr.Close()

	ml, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	// make spec
	nameMap := map[string]proto.ID{}
	for name, m := range mx {
		nameMap[name] = m.ID
	}

	spec1, err := proto.NewSpec(time.Now().UnixNano(), nameMap, []string{})
	assert.NoError(t, err)

	err = tr.UploadSpec(spec1)
	assert.NoError(t, err)

	_, err = tr.GetSpec(spec1.ID)
	assert.NoError(t, err)
}


