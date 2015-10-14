package model_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/barc/model"
)

func Test_Model_IsBlobs(t *testing.T)  {
	tree := fixtures.NewTree("is-blob", "")
	assert.NoError(t, tree.Populate())
	defer tree.Squash()

	names := lists.NewFileList().ListDir(tree.CWD)

	m, err := model.New(tree.CWD, false, 1024 * 1024, 16)
	assert.NoError(t, err)

	_, err = m.IsBlobs(names...)
	assert.NoError(t, err)
}