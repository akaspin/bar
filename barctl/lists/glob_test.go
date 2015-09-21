package lists_test
import (
	"testing"
	"github.com/akaspin/bar/barctl/lists"
)

func Test_Glob1(t *testing.T) {
	lists.ListFiles("../../", []string{})
}
