package lists_test

import (
	"github.com/akaspin/bar/client/lists"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Filelist1(t *testing.T) {
	lister := lists.NewFileList()
	assert.Equal(t, []string{"a", "b"}, lister.List([]string{
		"a",
		"b",
		".hidden",
		"deep/.hidden",
	}))
}
