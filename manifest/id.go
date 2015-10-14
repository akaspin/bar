package manifest
import (
	"github.com/akaspin/bar/proto/bar"
	"encoding/hex"
)

// SHA3-256
type ID string

func (i ID) String() string {
	return string(i)
}

func (i ID) Decode(data []byte) (err error) {
	data, err = hex.DecodeString(i.String())
	return
}

func (i ID) MarshalThrift() (res bar.ID, err error) {
	return hex.DecodeString(i.String())
}

func (i *ID) UnmarshalThrift(data bar.ID) (err error) {
	*i = ID(hex.EncodeToString(data))
	return
}