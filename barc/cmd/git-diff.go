package cmd
import (
	"flag"
	"io"
	"github.com/tamtam-im/logx"
)


/*
Git diff command. Register it in config:

	# .git/config
	[diff "bar"]
		command = barc git-diff ...

And invoke like:

	$ git diff --staged --cached --full-index

Git provides following arguments like:

	fix/fix.txt /var/folders/nx/vpl22rw925jgp5fwtczb16s40000gn/T//7aQoZ8_fix.txt 4b7cfc80527d8ab6fffa8222f090c101b42b362f 100644 fix/fix.txt 8f9e07de59a92d61ed7bc942df8717f9f304092d 100644

or

	fix/fix1.txt /dev/null . . fix/fix1.txt 8f9e07de59a92d61ed7bc942df8717f9f304092d 100644

where

	fix/fix.txt                             Filename in working tree

	/var/..._fix.txt (or /dev/null)         Temporary filename from git index
	4b7cfc80527d8... (or .)                 OID from git index
	100644 (or .)                           Git file mode

	fix/fix.txt (or /dev/null)              Filename in staging area.
											If this value differd - file is
											dirty
	8f9e07de59a92... (or .)
	100644 (or .)

Both objects must be manifests. Diff parses both and emits usual git diff like

	diff --git a/fix2.txt b/fix2.txt
	deleted file mode 100644
	index e56e15bb7ddb6bd0b6d924b18fcee53d8713d7ea..0000000000000000000000000000000000000000
	--- a/fix2.txt
	+++ /dev/null
	@@ -1 +0,0 @@
	-BAR:BLOB 859a7a7603028deeb3b66234cffa5191466d1a0538e449a19812273b0d98dc1c

or

	diff --git a/fix.txt b/fix.txt
	index 190a18037c64c43e6b11489df4bf0b9eb6d2c9bf..e56e15bb7ddb6bd0b6d924b18fcee53d8713d7ea 100644
	--- a/fix.txt
	+++ b/fix.txt
	@@ -1 +1 @@
	-BAR:BLOB ...
	+BAR:BLOB ...

or

	diff --git a/fix1.txt b/fix1.txt
	new file mode 100644
	index 0000000000000000000000000000000000000000..e56e15bb7ddb6bd0b6d924b18fcee53d8713d7ea
	--- /dev/null
	+++ b/fix1.txt
	@@ -0,0 +1 @@
	+BAR:BLOB ...
*/
type GitDiffCmd struct {
	fs *flag.FlagSet
}

func (c *GitDiffCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.fs = fs
	return
}

func (c GitDiffCmd) Do() (err error) {
	logx.Debug("diff-cmd", c.fs.Args())
	return
}
