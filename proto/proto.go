package proto
import (
	"github.com/akaspin/bar/proto/manifest"
	"golang.org/x/crypto/sha3"
	"sort"
	"encoding/hex"
)

// Server info
type Info struct {

	// Preferred endpoint
	Endpoint string

	// Preferred chunk size
	ChunkSize int64

	// Preferred number of connections from client
	MaxConn int
}

type BLOBChunk struct {
	BlobID string
	ChunkID string
	Size int64
	Data []byte
}

// Declare commit transaction request
type DeclareUploadTxRequest struct {

	// Commit id
	CommitID string

	// BLOB IDs
	IDs []string
}

// Response for upload request
type DeclareUploadTxResponse struct {

	// Endpoints to upload
	Endpoints []string

	// Missing blob IDs
	MissingIDs []string
}

// Commit upload (not implemented yet)
type CommitUploadTxRequest struct {
	UploadID string
	BindToId string
}

// Download request
type DownloadRequest struct {
	IDs []string
}

// Response to download request
type DownloadResponse struct {

	// Manifests for BLOBs on bard
	BLOBs []manifest.Manifest

	// Mappings between IDs and endpoints
	//
	//    <endpoint>: {id, id ...}
	//
	Endpoints map[string][]string
}

// Tree spec
type Spec struct {

	// Spec ID is SHA3-256 hash of sorted filepath:manifest-id
	ID string

	// File mapping
	BLOBs map[string]string
}

func NewSpec(in map[string]string) (res Spec, err error) {
	hasher := sha3.New256()
	var idBuf []byte

	var names sort.StringSlice
	for n, _ := range in {
		names = append(names, n)
	}
	sort.Sort(names)

	for _, n := range names {
		if _, err = hasher.Write([]byte(n)); err != nil {
			return
		}
		if idBuf, err = hex.DecodeString(in[n]); err != nil {
			return
		}
		if _, err = hasher.Write(idBuf); err != nil {
			return
		}
	}
	res = Spec{hex.EncodeToString(hasher.Sum(nil)), in}
	return
}