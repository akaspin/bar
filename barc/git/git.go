package git
import (
	"strings"
	"bufio"
	"bytes"
	"io"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"github.com/akaspin/bar/barc/lists"
	"github.com/akaspin/bar/proto"
	"github.com/tamtam-im/logx"
)

// Extracted diff entry
type DiffEntry struct {

	// Git OID to get cached manifest
	OID string

	// bar BLOB ID
	ID string

	// BLOB filename
	Filename string
}


// Git wrapper.
// TODO: use native git (https://github.com/gogits/git)
//
// This wrapper always run git in repository root
type Git struct {

	*lists.Mapper
}

// New git in given root.
func NewGit(cwd string) (res *Git, err error) {
	if cwd == "" {
		if cwd, err = os.Getwd(); err != nil {
			return
		}
	}

	raw, err := execCommand("git", "-C", cwd,
		"rev-parse", "--show-toplevel").Output()
	if err != nil {
		return
	}
	root := strings.TrimSpace(string(raw))



	res = &Git{lists.NewMapper(cwd, root)}
	return
}

// Run git command
func (g *Git) Run(sub string, arg ...string) (res string, err error) {
	logx.Tracef("executing `git %s`", strings.Join(
		append([]string{"-C", g.Root, "--no-pager", sub}, arg...), " "))

	c := execCommand("git",
		append([]string{"-C", g.Root, "--no-pager", sub}, arg...)...)

	var out, stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr
	err = c.Run()
	if err != nil {
		err = fmt.Errorf("%s %s", err, stderr.String())
	}
	res = out.String()
	return
}

// Refresh files in git index (use after squash or blow)
func (g *Git) UpdateIndex(what ...string) (err error) {
	rooted, err := g.ToRoot(what...)
	if _, err = g.Run("update-index", rooted...); err != nil {
		return
	}
	logx.Debugf("git index updated for %s", what)
	return
}

// Set config value
func (g *Git) SetConfig(key, val string) (err error) {
	_, err = g.Run("config", "--local", key, val)
	return
}

// Unset config value
func (g *Git) UnsetConfig(key string) (err error)  {
	_, err = g.Run("config", "--local", "--unset", key)
	return
}

