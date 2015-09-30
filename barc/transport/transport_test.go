package transport_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/barc/transport"
	"os"
	"bytes"
	"github.com/akaspin/bar/shadow"
)

func Test_Ping(t *testing.T) {
//	logx.SetLevel(logx.DEBUG)
	root := "fix-up-ping"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	tr := &transport.Transport{endpoint}
	res, err := tr.Ping()
	assert.NoError(t, err)
	assert.Equal(t, int64(1024*1024*2), res.ChunkSize)
}

func Test_Upload(t *testing.T) {
	root := "fix-up-local"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	bn := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 56)
	defer fixtures.KillBLOB(bn)

	m, err := fixtures.NewShadowFromFile(bn)

	tr := &transport.Transport{endpoint}
	err = tr.Push(bn, m)
	assert.NoError(t, err)

	m2, err := fixtures.NewShadowFromFile(
		fixtures.ServerStoredName(root, m.ID))
	assert.Equal(t, m.String(), m2.String())
}

func Test_Info(t *testing.T) {
	root := "fix-info-local"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	bn := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 56)
	defer fixtures.KillBLOB(bn)

	m, err := fixtures.NewShadowFromFile(bn)
	tr := &transport.Transport{endpoint}

	err = tr.Info(m.ID)
	assert.NoError(t, err)
}

func Test_Exists(t *testing.T) {
	root := "fix-exists-local"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	bn1 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 56)
	defer fixtures.KillBLOB(bn1)

	bn2 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 58)
	defer fixtures.KillBLOB(bn2)

	m1, err := fixtures.NewShadowFromFile(bn1)
	m2, err := fixtures.NewShadowFromFile(bn2)

	tr := &transport.Transport{endpoint}
	err = tr.Push(bn1, m1)
	assert.NoError(t, err)

	r1, err := tr.Check([]string{
		m1.ID, m2.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, r1[0], m1.ID)
}

func Test_DeclareCommitTx(t *testing.T) {
	root := "fix-declare-commit-tx-local"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	bn1 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 56)
	defer fixtures.KillBLOB(bn1)

	bn2 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 58)
	defer fixtures.KillBLOB(bn2)

	m1, err := fixtures.NewShadowFromFile(bn1)
	assert.NoError(t, err)
	m2, err := fixtures.NewShadowFromFile(bn2)
	assert.NoError(t, err)

	tr := &transport.Transport{endpoint}
	err = tr.Push(bn1, m1)
	assert.NoError(t, err)

	r1, err := tr.DeclareCommitTx("test", []string{
		m1.ID, m2.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, r1[0], m1.ID)
}

func Test_Transport_DownloadBLOB(t *testing.T) {
	root := "fix-declare-commit-tx-local"
	endpoint, stop := fixtures.RunServer(t, root)
	defer stop()
	defer os.RemoveAll(root)

	bn1 := fixtures.MakeBLOB(t, 1024 * 1024 *2 + 56)
	defer fixtures.KillBLOB(bn1)

	m1, err := fixtures.NewShadowFromFile(bn1)
	assert.NoError(t, err)

	tr := &transport.Transport{endpoint}

	// Upload fixture
	err = tr.Push(bn1, m1)
	assert.NoError(t, err)

	// Download uploaded blob using temporary buffer
	w := new(bytes.Buffer)
	err = tr.GetBLOB(m1.ID, m1.Size, w)
	assert.NoError(t, err)

	// Make shadow
	m2, err := shadow.New(bytes.NewReader(w.Bytes()), m1.Size)
	assert.NoError(t, err)

	assert.Equal(t, m1.ID, m2.ID)
}
