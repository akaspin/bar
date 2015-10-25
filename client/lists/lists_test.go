package lists_test
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/client/lists"
)

func Test_Filelist1(t *testing.T)  {
	lister := lists.NewFileList()
	assert.Equal(t, []string{"a", "b"}, lister.List([]string{
		"a",
		"b",
		".hidden",
		"deep/.hidden",
	}))
}
