package cmd
import (
	"flag"
	"io"
	"fmt"
	"github.com/akaspin/bar/barc/git"
	"github.com/akaspin/bar/shadow"
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
	chunkSize int64

	wd string
	fs *flag.FlagSet
	out io.Writer

	git *git.Git
}

func (c *GitDiffCmd) Bind(wd string, fs *flag.FlagSet, in io.Reader, out io.Writer) (err error) {
	c.wd = wd
	c.fs = fs
	c.out = out

	fs.Int64Var(&c.chunkSize, "chunk", shadow.CHUNK_SIZE, "preferred chunk size")
	return
}

func (c *GitDiffCmd) Do() (err error) {
	if c.git, err = git.NewGit(c.wd); err != nil {
		return
	}

	wtName := c.fs.Arg(0)

//	lName := c.fs.Arg(1)
	lOID := c.fs.Arg(2)
	lMode := c.fs.Arg(3)

//	rName := c.fs.Arg(4)
	rOID := c.fs.Arg(5)
	rMode := c.fs.Arg(6)

	var isNew, isDeleted bool
	if rOID == "." {
		isDeleted = true
	}
	if lOID == "." {
		isNew = true
	}

	fmt.Fprintf(c.out, "diff --git a/%s b/%s\n", wtName, wtName)

	if isDeleted {
		fmt.Fprintf(c.out, "deleted file mode %s\n", lMode)
	} else if isNew {
		fmt.Fprintf(c.out, "new file mode %s\n", rMode)
	}

	fmt.Fprintf(c.out, "index %s..%s\n", lOID, rOID)

	if isNew {
		fmt.Fprintf(c.out, "--- /dev/null\n")
	} else {
		fmt.Fprintf(c.out, "--- a/%s\n", wtName)
	}

	if isDeleted {
		fmt.Fprintf(c.out, "+++ /dev/null\n")
	} else {
		fmt.Fprintf(c.out, "+++ b/%s\n", wtName)
	}

	if isNew {
		fmt.Fprintln(c.out, "@@ -0,0 +1 @@")
	} else if isDeleted {
		fmt.Fprintln(c.out, "@@ -1 +0,0 @@")
	} else {
		fmt.Fprintln(c.out, "@@ -1 +1 @@")
	}

	var r io.Reader
	var m *shadow.Shadow
	if !isNew {
		if r, err = c.git.Cat(lOID); err != nil {
			return
		}
		if m, err = shadow.NewFromAny(r, c.chunkSize); err != nil {
			return
		}
		logx.Debugf("manifest from %s (source is manifest: %t)", lOID, m.IsFromShadow)
		fmt.Fprintf(c.out, "-BAR-SHADOW-BLOB %s\n", m.ID)
	}

	if !isDeleted {
		if r, err = c.git.Cat(rOID); err != nil {
			return
		}
		if m, err = shadow.NewFromAny(r, c.chunkSize); err != nil {
			return
		}
		logx.Debugf("manifest from %s (source is manifest: %t)", rOID, m.IsFromShadow)
		fmt.Fprintf(c.out, "+BAR-SHADOW-BLOB %s\n", m.ID)
	}

	return
}
