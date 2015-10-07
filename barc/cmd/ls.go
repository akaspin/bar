package cmd
import (
	"github.com/akaspin/bar/barc/model"
	"github.com/akaspin/bar/proto/manifest"
	"github.com/akaspin/bar/barc/lists"
	"fmt"
	"github.com/akaspin/bar/barc/transport"
	"text/tabwriter"
	"strings"
	"sort"
)

/*
List bar blobs

	$ barc ls my
	NAME                BLOB    SYNC    GIT     ID              SIZE
	my/blob             no      yes     ok      3d3463063cb1    1930746935
	my/other/blob       yes     no      dirty   309a34901901    3
	my/bad/blob         no      no      dirty   3d3463063cb1    1930746935


*/
type LsCmd struct {
	*BaseSubCommand

	useGit bool
	endpoint string
	pool int

	noHeader bool
	fullID bool

	noBlobs bool
	noManifests bool

	noRemote bool
	noName bool
	noID bool
	noSize bool
}

func NewLsCmd(s *BaseSubCommand) SubCommand {
	c := &LsCmd{BaseSubCommand: s}
	c.FS.BoolVar(&c.useGit, "git", false, "use git infrastructure")
	c.FS.StringVar(&c.endpoint, "endpoint", "http://localhost:3000/v1",
		"bard endpoint")
	c.FS.IntVar(&c.pool, "pool", 16, "pools size")

	c.FS.BoolVar(&c.noHeader, "no-header", false, "do not print header")
	c.FS.BoolVar(&c.fullID, "full-id", false, "print full BLOB IDs")

	c.FS.BoolVar(&c.noBlobs, "no-blobs", false, "do not collect BLOBs")
	c.FS.BoolVar(&c.noManifests, "no-manifests", false,
		"do not collect manifests")

	c.FS.BoolVar(&c.noRemote, "no-remote", false,
		"do not request bard for BLOBs states")
	c.FS.BoolVar(&c.noName, "no-name", false, "do not print BLOB filenames")
	c.FS.BoolVar(&c.noID, "no-id", false, "do not print BLOB IDs")
	c.FS.BoolVar(&c.noSize, "no-size", false, "do not print BLOB sizes")
	return c
}

func (c *LsCmd) Do() (err error) {

	if c.noBlobs && c.noManifests {
		err = fmt.Errorf("both -no-blobs and -no-manifests are on")
		return
	}

	mod, err := model.New(c.WD, c.useGit, manifest.CHUNK_SIZE, c.pool)
	if err != nil {
		return
	}

	feed := lists.NewFileList(c.FS.Args()...).ListDir(c.WD)

	var dirty map[string]struct{}
	if c.useGit {
		if feed, err = mod.Git.FilterByAttr("bar", feed...); err != nil {
			return
		}
		var dirtyst []string
		if dirtyst, err = mod.Git.DiffFiles(feed...); err != nil {
			return
		}
		dirty = map[string]struct{}{}
		for _, n := range dirtyst {
			dirty[n] = struct{}{}
		}
	}

	blobs, err := mod.CollectManifests(!c.noBlobs, !c.noManifests, feed...)
	if err != nil {
		return
	}

	onRemote := map[string]struct{}{}
	if !c.noRemote {
		var exists []string
		trans := transport.NewTransport(c.WD, c.endpoint, c.pool)
		if exists, err = trans.Check(blobs.IDMap().IDs()); err != nil {
			return
		}
		for _, id := range exists {
			onRemote[id] = struct{}{}
		}
	}

	// print this stuff
	w := new(tabwriter.Writer)
	w.Init(c.Stdout, 0, 8, 0, '\t', 0)

	var line []string
	toLine := func(term string) {
		line = append(line, term)
	}
	flushLine := func() {
		fmt.Fprintln(w, strings.Join(line, "\t"))
		line = []string{}
	}

	if !c.noHeader {
		if !c.noName {
			toLine("NAME")
		}
		if !c.noBlobs && !c.noManifests {
			toLine("BLOB")
		}
		if !c.noRemote {
			toLine("SYNC")
		}
		if c.useGit {
			toLine("GIT")
		}
		if !c.noID {
			toLine("ID")
		}
		if !c.noSize {
			toLine("SIZE")
		}
		flushLine()
	}

	var names sort.StringSlice
	for n, _ := range blobs {
		names = append(names, n)
	}
	names.Sort()

	var blobMap map[string]bool
	if !c.noBlobs && !c.noManifests {
		if blobMap, err = mod.IsBlobs(names); err != nil {
			return
		}
	}

	for _, name := range names {
		if !c.noName {
			toLine(name)
		}
		if !c.noBlobs && !c.noManifests {
			if blobMap[name] {
				toLine("yes")
			} else {
				toLine("no")
			}
		}
		if !c.noRemote {
			if _, synced := onRemote[blobs[name].ID]; synced {
				toLine("yes")
			} else {
				toLine("no")
			}
		}
		if c.useGit {
			if _, bad := dirty[name]; !bad {
				toLine("ok")
			} else {
				toLine("dirty")
			}
		}
		if !c.noID {
			if !c.fullID {
				toLine(blobs[name].ID[:12])
			} else {
				toLine(blobs[name].ID)
			}
		}
		if !c.noSize {
			toLine(fmt.Sprintf("%d", blobs[name].Size))
		}

		flushLine()
	}
	w.Flush()

	return
}

