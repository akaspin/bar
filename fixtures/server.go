package fixtures
import (
	"testing"
	"net/url"
	"github.com/akaspin/bar/bard/storage"
	"time"
	"github.com/stretchr/testify/assert"
"github.com/akaspin/bar/bard/server"
	"fmt"
	"path/filepath"
)

func RunServer(t *testing.T, root string) (endpoint *url.URL)  {
	p := storage.NewStoragePool(storage.NewBlockStorageFactory(root, 2), 200, time.Minute)
	port, err := GetOpenPort()
	assert.NoError(t, err)
	go server.Serve(fmt.Sprintf(":%d", port), p)
	endpoint, err = url.Parse(fmt.Sprintf("http://127.0.0.1:%d/v1", port))
	assert.NoError(t, err)
	return
}

func ServerStoredName(root string, id string) string {
	return filepath.Join(root, id[:2], id)
}
