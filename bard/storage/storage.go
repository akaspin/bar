package storage
import (
	"io"
	"github.com/akaspin/bar/proto"
	"time"
	"github.com/nu7hatch/gouuid"
)

// All operations in storage driver are independent to each other
type Storage interface {
	io.Closer

	IsSpecExists(id proto.ID) (ok bool, err error)

	IsBLOBExists(id proto.ID) (ok bool, err error)

//	CheckBLOBS(ids []string) (map[string]bool, error)

	// Write spec
	WriteSpec(s proto.Spec) (err error)

	// Read spec
	ReadSpec(id proto.ID) (res proto.Spec, err error)

	// Read proto
	ReadManifest(id proto.ID) (res *proto.Manifest, err error)

	// Get manifests by it's ids
	GetManifests(ids []proto.ID) (res []proto.Manifest, err error)

	// Declare new upload
	DeclareUpload(m proto.Manifest) (err error)

	// Write chunk for declared blob from given reader
	WriteChunk(blobID, chunkID proto.ID, size int64, r io.Reader) (err error)

	// Finish upload
	FinishUpload(id proto.ID) (err error)

	// Read Chunk from blob by size and offset
	ReadChunkFromBlob(blobID proto.ID, size, offset int64, w io.Writer) (err error)

	// Returns IDs of requested blobs except already stored
	// GetMissingBlobIDs(ids []proto.ID) (res []proto.ID, err error)

	//// New store API

	// Create new upload session. Returns upload ID and IDs of chunks
	// missing on bard. Upload id is simple uuid bytes
	CreateUploadSession(uploadID uuid.UUID, in []proto.Manifest, ttl time.Duration) (missingChunkIDs []proto.ID, err error)

	// Upload chunk
	UploadChunk(uploadID uuid.UUID, chunkID proto.ID, r io.Reader) (err error)

	// Finish upload session
	FinishUploadSession(uploadID uuid.UUID) (err error)
}


