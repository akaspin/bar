package lists_test
import (
	"testing"
	"github.com/akaspin/bar/bar/lists"
	"github.com/stretchr/testify/assert"
	"path/filepath"
)


func Test_Mapper_ToRoot(t *testing.T) {
	root := "/home/user"
	cwd := "/home/user/wd"
	mapper := lists.NewMapper(cwd, root)

	res, err := mapper.ToRoot("a1/b", "a2")

	assert.NoError(t, err)
	assert.Equal(t, []string{"wd/a1/b", "wd/a2"}, res)

}

func Test_Mapper_FromRoot(t *testing.T) {
	root := "/home/user"
	cwd := "/home/user/wd"
	mapper := lists.NewMapper(cwd, root)

	res, err := mapper.FromRoot("wd/a1/b", "wd/a2")

	assert.NoError(t, err)
	assert.Equal(t, []string{"a1/b", "a2"}, res)
}

func Test_Mapper_Wrap(t *testing.T) {
	data := []string{"file/one", "file/with spaces.bin", "dir/with spaces/file"}
	var res []string
	for _, i := range data {
		res = append(res, filepath.FromSlash(i))
	}

	t.Log(res)
}