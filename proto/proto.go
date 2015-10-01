// bard proto stubs
package proto

// Server info
type Info struct {

	// Preferred chunk size
	ChunkSize int64

	// Preferred number of connections from client
	MaxConn int
}

// Declare commit transaction request
type DeclareCommitTxRequest struct {

	// Commit id
	CommitID string

	// BLOB IDs
	IDs []string
}

type DeclareCommitTxResponse struct {
	MissingIDs []string
}
