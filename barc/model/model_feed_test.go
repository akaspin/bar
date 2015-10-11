package model_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/akaspin/bar/barc/model"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/barc/lists"
)


func Test_Model_FeedManifests(t *testing.T)  {
	tree := fixtures.NewTree("feed-manifests", "")
	assert.NoError(t, tree.Populate())

	names := lists.NewFileList().ListDir(tree.CWD)

	m, err := model.New(tree.CWD, false, 1024 * 1024, 16)
	assert.NoError(t, err)
	lx, err := m.FeedManifests(true, true, true, names...)
	assert.NoError(t, err)

	assert.Len(t, lx.Names(), 16)
}

func Test_Model_FeedManifests_Nil(t *testing.T)  {
	tree := fixtures.NewTree("feed-manifests", "")
	assert.NoError(t, tree.Populate())

	names := lists.NewFileList().ListDir(tree.CWD)
	tree.KillBLOB("file-one.bin")

	m, err := model.New(tree.CWD, false, 1024 * 1024, 16)
	assert.NoError(t, err)
	lx, err := m.FeedManifests(true, true, false, names...)
	assert.Error(t, err)
	assert.Len(t, lx.Names(), 15)
}

func Test_Model_FeedManifests_Large(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	tree := fixtures.NewTree("collect-manifests-large", "")
	defer tree.Squash()

	assert.NoError(t, tree.Populate())
	assert.NoError(t, tree.PopulateN(1024 * 1024 * 300, 1))

	names := lists.NewFileList().ListDir(tree.CWD)

	m, err := model.New(tree.CWD, false, 1024 * 1024, 16)
	assert.NoError(t, err)
	lx, err := m.FeedManifests(true, true, true, names...)
	assert.NoError(t, err)

	assert.Len(t, lx.Names(), 17)
}

func Test_Model_FeedManifests_Many(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	tree := fixtures.NewTree("collect-manifests-large", "")
	defer tree.Squash()
	assert.NoError(t, tree.Populate())
	assert.NoError(t, tree.PopulateN(10, 1000))

	names := lists.NewFileList().ListDir(tree.CWD)

	m, err := model.New(tree.CWD, false, 1024 * 1024, 16)
	assert.NoError(t, err)
	lx, err := m.FeedManifests(true, true, true, names...)
	assert.NoError(t, err)

	assert.Len(t, lx.Names(), 1016)
}

func Benchmark_FeedManifests_Many(b *testing.B)  {
	n := 100000

	tree := fixtures.NewTree("collect-manifests-many-B", "")
	defer tree.Squash()
	assert.NoError(b, tree.Populate())
	assert.NoError(b, tree.PopulateN(10, n))
	names := lists.NewFileList().ListDir(tree.CWD)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		m, err := model.New(tree.CWD, false, 1024 * 1024, 16)
		assert.NoError(b, err)
		lx, err := m.FeedManifests(true, true, true, names...)
		b.Log(len(lx), err)
		assert.NoError(b, err)
		b.StopTimer()
	}
}

func Benchmark_FeedManifests_Large(b *testing.B)  {
	tree := fixtures.NewTree("collect-manifests-large-B", "")
	defer tree.Squash()
	assert.NoError(b, tree.Populate())
	assert.NoError(b, tree.PopulateN(1024 * 1024 * 500, 5))

	names := lists.NewFileList().ListDir(tree.CWD)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		m, err := model.New(tree.CWD, false, 1024 * 1024, 16)
		assert.NoError(b, err)
		_, err = m.FeedManifests(true, true, true, names...)
		assert.NoError(b, err)
		b.StopTimer()
	}
}