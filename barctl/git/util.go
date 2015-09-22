package git
import (
	"strings"
	"os/exec"
)


func execGit(root, sub string, arg ...string) *exec.Cmd {
	var err error
	if root == "" {
		root, err = GetGitTop()
		if err != nil {
			panic(err)
		}
	}
	c := execCommand("git", append([]string{"-C", root, sub}, arg...)...)
	return c
}


// Get git top directory
//
//    git rev-parse --show-toplevel
//
func GetGitTop() (res string, err error) {
	raw, err := execCommand("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return
	}
	res = strings.TrimSpace(string(raw))
	return
}