// Get git hook contents
func (g *Git) GetHook(name string) (res string, err error) {
	f, err := os.Open(g.hookName(name))
	if err != nil {
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	res = string(data)
	return
}

// Install git hook
func (g *Git) SetHook(name string, contents string) (err error) {
	f, err := os.OpenFile(g.hookName(name),
		os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0776)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	return
}

func (g *Git) CleanHook(name string) (err error) {
	err = os.Remove(g.hookName(name))
	return
}

func (g *Git) hookName(name string) (res string) {
	return filepath.Join(g.Root, ".git", "hooks", name)
}

// Returns dirty files with filter=bar
func (g *Git) DiffFilesWithAttr(arg ...string) (res []string, err error) {
	delta, err := g.DiffFiles(arg...)
	if err != nil {
		return
	}
	if len(delta) == 0 {
		return
	}

	res, err = g.FilterByAttr("bar", delta...)
	if err != nil {
		return
	}
	res, err = g.FromRoot(res...)
	return
}

// Get list of non-staged files in working tree
//
//    $ git diff-files --name-only
//
// This command always takes filenames relative to CWD
func (g *Git) DiffFiles(what ...string) (res []string, err error) {
	if what, err = g.ToRoot(what...); err != nil {
		return
	}

	rawFiles, err := g.Run("diff-files",
		append([]string{"--name-only", "-z"}, what...)...)
	if err != nil {
		return
	}
	for _, f := range strings.Split(rawFiles, "\x00") {
		if f != "" {
			res = append(res, f)
		}
	}
	res, err = g.ToRoot(res...)
	return
}

// Filter files by diff attribute
//
//    $ git check-attr diff <files> | grep "diff: <diff>"
//
// This command takes and returns filenames relative to CWD
func (g *Git) FilterByAttr(diff string, filenames ...string) (res []string, err error) {
	if filenames, err = g.ToRoot(filenames...); err != nil {
		return
	}

	rawAttrs, err := g.Run("check-attr",
		append([]string{"filter"}, filenames...)...)
	if err != nil {
		return
	}

	attrReader := bufio.NewReader(strings.NewReader(rawAttrs))
	var data []byte
	var line string
	suffix := ": filter: " + diff
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
	res, err = g.FromRoot(res...)
	return
}

// Get file OID in Git index
//
//    $ git ls-files --cached -s --full-name <file>
//

func (g *Git) GetOID(filename string) (res string, err error) {
	rooted, err := g.ToRoot(filename)
	if err != nil {
		return
	}
	raw, err := g.Run("ls-files", "--cached", "-s", "--full-name", "-z", rooted[0])
	if err != nil {
		return
	}
	var trash1, trash2 string
	_, err = fmt.Sscanf(raw, "%s %s %s", &trash1, &res, &trash2)
	return
}

// Cat file from index
//
//    $ git cat-file -p <OID>
//
func (g *Git) Cat(oid string) (res io.Reader, err error) {
	raw, err := g.Run("cat-file", "-p", oid)
	if err != nil {
		return
	}
	res = strings.NewReader(raw)
	return
}

// Run diff against HEAD
//
//    $ git diff --cached --staged --full-index --no-color HEAD
//
func (g *Git) Diff() (res io.Reader, err error) {
	raw, err := g.Run("diff",
		"--cached", "--staged", "--full-index", "--no-color",
		"-U99999999999999", "--no-prefix")
	if err != nil {
		return
	}
	res = strings.NewReader(raw)
	return
}

/*
Extract manifests from diff

	$ git diff --staged --cached -U99999999999999 --no-prefix
	diff --git fix/aa.txt fix/aa.txt
	index 5afd5028d71cfadf73c0e3abd70f852d67357909..63c202a03152c8635eb78dcfef35859f0c68f5cf 100644
	--- fix/aa.txt
	+++ fix/aa.txt
	@@ -1,10 +1,10 @@
	 BAR:MANIFEST

	-id aba7aeb8a7948dd0cdb8eeb9239e5d1dab2bd840f13930f86f6e67ba40ea5350
	-size 4
	+id f627c8f9355399ef45e1a6b6e5a9e6a3abcb3e1b6255603357bffa9f2211ba7e
	+size 6


	-id aba7aeb8a7948dd0cdb8eeb9239e5d1dab2bd840f13930f86f6e67ba40ea5350
	-size 4
	+id f627c8f9355399ef45e1a6b6e5a9e6a3abcb3e1b6255603357bffa9f2211ba7e
	+size 6
	 offset 0

		...

	diff --git fix/aa2.txt fix/aa2.txt
	new file mode 100644
	index 0000000000000000000000000000000000000000..63c202a03152c8635eb78dcfef35859f0c68f5cf
	--- /dev/null
	+++ fix/aa2.txt         <--- Start consuming
	@@ -0,0 +1,10 @@        <--- if consuming - next line
	+BAR:MANIFEST           <--- Start consuming manifest by "" and "+"
	+
	+id f627c8f9355399ef45e1a6b6e5a9e6a3abcb3e1b6255603357bffa9f2211ba7e
	+size 6
	+
	+
	+id f627c8f9355399ef45e1a6b6e5a9e6a3abcb3e1b6255603357bffa9f2211ba7e
	+size 6
	+offset 0
	+


*/
func (g *Git) ManifestsFromDiff(r io.Reader) (res lists.BlobMap, err error) {
	res = lists.BlobMap{}

	var data []byte
	var line string
	buf := bufio.NewReader(r)

	var filename string
	var w *bytes.Buffer
	var hunt, eof bool

	for {
		data, _, err = buf.ReadLine()
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return
		} else {
			line = strings.TrimSpace(string(data))
		}

		if eof || strings.HasPrefix(line, "diff --git ") {
			if hunt {
				mn, err := proto.NewFromManifest(bytes.NewReader(w.Bytes()))
				if err != nil {
					return res, err
				}
				// need to lift filenames from root
				f1, err := g.FromRoot(filename)
				if err != nil {
					return res, err
				}
				res[f1[0]] = *mn
			}
			filename = ""
			w = nil
			hunt = false
			if eof {
				return
			}
			continue
		}

		if strings.HasPrefix(line, "+++ ") {
			// new file - hunt
			filename = strings.TrimPrefix(line, "+++ ")
			continue
		}

		if strings.TrimPrefix(line, "+") == "BAR:MANIFEST" {
			// HUNT
			hunt = true
			w = new(bytes.Buffer)
			if _, err = fmt.Fprintln(w, "BAR:MANIFEST"); err != nil {
				return
			}
			continue
		}

		if hunt {
			line := strings.TrimPrefix(line, "+")
			if !strings.HasPrefix(line, "-") {
				fmt.Fprintln(w, line)
			}
		}
	}

	return
}

// Get BLOBs to upload for pre-commit hook.
//
//    diff --git a/fixtures/bb.txt b/fixtures/bb.txt
//    index 1b28d39c1a2600a86355cd90b25d32e273e91f57..39599d03bfcccc04f209e2bbf74b75b7878b837f 100644
//    --- a/fixtures/bb.txt
//    +++ b/fixtures/bb.txt
//    @@ -1 +1 @@
//    -BAR-SHADOW-BLOB 8d52e76479a51b51135c493c56c2ee32f64866af0d518f97e0c3432bc057b0b7
//    +BAR-SHADOW-BLOB a554e7d8ecf0c26939167320c04c386f4d19efc74881e905fa5c5934501abeca
//
// where:
//
//    diff --git <skip>                 <- detect new file
//    index <skip>...<OID> <skip>       <- extract OID
//    <skip>
//    +++ b/<Filename>                  <- extract filename
//    <skip>
//    +BAR-SHADOW-BLOB <ID>             <- Assume as BLOB and extract ID
//
func (g *Git) ParseDiff(r io.Reader) (res []DiffEntry, err error) {
	var data []byte
	var line string
	buf := bufio.NewReader(r)

	var oid, id, filename string
	for {
		data, _, err = buf.ReadLine()
		if err == io.EOF {
			err = nil
			return
		} else if err != nil {
			return
		}
		line = strings.TrimSpace(string(data))

		if strings.HasPrefix(line, "diff --git ") {
			// New file
			oid, id, filename = "", "", ""
			continue
		}

		if strings.HasPrefix(line, "index") {
			oid = line[48:88]
			continue
		}

		if strings.HasPrefix(line, "+++ b/") {
			filename = strings.TrimPrefix(line, "+++ b/")
			continue
		}

		if strings.HasPrefix(line, "+BAR-SHADOW-BLOB ") {
			id = strings.TrimPrefix(line, "+BAR-SHADOW-BLOB ")
			res = append(res, DiffEntry{oid, id, filename})
			continue
		}
	}
	return
}

