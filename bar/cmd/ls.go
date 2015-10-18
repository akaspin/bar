package cmd
import (
	"github.com/akaspin/bar/bar/model"
	"github.com/akaspin/bar/proto"
	"github.com/akaspin/bar/bar/lists"
	"fmt"
	"github.com/akaspin/bar/bar/transport"
	"text/tabwriter"
	"strings"
	"sort"
	"flag"
)

/*
List bar blobs

	$ bar ls my
	NAME                BLOB    SYNC    GIT     ID              SIZE
	my/blob             no      yes     ok      3d3463063cb1    1930746935
	my/other/blob       yes     no      dirty   309a34901901    3
	my/bad/blob         no      no      dirty   3d3463063cb1    1930746935


*/
type LsCmd struct {
	*Base
	useGit bool

	noHeader bool
	fullID bool

	noBlobs bool
	noManifests bool

	noRemote bool
	noName bool
	noID bool
	noSize bool
}

func NewLsCmd(s *Base) SubCommand {
	c := &LsCmd{Base: s}

	return c
}

func (c *LsCmd) Init(fs *flag.FlagSet) {
	fs.BoolVar(&c.useGit, "git", false, "use git infrastructure")

	fs.BoolVar(&c.noHeader, "no-header", false, "do not print header")
	fs.BoolVar(&c.fullID, "full-id", false, "print full BLOB IDs")

	fs.BoolVar(&c.noBlobs, "no-blobs", false, "do not collect BLOBs")
	fs.BoolVar(&c.noManifests, "no-manifests", false,
		"do not collect manifests")

	fs.BoolVar(&c.noRemote, "no-remote", false,
		"do not request bard for BLOBs states")
	fs.BoolVar(&c.noName, "no-name", false, "do not print BLOB filenames")
	fs.BoolVar(&c.noID, "no-id", false, "do not print BLOB IDs")
	fs.BoolVar(&c.noSize, "no-size", false, "do not print BLOB sizes")

}

func (c *LsCmd) Do(args []string) (err error) {

	if c.noBlobs && c.noManifests {
		err = fmt.Errorf("both -no-blobs and -no-manifests are on")
		return
	}

	mod, err := model.New(c.WD, c.useGit, proto.CHUNK_SIZE, c.PoolSize)
	if err != nil {
		return
	}

	feed := lists.NewFileList(args...).ListDir(c.WD)

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

	blobs, err := mod.FeedManifests(!c.noBlobs, !c.noManifests, true, feed...)
	if err != nil {
		return
	}

	missingOnRemote := map[proto.ID]struct{}{}
	if !c.noRemote {
		var exists []proto.ID
		trans := transport.NewTransport(mod, "", c.Endpoints, c.PoolSize)
		if exists, err = trans.Check(blobs.IDMap().IDs()); err != nil {
			return
		}
		for _, id := range exists {
			missingOnRemote[id] = struct{}{}
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
		if blobMap, err = mod.IsBlobs(names...); err != nil {
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
			if _, missing := missingOnRemote[blobs[name].ID]; missing {
				toLine("no")
			} else {
				toLine("yes")
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
				toLine(blobs[name].ID.String()[:12])
			} else {
				toLine(blobs[name].ID.String())
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

func (c *LsCmd) Description() string {
	return "show information about bar-tracked blobs"
}
