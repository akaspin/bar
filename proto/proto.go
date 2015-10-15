package proto


type ChunkInfo struct {
	BlobID ID
	Chunk
}

type ChunkData struct {
	ChunkInfo
	Data []byte
}

