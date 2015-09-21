package git_test
import (
	"testing"
	"github.com/akaspin/bar/fixtures"
	"github.com/akaspin/bar/barctl/git"
	"strings"
	"github.com/stretchr/testify/assert"
)

func Test_Diff1(t *testing.T) {
	src := fixtures.CleanManifest(`
		diff --git a/.gitattributes b/.gitattributes
		index 6a54f64..cb16f98 100644
		--- a/.gitattributes
		+++ b/.gitattributes
		@@ -1,2 +1,3 @@
		-/fixtures/*    filter=bar
		-/fixtures/*    diff=bar
		+[barrel]filter=bar diff=bar
		+
		+/fixtures/*    filter=bar diff=bar
		diff --git a/fixtures/aa.txt b/fixtures/aa.txt
		deleted file mode 100644
		index 57a86b4..0000000
		--- a/fixtures/aa.txt
		+++ /dev/null
		@@ -1 +0,0 @@
		-BAR-SHADOW-BLOB 04f4efe30dd589f5e3102c1e1ecbdb07846c6338b37301daf4348e5b08f26a06
		diff --git a/fixtures/bb.txt b/fixtures/bb.txt
		new file mode 100644
		index 0000000..25eeb55
		--- /dev/null
		+++ b/fixtures/bb.txt
		@@ -0,0 +1 @@
		+BAR-SHADOW-BLOB 845b6321fec6848de5bf28d0312f19e9faee872ee68a308b7aa935a59c601f78
		diff --git a/test.txt b/test.txt
		index df9612f..c0ed5b9 100644
		--- a/test.txt
		+++ b/test.txt
		@@ -1 +1 @@
		-dsgsfdgh
		+dsgsfdghdsgsd
	`)
	res, err := git.ParseCommitDiff(strings.NewReader(src))
	assert.NoError(t, err)
	assert.EqualValues(t, res, map[string]string{
		"845b6321fec6848de5bf28d0312f19e9faee872ee68a308b7aa935a59c601f78": "fixtures/bb.txt",
	})
}
