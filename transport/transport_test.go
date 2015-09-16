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
	"github.com/akaspin/bar/transport"
	"github.com/akaspin/bar/shadow"
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
	endpoint := runServer(t, "fix-up-local")
	t.Log(endpoint)
	tr := &transport.Transport{endpoint}

	bn, err := fixtures.MakeBLOB(1024 * 1024 *2 + 56)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m, err := shadow.NewShadow(bn)

	err = tr.Push(bn, m)
	assert.NoError(t, err)
}
