package git
import (
	"fmt"
	"io"
	"bytes"
	"strings"
	"bufio"
)

func GetDirtyFiles(root string) (res []string, err error) {
	rawFiles, err := execGit(root, "diff-files", "--name-only", "-z").Output()
	if err != nil {
		return
	}
	res = strings.Split(string(rawFiles), "\x00")
	return
}


func IsBLOBsClean(root, detect string, files []string) (ok bool, err error) {
	rawAttrs, err := execGit(root, "check-attr",
		append([]string{"diff"}, files...)...).Output()

	attrReader := bufio.NewReader(bytes.NewReader(rawAttrs))
	var data []byte
	var trash string
	for {
		data, _, err = attrReader.ReadLine()
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}
		trash = string(data)
		if !strings.HasPrefix(trash, ":") && strings.HasSuffix(trash, "diff: " + detect) {
			return
		}
	}
	ok = true
	return
}

func IsClean(root, filename string) (res bool, err error) {
	c := execGit(root, "diff-files", "--name-only", filename)

	raw, err := c.Output()
	if err != nil {
		return
	}
	res = len(raw) == 0
	return
}

func GetFileOID(root, filename string) (res string, err error) {
	c := execGit(root, "ls-files", "--cached", "-s", "--full-name", "-z", filename)
	raw, err := c.Output()
	if err != nil {
		return
	}
	var trash1, trash2 string
	_, err = fmt.Sscanf(string(raw), "%s %s %s", &trash1, &res, &trash2)
	return
}

func CatFile(oid string) (res io.Reader, err error) {
	c := execCommand("git", "cat-file", "-p", oid)
	raw, err := c.Output()
	if err != nil {
		return
	}
	res = bytes.NewReader(raw)
	return
}

