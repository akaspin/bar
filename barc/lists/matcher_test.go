package lists_test
import (
	"testing"
	"regexp"
	"github.com/stretchr/testify/assert"
)

func Test_Match1(t *testing.T) {
	fix := "my/weird file/long/path.txt"

	pattens := map[string]bool{
		"^my/weird file/long/path.txt$": true,  // exact match
		"^my/.+/path.txt$": true,               // my/**/path.txt
		"^my/.+/path$": false,                  // my/**/path
		"^my/[^/]+/path.txt$": false,           // my/*/path.txt
		"^my/[^/]+/[^/]+/path.txt$": true,      // my/*/*/path.txt
	}

	for p, er := range pattens {
		res, err := regexp.MatchString(p, fix)
		assert.NoError(t, err)
		assert.Equal(t, er, res)
	}
}

func Test_Match_Dotfiles(t *testing.T) {
	dots := []string{
		".git",
		"any/long/.dotpath/file.txt",
		"path/.dotfile",
		".dotpath/file.txt",
	}

	for _, d := range dots {
		res, err := regexp.MatchString("^.*\\..*$", d)
		assert.NoError(t, err)
		assert.True(t, res, d)
	}
}