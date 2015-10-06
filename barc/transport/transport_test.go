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


//func Test_Exists(t *testing.T) {
//	root := "fix-exists-local"
//	endpoint, stop := fixtures.RunServer(t, root)
//	defer stop()
//	defer os.RemoveAll(root)
//
//	bn1 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 56)
//	defer fixtures.KillBLOB(bn1)
//
//	bn2 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 58)
//	defer fixtures.KillBLOB(bn2)
//
//	m1, err := fixtures.NewShadowFromFile(bn1)
//	m2, err := fixtures.NewShadowFromFile(bn2)
//
//	tr := &transport.Transport{endpoint}
//	err = tr.Push(bn1, m1)
//	assert.NoError(t, err)
//
//	r1, err := tr.Check([]string{
//		m1.ID, m2.ID,
//	})
//	assert.NoError(t, err)
//	assert.Equal(t, r1[0], m1.ID)
//}
//
//func Test_DeclareCommitTx(t *testing.T) {
//	root := "fix-declare-commit-tx-local"
//	endpoint, stop := fixtures.RunServer(t, root)
//	defer stop()
//	defer os.RemoveAll(root)
//
//	bn1 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 56)
//	defer fixtures.KillBLOB(bn1)
//
//	bn2 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 58)
//	defer fixtures.KillBLOB(bn2)
//
//	m1, err := fixtures.NewShadowFromFile(bn1)
//	assert.NoError(t, err)
//	m2, err := fixtures.NewShadowFromFile(bn2)
//	assert.NoError(t, err)
//
//	tr := &transport.Transport{endpoint}
//	err = tr.Push(bn1, m1)
//	assert.NoError(t, err)
//
//	r1, err := tr.DeclareCommitTx("test", []string{
//		m1.ID, m2.ID,
//	})
//	assert.NoError(t, err)
//	assert.Equal(t, r1[0], m2.ID)
//}
//
//func Test_Transport_DownloadBLOB(t *testing.T) {
//	root := "fix-declare-commit-tx-local"
//	endpoint, stop := fixtures.RunServer(t, root)
//	defer stop()
//	defer os.RemoveAll(root)
//
//	bn1 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 56)
//	defer fixtures.KillBLOB(bn1)
//
//	m1, err := fixtures.NewShadowFromFile(bn1)
//	assert.NoError(t, err)
//
//	tr := &transport.Transport{endpoint}
//
//	// Upload fixture
//	err = tr.Push(bn1, m1)
//	assert.NoError(t, err)
//
//	// Download uploaded blob using temporary buffer
//	w := new(bytes.Buffer)
//	err = tr.GetBLOB(m1.ID, m1.Size, w)
//	assert.NoError(t, err)
//
//	// Make shadow
//	m2, err := manifest.NewFromAny(bytes.NewReader(w.Bytes()), manifest.CHUNK_SIZE)
//	assert.NoError(t, err)
//
//	assert.Equal(t, m1.ID, m2.ID)
//}
