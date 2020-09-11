package sim

import (
	"fmt"
	"testing"

	"github.com/Dipper-Labs/Dipper-Protocol/tests"
)

func TestMock(t *testing.T) {
	t.Parallel()

	_, servAddr, port, dipdHome, dipcliHome, p2pAddr := initFixtures(t)

	dipdStartCmd := fmt.Sprintf("dipd start --home=%s --rpc.laddr=%v --p2p.laddr=%v", dipdHome, servAddr, p2pAddr)
	proc := tests.GoExecuteTWithStdout(t, dipdStartCmd)
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(1, port)

	fooAddr := executeGetAccAddress(t, fmt.Sprintf("dipcli keys show foo -a --home=%s", dipcliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("dipcli query account %s --home=%s --node tcp://localhost:%s -o json", fooAddr, dipcliHome, port))
	fmt.Println(fooAcc.Coins.String())
}
