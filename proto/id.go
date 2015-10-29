package proto

import (
	"encoding/hex"
	"github.com/akaspin/bar/proto/wire"
)

// SHA3-256
type ID string

func (i ID) String() string {
	return string(i)
}

func (i *ID) UnmarshalBinary(data []byte) (err error) {
	*i = ID(hex.EncodeToString(data))
	return
}

func (i ID) Decode(data []byte) (err error) {
	data, err = hex.DecodeString(i.String())
	return
}

func (i ID) MarshalThrift() (res wire.ID, err error) {
	return hex.DecodeString(i.String())
}

func (i *ID) UnmarshalThrift(data wire.ID) (err error) {
	*i = ID(hex.EncodeToString(data))
	return
}

type IDSlice []ID

// NOTE: strange behaviour of thrift compiller should be []ID.
func (i IDSlice) MarshalThrift() (res [][]byte, err error) {
	for _, id := range i {
		var id1 wire.ID
		if id1, err = id.MarshalThrift(); err != nil {
			return
		}
		res = append(res, []byte(id1))
	}

	return
}

func (i *IDSlice) UnmarshalThrift(data [][]byte) (err error) {
	for _, d := range data {
		var id ID
		if err = (&id).UnmarshalThrift(d); err != nil {
			return
		}
		*i = append(*i, id)
	}
	return
}
