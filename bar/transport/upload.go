package transport
import (
	"github.com/nu7hatch/gouuid"
	"github.com/akaspin/bar/bar/lists"
	"time"
	"github.com/akaspin/bar/proto"
)

type Upload struct  {
	*Transport
	*uuid.UUID
	ttl time.Duration
}

func NewUpload(t *Transport, ttl time.Duration) (res *Upload) {
	id, _ := uuid.NewV4()
	res = &Upload{t, id, ttl}
	return
}

func (u *Upload) SendCreateUpload(links lists.BlobMap) (missing []proto.ID, err error) {
	cli, err := u.Transport.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	mans := proto.ManifestSlice(links.GetManifestSlice())
	tSlice, err := mans.MarshalThrift()
	if err != nil {
		return
	}

	res1, err := cli.CreateUpload((*u.UUID)[:], tSlice, int64(u.ttl))
	if err != nil {
		return
	}
    var r2 proto.IDSlice
	err = (&r2).UnmarshalThrift(res1)
	missing = r2
	return
}

func (u *Upload) UploadChunk(name string, chunk proto.Chunk) (err error) {
	buf := make([]byte, chunk.Size)
	if err = u.Transport.model.ReadChunk(name, chunk, buf); err != nil {
		return
	}

	cli, err := u.Transport.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()

	id, err := chunk.ID.MarshalThrift()
	if err != nil {
		return
	}

	err = cli.UploadChunk((*u.UUID)[:], id, buf)
	return
}

func (u *Upload) Commit() (err error) {
	cli, err := u.Transport.tPool.Take()
	if err != nil {
		return
	}
	defer cli.Release()
	err = cli.FinishUpload((*u.UUID)[:])
	return
}