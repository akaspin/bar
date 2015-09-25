package git
import (
	"os/exec"
	"strings"
	"bufio"
	"bytes"
	"io"
	"fmt"
)

type CommitBLOB struct {

	// Git OID to get cached manifest
	OID string

	// bar BLOB ID
	ID string

	// BLOB filename
	Filename string
}

// Get git top directory
//
//    git rev-parse --show-toplevel
//
func getGitTop() (res string, err error) {
	raw, err := execCommand("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return
	}
	res = strings.TrimSpace(string(raw))
	return
}

// Git wrapper
type Git struct {
	Root string
}

// New git. Panic on error
func NewGit(root string) (res *Git, err error) {
	if root == "" {
		if root, err = getGitTop(); err != nil {
			return
		}
	}
	res = &Git{root}
	return
}

// Run git command
func (g *Git) Run(sub string, arg ...string) *exec.Cmd {
	// exec.Cmd.Dir has no effect for git.
	return execCommand("git",
		append([]string{"-C", g.Root, "--no-pager", sub}, arg...)...)
}

// Get list of non-staged files in working tree
//
//    $ git diff-files --name-only
//
func (g *Git) DirtyFiles() (res []string, err error) {
	rawFiles, err := g.Run("diff-files", "--name-only", "-z").Output()
	if err != nil {
		return
	}
	for _, f := range strings.Split(string(rawFiles), "\x00") {
		if f != "" {
			res = append(res, f)
		}
	}
	return
}

// Filter files by diff attribute
//
//    $ git check-attr diff <files> | grep "diff: <diff>"
//
func (g *Git) FilterByDiff(diff string, filenames ...string) (res []string, err error) {
	rawAttrs, err := g.Run("check-attr",
		append([]string{"diff"}, filenames...)...).Output()

	attrReader := bufio.NewReader(bytes.NewReader(rawAttrs))
	var data []byte
	var line string
	suffix := ": diff: " + diff
	for {
		data, _, err = attrReader.ReadLine()
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}
		line = string(data)
		if !strings.HasPrefix(line, ":") && strings.HasSuffix(line, suffix) {
			res = append(res, strings.TrimSuffix(line, suffix))
		}
	}
	return
}

// Get file OID in Git index
//
//    $ git ls-files --cached -s --full-name <file>
//
func (g *Git) OID(filename string) (res string, err error) {
	c := g.Run("ls-files", "--cached", "-s", "--full-name", "-z", filename)
	raw, err := c.Output()
	if err != nil {
		return
	}
	var trash1, trash2 string
	_, err = fmt.Sscanf(string(raw), "%s %s %s", &trash1, &res, &trash2)
	return
}

// Cat file from index
//
//    $ git cat-file -p <OID>
//
func (g *Git) Cat(oid string) (res io.Reader, err error) {
	c := execCommand("git", "cat-file", "-p", oid)
	raw, err := c.Output()
	if err != nil {
		return
	}
	res = bytes.NewReader(raw)
	return
}

// Run diff against HEAD
//
//    $ git diff --cached --staged --full-index --no-color HEAD
//
func (g *Git) Diff() (res io.Reader, err error) {
	raw, err := g.Run("diff",
		"--cached", "--staged", "--full-index", "--no-color", "HEAD").Output()
	if err != nil {
		return
	}
	res = bytes.NewReader(raw)
	return
}

// Get BLOBs to upload for pre-commit hook.
//
//    diff --git a/fixtures/bb.txt b/fixtures/bb.txt
//    index 1b28d39c1a2600a86355cd90b25d32e273e91f57..39599d03bfcccc04f209e2bbf74b75b7878b837f 100644
//    --- a/fixtures/bb.txt
//    +++ b/fixtures/bb.txt
//    @@ -1 +1 @@
//    -BAR-SHADOW-BLOB 8d52e76479a51b51135c493c56c2ee32f64866af0d518f97e0c3432bc057b0b7
//    +BAR-SHADOW-BLOB a554e7d8ecf0c26939167320c04c386f4d19efc74881e905fa5c5934501abeca
//
// where:
//
//    diff --git <skip>                 <- detect new file
//    index <skip>...<OID> <skip>       <- extract OID
//    <skip>
//    +++ b/<Filename>                  <- extract filename
//    <skip>
//    +BAR-SHADOW-BLOB <ID>             <- Assume as BLOB and extract ID
//
func (g *Git) ParseDiff(r io.Reader) (res []CommitBLOB, err error) {
	var data []byte
	var line string
	buf := bufio.NewReader(r)

	var oid, id, filename string
	for {
		data, _, err = buf.ReadLine()
		if err == io.EOF {
			err = nil
			return
		} else if err != nil {
			return
		}
		line = strings.TrimSpace(string(data))

		if strings.HasPrefix(line, "diff --git ") {
			// New file
			oid, id, filename = "", "", ""
			continue
		}

		if strings.HasPrefix(line, "index") {
			oid = line[48:88]
			continue
		}

		if strings.HasPrefix(line, "+++ b/") {
			filename = strings.TrimPrefix(line, "+++ b/")
			continue
		}

		if strings.HasPrefix(line, "+BAR-SHADOW-BLOB ") {
			id = strings.TrimPrefix(line, "+BAR-SHADOW-BLOB ")
			res = append(res, CommitBLOB{oid, id, filename})
			continue
		}
	}
	return
}
