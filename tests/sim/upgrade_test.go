package sim

import (
	"fmt"
	"testing"

	"github.com/Dipper-Labs/Dipper-Protocol/tests"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/stretchr/testify/require"
)

func accountQueryCmd(t *testing.T, acc, dipcliHome, port string) string {
	fooAddr := executeGetAccAddress(t, fmt.Sprintf("dipcli keys show %s -a --home=%s", acc, dipcliHome))
	return fmt.Sprintf("dipcli query account %s --home=%s --node tcp://localhost:%s -o json", fooAddr, dipcliHome, port)
}

func Test_Upgrade(t *testing.T) {
	t.Parallel()
	_, servAddr, port, dipdHome, dipcliHome, p2pAddr := initFixtures(t)

	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("dipd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", dipdHome, servAddr, p2pAddr))
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(1, port)

	cmdGetAccount := accountQueryCmd(t, "foo", dipcliHome, port)
	fooAccount := executeGetAccount(t, cmdGetAccount)
	//check foo account init coins
	require.Equal(t, fooAccount.Coins.AmountOf(sdk.NativeTokenName), sdk.NewInt(DefaultGenAccountAmount-1000000000000))
}
