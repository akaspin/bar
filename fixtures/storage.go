package fixtures
import (
	"path/filepath"
)

func StoredName(root string, id string) string {
	return filepath.Join(root, id[:2], id)
}
