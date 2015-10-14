package manifest
import (
	"github.com/akaspin/bar/proto/bar"
)

// Data chunk
type Data struct  {
	ID ID
	Size int64
}

func (d Data) MarshalThrift() (tData bar.DataInfo, err error)  {
	tData.Id = bar.ID(d.ID)
	tData.Size = d.Size
	return
}

func (d *Data) UnmarshalThrift(tData bar.DataInfo) (err error) {
	d.ID = ID(tData.Id)
	d.Size = tData.Size
	return
}