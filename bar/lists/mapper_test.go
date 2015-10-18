package lists_test
import (
	"testing"
	"github.com/akaspin/bar/bar/lists"
	"github.com/stretchr/testify/assert"
)


func Test_MapperToRoot(t *testing.T) {
	root := "/home/user"
	cwd := "/home/user/wd"
	mapper := lists.NewMapper(cwd, root)

	res, err := mapper.ToRoot("a1/b", "a2")

	assert.NoError(t, err)
	assert.Equal(t, []string{"wd/a1/b", "wd/a2"}, res)

}

func Test_MapperFromRoot(t *testing.T) {
	root := "/home/user"
	cwd := "/home/user/wd"
	mapper := lists.NewMapper(cwd, root)

	res, err := mapper.FromRoot("wd/a1/b", "wd/a2")

	assert.NoError(t, err)
	assert.Equal(t, []string{"a1/b", "a2"}, res)

}