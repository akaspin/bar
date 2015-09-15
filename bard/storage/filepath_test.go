package storage_test
import (
	"testing"
	"path/filepath"
	"github.com/stretchr/testify/assert"
)

func Test_Filepath_EmptyChunk(t *testing.T) {
	id := "third"
	res := filepath.Join("first/second", id[:0], id)
	assert.Equal(t, "first/second/third", res)
}
