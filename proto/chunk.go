package proto
import (
	"fmt"
	"github.com/akaspin/bar/proto/wire"
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

func (c Chunk) MarshalThrift() (data wire.Chunk, err error) {
	dataInfo, err := c.Data.MarshalThrift()
	if err != nil {
		return
	}
	data = wire.Chunk{&dataInfo, c.Offset}
	return
}

func (c *Chunk) UnmarshalThrift(data wire.Chunk) (err error) {
	if err = c.Data.UnmarshalThrift(*data.Info); err != nil {
		return
	}
	c.Offset = data.Offset
	return
}

////

// slice of unique chunks
type ChunkSlice []Chunk

func (s ChunkSlice) MarshalThrift() (data []*wire.Chunk, err error) {
	for _, chunk := range s {
		var c wire.Chunk
		if c, err = chunk.MarshalThrift(); err != nil {
			return
		}
		data = append(data, &c)
	}
	return
}

func (s *ChunkSlice) UnmarshalThrift(data []*wire.Chunk) (err error) {
	for _, tChunk := range data {
		var c Chunk
		if err = (&c).UnmarshalThrift(*tChunk); err != nil {
			return
		}
	}
	return
}


