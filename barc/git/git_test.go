package git_test
import (
	"testing"
	"github.com/akaspin/bar/barc/git"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/fixtures"
	"os"
	"path/filepath"
)

func Test_DirtyFiles(t *testing.T) {
	g, err := git.NewGit("")
	assert.NoError(t, err)

	dirty, err := g.DiffFiles()
	assert.NoError(t, err)
	for _, f := range dirty {
		assert.NotEqual(t, "", f)
	}
}

func Test_FilterByDiff(t *testing.T) {
	wd, _ := os.Getwd()
	cm := filepath.Clean(filepath.Join(wd, "../../barc"))

	g, err := git.NewGit(cm)
	assert.NoError(t, err)

	res, err := g.FilterByAttr("unspecified", []string{
		"cmd/git-cat.go",
		"cmd/git-clean.go",
	}...)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"cmd/git-cat.go",
		"cmd/git-clean.go",
	}, res)
}

func Test_ParseDiff(t *testing.T) {
	in := fixtures.CleanInput(`
	$ git diff --cached --staged --full-index HEAD
	diff --git a/fixtures/bb.txt b/fixtures/bb.txt
	index 1b28d39c1a2600a86355cd90b25d32e273e91f57..39599d03bfcccc04f209e2bbf74b75b7878b837f 100644
	--- a/fixtures/bb.txt
	+++ b/fixtures/bb.txt
	@@ -1 +1 @@
	-BAR-SHADOW-BLOB 8d52e76479a51b51135c493c56c2ee32f64866af0d518f97e0c3432bc057b0b7
	+BAR-SHADOW-BLOB a554e7d8ecf0c26939167320c04c386f4d19efc74881e905fa5c5934501abeca
	diff --git a/fixtures/egqwert b/fixtures/egqwert
	index 1888497310078a6f2354891fd081f6298f04b1f7..cdf4722b45866f34a35d21a0d16413af617ec863 100644
	--- a/fixtures/egqwert
	+++ b/fixtures/egqwert
	@@ -1 +1 @@
	-BAR-SHADOW-BLOB d0aeb1f7864a7ad42a6527881583b7ef1eae5551aea59cc61bd8083f4653d28d
	+BAR-SHADOW-BLOB a31f5bb02c2bae1438af99b6bd0cb938872197819c8dcc24cb2fd6d740d7868e
	diff --git a/test.txt b/test.txt
	index c0ed5b92e152f49be95950bb4624dadeb372b5a0..868c0795be44e6d304cc73e992d474826f36e070 100644
	--- a/test.txt
	+++ b/test.txt
	@@ -1 +1 @@
	-dsgsfdghdsgsd
	+dsgsfdghdsgsddsgs
	`)
	g, err := git.NewGit("")
	assert.NoError(t, err)

	res, err := g.ParseDiff(in)
	assert.NoError(t, err)
	assert.Equal(t, []git.DiffEntry{
		git.DiffEntry{
			OID:"39599d03bfcccc04f209e2bbf74b75b7878b837f",
			ID:"a554e7d8ecf0c26939167320c04c386f4d19efc74881e905fa5c5934501abeca",
			Filename:"fixtures/bb.txt"},
		git.DiffEntry{
			OID:"cdf4722b45866f34a35d21a0d16413af617ec863",
			ID:"a31f5bb02c2bae1438af99b6bd0cb938872197819c8dcc24cb2fd6d740d7868e",
			Filename:"fixtures/egqwert"}}, res)
}
