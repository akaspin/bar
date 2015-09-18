package transport_test
import (
	"github.com/akaspin/bar/bard/storage"
	"time"
	"github.com/akaspin/bar/bard/server"
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"fmt"
	"net/url"
	"github.com/akaspin/bar/barctl/transport"
	"github.com/akaspin/bar/shadow"
	"os"
)

func runServer(t *testing.T, root string) (endpoint *url.URL)  {
	p := storage.NewStoragePool(storage.NewBlockStorageFactory(root, 2), 200, time.Minute)
	port, err := fixtures.GetOpenPort()
	assert.NoError(t, err)
	go server.Serve(fmt.Sprintf(":%d", port), p)
	endpoint, err = url.Parse(fmt.Sprintf("http://127.0.0.1:%d/v1", port))
	assert.NoError(t, err)
	return
}

func Test_Upload(t *testing.T) {
	root := "fix-up-local"
	endpoint := runServer(t, root)
	defer os.RemoveAll(root)

	bn, err := fixtures.MakeBLOB(1024 * 1024 *2 + 56)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m, err := shadow.NewShadowFromFile(bn, true, shadow.CHUNK_SIZE)

	tr := &transport.Transport{endpoint}
	err = tr.Push(bn, m)
	assert.NoError(t, err)

	m2, err := shadow.NewShadowFromFile(
		fixtures.StoredName(root, m.ID), true, shadow.CHUNK_SIZE)
	assert.Equal(t, m.String(), m2.String())
}

func Test_Info(t *testing.T) {
	root := "fix-info-local"
	endpoint := runServer(t, root)
	defer os.RemoveAll(root)

	bn, err := fixtures.MakeBLOB(1024 * 1024 *2 + 56)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m, err := shadow.NewShadowFromFile(bn, false, shadow.CHUNK_SIZE)
	tr := &transport.Transport{endpoint}

	err = tr.Info(m.ID)
	assert.NoError(t, err)
}
