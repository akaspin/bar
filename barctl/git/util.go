package git
import "strings"

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
