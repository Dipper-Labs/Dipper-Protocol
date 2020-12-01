package common

import (
	"bytes"
	"testing"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

func TestCreateAddress2(t *testing.T) {
	type testcase struct {
		origin   string
		salt     string
		code     string
		expected string
	}

	for i, tt := range []testcase{
		{
			origin:   "dip1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq8nxcq2",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x00",
			expected: "dip1l2l35x4c06nzdemz4uy3grc0vk2gzphc5cncut",
		},
		{
			origin:   "dip1l2l35x4c06nzdemz4uy3grc0vk2gzphc5cncut",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x00",
			expected: "dip17tj3hmt5fdae3e7j2p3lw2xnt84aanr7ud36cj",
		},
		{
			origin:   "dip17tj3hmt5fdae3e7j2p3lw2xnt84aanr7ud36cj",
			salt:     "0xfeed000000000000000000000000000000000000",
			code:     "0x00",
			expected: "dip1y3rq2v0yy2ugcpkjyfsxh7p3jpfu08zpw2ulx9",
		},
		{
			origin:   "dip1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq8nxcq2",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0xdeadbeef",
			expected: "dip1lxz5m0z5hyj4hgagthm5n4esghpazz8enefnzh",
		},
		{
			origin:   "dip1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq8nxcq2",
			salt:     "0xcafebabe",
			code:     "0xdeadbeef",
			expected: "dip1jh9ny6pmtm42aa2myl2k0mnq0dmrtag0rdvw4h",
		},
		{
			origin:   "dip1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq8nxcq2",
			salt:     "0xcafebabe",
			code:     "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			expected: "dip1ed399kjapfxqqsea6dkk2uxrzepq8cjwddawyr",
		},
		{
			origin:   "dip1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq8nxcq2",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x",
			expected: "dip1sv0mfghrpsp0usxcpc9eax8u2d9t4mx8ean2ym",
		},
	} {
		origin, _ := sdk.AccAddressFromBech32(tt.origin)
		salt := sdk.BytesToHash(FromHex(tt.salt))
		codeHash := crypto.Sha256(FromHex(tt.code))
		address := CreateAddress2(origin, salt, codeHash)

		expected, err := sdk.AccAddressFromBech32(tt.expected)
		if err != nil {
			t.Log(err)
		}
		if !bytes.Equal(expected.Bytes(), address.Bytes()) {
			t.Errorf("test %d: expected %s, got %s", i, expected.String(), address.String())
		}

	}
}
