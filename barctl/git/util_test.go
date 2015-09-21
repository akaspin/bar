package git_test
import (
	"testing"
	"github.com/akaspin/bar/barctl/git"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
)

func Test_GetGitTop(t *testing.T) {
	res, err := git.GetGitTop()
	assert.NoError(t, err)
	cwd, _ := os.Getwd()
	assert.Equal(t, res, filepath.Clean(filepath.Join(cwd, "../../../bar")))
}
