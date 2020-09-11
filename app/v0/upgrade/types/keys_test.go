package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUintToHexString(t *testing.T) {
	require.Equal(t, UintToHexString(15), "000000000000000f")
}
