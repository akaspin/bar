package manifest
import "fmt"

const CHUNK_SIZE = 1024 * 1024 * 2

// Manifest chunk
type Chunk struct  {

	// Chunk ID
	ID string

	// Chunk Size
	Size int64

	// Offset
	Offset int64
}

func (c Chunk) String() string {
	return fmt.Sprintf("id %s\nsize %d\noffset %d\n\n", c.ID, c.Size, c.Offset)
}
