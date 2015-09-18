package fixtures
import (
	"encoding/hex"
	"path/filepath"
)

func StoredName(root string, id []byte) string {
	s := hex.EncodeToString(id)
	return filepath.Join(root, s[:2], s)
}
