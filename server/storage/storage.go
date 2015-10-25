package storage
import (
	"io"
	"github.com/akaspin/bar/proto"
	"time"
	"github.com/nu7hatch/gouuid"
)

// All operations in storage driver are independent to each other
type Storage interface {

	IsSpecExists(id proto.ID) (ok bool, err error)

	// Write spec
	WriteSpec(s proto.Spec) (err error)

	// Read spec
	ReadSpec(id proto.ID) (res proto.Spec, err error)

	// Get manifests by it's ids
	GetManifests(ids []proto.ID) (res []proto.Manifest, err error)

	// Read Chunk from blob by size and offset
	ReadChunkFromBlob(blobID proto.ID, size, offset int64, w io.Writer) (err error)

	// Returns IDs of requested blobs except already stored
	GetMissingBlobIDs(ids []proto.ID) (res []proto.ID, err error)

	// Create new upload session. Returns upload ID and IDs of chunks
	// missing on bard. Upload id is simple uuid bytes
	CreateUploadSession(uploadID uuid.UUID, in []proto.Manifest, ttl time.Duration) (missingChunkIDs []proto.ID, err error)

	// Upload chunk
	UploadChunk(uploadID uuid.UUID, chunkID proto.ID, r io.Reader) (err error)

	// Finish upload session
	FinishUploadSession(uploadID uuid.UUID) (err error)
}


