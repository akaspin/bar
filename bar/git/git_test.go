package git_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/bar/git"
	"github.com/akaspin/bar/bar/lists"
	"os"
	"path/filepath"
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/proto"
)

func Test_Git_LsTree(t *testing.T) {
	tree := fixtures.NewTree("Git_LsTree", "")
	assert.NoError(t, tree.Populate())
	defer tree.Squash()
	g, err := gitFixture(tree)
	assert.NoError(t, err)

	res1, err := g.LsTree("HEAD", "in-other", "one", "two/file-four with spaces.bin")
	assert.NoError(t, err)
	assert.EqualValues(t,
		[]string{
			"one/file-four with spaces.bin",
			"one/file-one.bin",
			"one/file-three.bin",
			"one/file-two.bin",
			"two/file-four with spaces.bin"}, res1)
	res2, err := g.LsTree("other", "in-other", "one")
	assert.NoError(t, err)
	assert.EqualValues(t,
		[]string{
			"in-other/blob.bin",
			"one/file-four with spaces.bin",
			"one/file-one.bin",
			"one/file-three.bin",
			"one/file-two.bin"}, res2)
}

func Test_Git_DiffFiles(t *testing.T) {
	tree := fixtures.NewTree("Git_DiffFiles", "")
	assert.NoError(t, tree.Populate())
	defer tree.Squash()
	g, err := gitFixture(tree)
	assert.NoError(t, err)

	fixtures.MakeNamedBLOB(filepath.Join(tree.CWD, "in-other", "blob.bin"), 10)
	fixtures.MakeNamedBLOB(filepath.Join(tree.CWD, "two", "file-four with spaces.bin"), 110)
	dirty, err := g.DiffFiles(
		filepath.Join("in-other", "blob.bin"),
		filepath.Join("two", "file-four with spaces.bin"),
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{"two/file-four with spaces.bin"}, dirty)
}

func Test_Git_DiffIndex(t *testing.T) {
	tree := fixtures.NewTree("Git_DiffFiles", "")
	assert.NoError(t, tree.Populate())
	defer tree.Squash()
	g, err := gitFixture(tree)
	assert.NoError(t, err)

	fixtures.MakeNamedBLOB(filepath.Join(tree.CWD, "in-other", "blob.bin"), 10)
	fixtures.MakeNamedBLOB(filepath.Join(tree.CWD, "two", "file-four with spaces.bin"), 110)
	g.Run("add -A")

	dirty, err := g.DiffIndex(
		filepath.Join("in-other", "blob.bin"),
		filepath.Join("two", "file-four with spaces.bin"),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"two/file-four with spaces.bin"}, dirty)
}

func Test_Git_Checkout_Files(t *testing.T) {
	tree := fixtures.NewTree("Git_Checkout_files", "")
	assert.NoError(t, tree.Populate())
	defer tree.Squash()
	g, err := gitFixture(tree)
	assert.NoError(t, err)
	err = g.Checkout("other",
		filepath.Join("in-other", "blob.bin"),
		filepath.Join("two", "file-four with spaces.bin"),
	)
	assert.NoError(t, err)

	// assert on branch master
	current, _, err := g.GetBranches()
	assert.NoError(t, err)
	assert.Equal(t, "master", current)

	diff, err := g.DiffIndex(
		filepath.Join("in-other", "blob.bin"),
		filepath.Join("two", "file-four with spaces.bin"),
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"in-other/blob.bin",
		"two/file-four with spaces.bin"}, diff)
}

func Test_Git_Divert1(t *testing.T)  {

	tree := fixtures.NewTree("Git_divert", "")
	assert.NoError(t, tree.Populate())
	defer tree.Squash()

	g, err := gitFixture(tree)
	assert.NoError(t, err)

//	logx.SetLevel(logx.TRACE)

	// get blobmap for further checks
	mod, err := model.New(tree.CWD, false, proto.CHUNK_SIZE, 16)
	assert.NoError(t, err)

	names := lists.NewFileList().ListDir(tree.CWD)
	mans1, err := mod.FeedManifests(true, true, true, names...)
	assert.NoError(t, err)

	// Run divert on "in-other" and "one"
	divert := git.NewDivert(g)
	err = divert.Begin("other", "in-other", "one", "two/file-four with spaces.bin")
	assert.NoError(t, err)

	// Make two blobs and collect their manifests
	bn1 := filepath.Join(tree.CWD, "in-other", "blob.bin")
	bn2 := filepath.Join(tree.CWD, "one", "file-one.bin")
	fixtures.MakeNamedBLOB(bn1, 110)
	fixtures.MakeNamedBLOB(bn2, 200)

	oMan1, err := fixtures.NewShadowFromFile(bn1)
	assert.NoError(t, err)
	oMan2, err := fixtures.NewShadowFromFile(bn2)
	assert.NoError(t, err)

	// commit
	spec, err := divert.ReadSpec()
	assert.NoError(t, err)

	err = divert.Commit(spec, "from-master")
	assert.NoError(t, err)

	err = divert.Cleanup(spec)
	assert.NoError(t, err)

	err = divert.CleanSpec()
	assert.NoError(t, err)

	// Final checks
	branch, _, err := g.GetBranches()
	assert.NoError(t, err)
	assert.Equal(t, "master", branch)

	// check master files
	names = lists.NewFileList().ListDir(tree.CWD)
	mans2, err := mod.FeedManifests(true, true, true, names...)
	assert.NoError(t, err)
	assert.EqualValues(t, mans1, mans2)

	// check stored branch
	err = g.Checkout("other")
	assert.NoError(t, err)

	oMan1p, err := fixtures.NewShadowFromFile(bn1)
	assert.NoError(t, err)
	assert.EqualValues(t, oMan1, oMan1p)

	oMan2p, err := fixtures.NewShadowFromFile(bn2)
	assert.NoError(t, err)
	assert.EqualValues(t, oMan2, oMan2p)
}

func gitFixture(tree *fixtures.Tree) (res *git.Git, err error)  {
	res = &git.Git{lists.NewMapper(tree.CWD, tree.CWD)}
	if _, err = res.Run("init"); err != nil {
		return
	}
	if _, err = res.Run("add", "-A"); err != nil {
		return
	}
	if _, err = res.Run("commit", "-m", "master1"); err != nil {
		return
	}
	if _, err = res.Run("checkout", "-b", "other"); err != nil {
		return
	}
	os.MkdirAll(filepath.Join(tree.CWD, "in-other"), 0755)
	fixtures.MakeNamedBLOB(filepath.Join(tree.CWD, "in-other", "blob.bin"), 1670)
	fixtures.MakeNamedBLOB(filepath.Join(tree.CWD, "one", "file-one.bin"), 170)
	fixtures.MakeNamedBLOB(filepath.Join(tree.CWD, "two", "file-four with spaces.bin"), 170)
	if _, err = res.Run("add", "-A"); err != nil {
		return
	}
	if _, err = res.Run("commit", "-m", "other1"); err != nil {
		return
	}
	if _, err = res.Run("checkout", "master"); err != nil {
		return
	}
	return
}