package cmd_test
import (
	"testing"
	"os"
)

func skip(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skip integration")
	}
}
