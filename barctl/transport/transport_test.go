package transport_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/barctl/transport"
	"github.com/akaspin/bar/shadow"
	"os"
)

//func runServer(t *testing.T, root string) (endpoint *url.URL)  {
//	p := storage.NewStoragePool(storage.NewBlockStorageFactory(root, 2), 200, time.Minute)
//	port, err := fixtures.GetOpenPort()
//	assert.NoError(t, err)
//	go server.Serve(fmt.Sprintf(":%d", port), p)
//	endpoint, err = url.Parse(fmt.Sprintf("http://127.0.0.1:%d/v1", port))
//	assert.NoError(t, err)
//	return
//}

func Test_Upload(t *testing.T) {
	root := "fix-up-local"
	endpoint := fixtures.RunServer(t, root)
	defer os.RemoveAll(root)

	bn, err := fixtures.MakeBLOB(1024 * 1024 *2 + 56)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m, err := fixtures.NewShadowFromFile(bn, true, shadow.CHUNK_SIZE)

	tr := &transport.Transport{endpoint}
	err = tr.Push(bn, m)
	assert.NoError(t, err)

	m2, err := fixtures.NewShadowFromFile(
		fixtures.StoredName(root, m.ID), true, shadow.CHUNK_SIZE)
	assert.Equal(t, m.String(), m2.String())
}

func Test_Info(t *testing.T) {
	root := "fix-info-local"
	endpoint := fixtures.RunServer(t, root)
	defer os.RemoveAll(root)

	bn, err := fixtures.MakeBLOB(1024 * 1024 *2 + 56)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m, err := fixtures.NewShadowFromFile(bn, false, shadow.CHUNK_SIZE)
	tr := &transport.Transport{endpoint}

	err = tr.Info(m.ID)
	assert.NoError(t, err)
}

func Test_Exists(t *testing.T) {
	root := "fix-exists-local"
	endpoint := fixtures.RunServer(t, root)
	defer os.RemoveAll(root)

	bn1, err := fixtures.MakeBLOB(1024 * 1024 *2 + 56)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn1)

	bn2, err := fixtures.MakeBLOB(1024 * 1024 *2 + 58)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn2)

	m1, err := fixtures.NewShadowFromFile(bn1, true, shadow.CHUNK_SIZE)
	m2, err := fixtures.NewShadowFromFile(bn2, true, shadow.CHUNK_SIZE)

	tr := &transport.Transport{endpoint}
	err = tr.Push(bn1, m1)
	assert.NoError(t, err)

	r1, err := tr.Check([]string{
		m1.ID, m2.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, r1[0], m1.ID)
}
