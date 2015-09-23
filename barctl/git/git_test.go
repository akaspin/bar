package git_test
import (
	"testing"
	"github.com/akaspin/bar/barctl/git"
	"github.com/stretchr/testify/assert"
)

func Test_DirtyFiles(t *testing.T) {
	g, err := git.NewGit("")
	assert.NoError(t, err)

	dirty, err := g.DirtyFiles()
	assert.NoError(t, err)
	for _, f := range dirty {
		assert.NotEqual(t, "", f)
	}
}

func Test_FilterByDiff(t *testing.T) {
	g, err := git.NewGit("")
	assert.NoError(t, err)

	res, err := g.FilterByDiff("unspecified", []string{
		"barctl/cmd/git-cat.go",
		"barctl/cmd/git-clean.go",
	}...)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"barctl/cmd/git-cat.go",
		"barctl/cmd/git-clean.go",
	}, res)
}
