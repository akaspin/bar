package git
import (
	"os/exec"
	"strings"
	"bufio"
	"bytes"
	"io"
	"fmt"
)


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
