package manifest
import (
	"github.com/akaspin/bar/proto/bar"
	"encoding/hex"
)

// Data chunk
type Data struct  {
	ID string
	Size int64
}

func (d Data) MarshalThrift() (tData bar.DataInfo, err error)  {
	if tData.Id, err = hex.DecodeString(d.ID); err != nil {
		return
	}
	tData.Size = d.Size
	return
}

func (d *Data) UnmarshalThrift(tData bar.DataInfo) (err error) {
	d.ID = hex.EncodeToString(tData.Id)
	d.Size = tData.Size
	return
}