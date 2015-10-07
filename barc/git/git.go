	"github.com/akaspin/bar/proto/manifest"
	"github.com/tamtam-im/logx"
	logx.Tracef("executing `git %s`", strings.Join(
		append([]string{"-C", g.Root, "--no-pager", sub}, arg...), " "))

	rooted, err := g.ToRoot(what...)
	if _, err = g.Run("update-index", rooted...); err != nil {
		return
	}
	logx.Debugf("git index updated for %s", what)
		"--cached", "--staged", "--full-index", "--no-color",
		"-U99999999999999", "--no-prefix")
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
func (g *Git) ManifestsFromDiff(r io.Reader) (res lists.Links, err error) {
	res = lists.Links{}

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
				mn, err := manifest.NewFromManifest(bytes.NewReader(w.Bytes()))
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
