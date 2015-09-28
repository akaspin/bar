package fixtures_test
import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/fixtures"
	"path/filepath"
)


func Test_MakeTree(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	tree := fixtures.NewTree(filepath.Join(wd, "testdata/pure"))
	err = tree.Populate()
	assert.NoError(t, err)

	err = tree.Squash()
	assert.NoError(t, err)
}

func Test_MakeTreeGit(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	tree := fixtures.NewTree(filepath.Join(wd, "testdata/git"))
	err = tree.Populate()
	assert.NoError(t, err)

	err = tree.InitGit()
	assert.NoError(t, err)

	err = tree.Squash()
	assert.NoError(t, err)
}
