package git_test
import (
	"testing"
	"github.com/akaspin/bar/barctl/git"
	"github.com/stretchr/testify/assert"
)

func Test_IsDirty1(t *testing.T)  {
	res, err := git.IsClean("", "barctl/git/util_test.go")
	assert.NoError(t, err)
	assert.True(t, res)
}

func Test_GetOID(t *testing.T)  {
	res, err := git.GetFileOID("", "barctl/git/util_test.go")
	assert.NoError(t, err)
	assert.NotEqual(t, "", res)
}

func Test_BLOBsClean(t *testing.T) {
	ok, err := git.IsBLOBsClean("", "unspecified", []string{
		"barctl/cmd/git-cat.go",
		"barctl/cmd/git-clean.go",
		"",
	})
	assert.NoError(t, err)
	t.Log(ok)
}
