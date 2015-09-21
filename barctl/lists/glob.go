package lists
import (
	"strings"
	"path/filepath"
	"os"
	"fmt"
)

// Get files
func ListFiles(root string, what []string) (res []string, err error) {
	var globs []string
	var paths []string

	for _, w := range what {
		if strings.ContainsAny(w, "*?[]^") {
			globs = append(globs, w)
		} else {
			paths = append(paths, w)
		}
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		fmt.Println(filepath.Clean(path))
		return nil
	})

	return
}

