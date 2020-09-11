package keys

import (
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

var cdc *codec.Codec

func init() {
	cdc = codec.New()
	codec.RegisterCrypto(cdc)
	cdc.Seal()
}

// MarshalJSON  - marshal keys
func MarshalJSON(o interface{}) ([]byte, error) {
	return cdc.MarshalJSON(o)
}

// UnmarshalJSON - unmarshal json
func UnmarshalJSON(bz []byte, ptr interface{}) error {
	return cdc.UnmarshalJSON(bz, ptr)
}
