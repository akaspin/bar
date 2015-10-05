package proto
import "github.com/akaspin/bar/proto/manifest"

// Server info
type Info struct {

	// Preferred chunk size
	ChunkSize int64

	// Preferred number of connections from client
	MaxConn int
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

// Tree spec for git-less usage
type Spec struct {

	// Spec url.
	URL string

	//
	BLOBs map[string]manifest.Manifest
}