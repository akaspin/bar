package manifest
import (
	"encoding/hex"
	"encoding/json"
)

type ID []byte

func (i ID) String() string  {
	return hex.EncodeToString(i)
}

func (i ID) MarshalJSON() (res []byte, err error) {
	return json.Marshal(hex.EncodeToString(i))
}

func (i *ID) UnmarshalJSON(b []byte) (err error) {
	var j string
	if err = json.Unmarshal(b, &j); err != nil {
		return
	}
	*i, err = hex.DecodeString(j)
	return
}