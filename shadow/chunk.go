package shadow
import "fmt"

type Chunk struct  {
	ID string
	Size int64
	Offset int64
}

func (c Chunk) String() string {
	return fmt.Sprintf("id %s\nsize %d\noffset %d\n\n", c.ID, c.Size, c.Offset)
}

// Guess chunk size by BLOB size. For now 1M
func GuessChunkSize(size int64) (res int64) {
	return CHUNK_SIZE
}
