package manifest
import (
	"fmt"
	"github.com/akaspin/bar/proto/bar"
)

const CHUNK_SIZE = 1024 * 1024 * 2

// Manifest chunk
type Chunk struct  {

	Data

	// Offset
	Offset int64
}

func (c Chunk) String() string {
	return fmt.Sprintf("id %s\nsize %d\noffset %d\n\n", c.ID, c.Size, c.Offset)
}

func (c Chunk) MarshalThrift() (res bar.Chunk, err error) {
	data, err := c.Data.MarshalThrift()
	if err != nil {
		return
	}
	res = bar.Chunk{&data, c.Offset}
	return
}

func (c *Chunk) UnmarshalThrift(tC bar.Chunk) (err error) {
	if err = c.Data.UnmarshalThrift(*tC.Info); err != nil {
		return
	}
	c.Offset = tC.Offset
	return
}
