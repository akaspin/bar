package model_test

import (
	"github.com/akaspin/bar/client/lists"
	"github.com/akaspin/bar/client/model"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Model_IsBlobs(t *testing.T) {
	tree := fixtures.NewTree("is-blob", "")
	assert.NoError(t, tree.Populate())
	defer tree.Squash()

	names := lists.NewFileList().ListDir(tree.CWD)

	m, err := model.New(tree.CWD, false, 1024*1024, 16)
	assert.NoError(t, err)

	_, err = m.IsBlobs(names...)
	assert.NoError(t, err)
}
