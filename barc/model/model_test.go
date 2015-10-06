package model_test
import (
	"testing"
	"os"
	"github.com/akaspin/bar/fixtures"
	"github.com/akaspin/bar/barc/model"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/barc/lists"
	"path/filepath"
)

func Test_Model_CollectManifests(t *testing.T)  {
	wd, _ := os.Getwd()
	root := filepath.Join(wd, "tree1")

	tree := fixtures.NewTree(root)
	assert.NoError(t, tree.Populate())
	defer tree.Squash()

	names := lists.NewFileList().ListDir(root)

	m, err := model.New(root, false, 1024 * 1024, 16)
	assert.NoError(t, err)
	lx, err := m.CollectManifests(true, true, names...)
	assert.NoError(t, err)

	assert.Len(t, lx.IDMap(), 3)

}