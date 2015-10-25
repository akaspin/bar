package transport_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/client/transport"
	"github.com/akaspin/bar/client/model"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/client/lists"
	"time"
	"fmt"
)

func seed(t *testing.T, root string) (halt func(), tree *fixtures.Tree, ml *model.Model, srv *fixtures.FixtureServer, trans *transport.Transport) {
	tree = fixtures.NewTree(root, "")
	assert.NoError(t, tree.Populate())

	ml, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 32)
	assert.NoError(t, err)

	srv, err = fixtures.NewFixtureServer(root)
	assert.NoError(t, err)
	trans = transport.NewTransport(ml, srv.HTTPEndpoint, srv.RPCEndpoints[0], 16)

	halt = func() {
		trans.Close()
		srv.Stop()
		tree.Squash()
	}
	return
}

func Test_Transport_ServerInfo(t *testing.T) {
	root := "testdata-Ping"

	srv, err := fixtures.NewFixtureServer(root)
	assert.NoError(t, err)
	defer srv.Stop()
	mod, err := model.New("", false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)
	tr := transport.NewTransport(mod, srv.HTTPEndpoint, srv.RPCEndpoints[0], 16)
	defer tr.Close()

	res, err := tr.ServerInfo()
	assert.NoError(t, err)
	assert.Equal(t, int64(1024*1024*2), res.ChunkSize)
}

func Test_Transport_CreateUpload(t *testing.T) {
	root := "DeclareUpload"

	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	upload := transport.NewUpload(tr, time.Hour)

	toUp, err := upload.SendCreateUpload(mx)
	assert.NoError(t, err)
	assert.Len(t, toUp, 4)
}

func Test_Transport_UploadChunk(t *testing.T) {
	root := "UploadChunk"

	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	upload := transport.NewUpload(tr, time.Hour)

	missing, err := upload.SendCreateUpload(mx)
	assert.NoError(t, err)

	toUp := mx.GetChunkLinkSlice(missing)
	for _, tu := range toUp {
		err = upload.UploadChunk(tu.Name, tu.Chunk)
		assert.NoError(t, err)
	}
}

func Test_Transport_FinishUpload(t *testing.T) {
	root := "FinishUpload"

	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	upload := transport.NewUpload(tr, time.Hour)

	missing, err := upload.SendCreateUpload(mx)
	assert.NoError(t, err)

	toUp := mx.GetChunkLinkSlice(missing)
	for _, tu := range toUp {
		err = upload.UploadChunk(tu.Name, tu.Chunk)
		assert.NoError(t, err)
	}
	assert.NoError(t, upload.Commit())
}

func Test_Transport_Upload(t *testing.T) {
	root := "Upload"
	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)
}

func Test_Transport_Download(t *testing.T) {
	root := "Download"

	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	// Kill some blobs
	tree.KillBLOB("file-two.bin")
	tree.KillBLOB("one/file-two.bin")
	tree.KillBLOB("one/file-three.bin")

	err = tr.Download(lists.BlobMap{
		"file-two.bin": mx["file-two.bin"],
		"one/file-two.bin": mx["one/file-two.bin"],
		"one/file-three.bin": mx["one/file-three.bin"],
	})
	assert.NoError(t, err)
}

func Test_Transport_Download_Many(t *testing.T) {
//	t.Skip()
	root := "Download-many"

	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

	assert.NoError(t, tree.PopulateN(10, 1000))

	mx, err := ml.FeedManifests(true, true, true, lists.NewFileList().ListDir(tree.CWD)...)
	assert.NoError(t, err)

	err = tr.Upload(mx)
	assert.NoError(t, err)

	// Kill some blobs
	tree.KillBLOB("file-two.bin")
	tree.KillBLOB("one/file-two.bin")
	tree.KillBLOB("one/file-three.bin")
	req := lists.BlobMap{
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

func Test_Transport_Check(t *testing.T) {
	root := "Check"
	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

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
		"eebd7b0c388d7f4d4ede4681b472969d5f09228c0473010d670a6918a3c05e7a",
	}, res)
}

func Test_Transport_UploadSpec(t *testing.T) {
//	t.Skip()
	root := "Spec"
	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

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

//	_, err = tr.GetSpec(spec1.ID)
//	assert.NoError(t, err)
}

func Test_Transport_FetchSpec(t *testing.T) {
//	t.Skip()
	root := "Spec"
	halt, tree, ml, _, tr := seed(t, root)
	defer halt()

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

	spec2, err := tr.GetSpec(spec1.ID)
	assert.NoError(t, err)

	assert.Equal(t, spec1.ID, spec2.ID)
}


