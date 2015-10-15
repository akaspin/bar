package proto
import (
	"github.com/akaspin/bar/proto/wire"
)

// Data chunk
type Data struct  {
	ID ID
	Size int64
}

func (d Data) MarshalThrift() (tData wire.DataInfo, err error)  {
	tData.Id = wire.ID(d.ID)
	tData.Size = d.Size
	return
}

func (d *Data) UnmarshalThrift(tData wire.DataInfo) (err error) {
	d.ID = ID(tData.Id)
	d.Size = tData.Size
	return
}