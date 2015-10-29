package storage_test

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_Filepath_EmptyChunk(t *testing.T) {
	id := "third"
	res := filepath.Join("first/second", id[:0], id)
	assert.Equal(t, "first/second/third", res)
}
