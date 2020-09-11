package sim

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/bank"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/crisis"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/gov"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/params"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/slashing"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	"github.com/Dipper-Labs/Dipper-Protocol/server"
	"github.com/Dipper-Labs/Dipper-Protocol/tests"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	DefaultKeyPass                = "12345678"
	DefaultGenAccountAmount int64 = 100000000000000000
)

type KeyOutput struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
	Seed    string `json:"seed,omitempty"`
}

type GenesisFileAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`
}

func getTestingHomeDirs(name string) (string, string) {
	tmpDir := os.TempDir()
	dipdHome := fmt.Sprintf("%s%s%s.test_dipd", tmpDir, name, string(os.PathSeparator))
	dipcliHome := fmt.Sprintf("%s%s%s.test_dipcli", tmpDir, name, string(os.PathSeparator))
	return dipdHome, dipcliHome
}

func initFixtures(t *testing.T) (chainID, servAddr, port, dipdHome, dipcliHome, p2p2Addr string) {
	dipdHome, dipcliHome = getTestingHomeDirs(t.Name())
	tests.ExecuteT(t, fmt.Sprintf("rm -rf %s ", dipdHome), "")
	tests.ExecuteT(t, fmt.Sprintf("rm -rf %s ", dipcliHome), "")

	executeWriteCheckErr(t, fmt.Sprintf("dipcli keys add --home=%s foo", dipcliHome), DefaultKeyPass)
	executeWriteCheckErr(t, fmt.Sprintf("dipcli keys add --home=%s bar", dipcliHome), DefaultKeyPass)

	chainID = executeInit(t, fmt.Sprintf("dipd init dip-foo -o --home=%s", dipdHome))
	tests.ExecuteT(t, fmt.Sprintf("dipcli config chain-id %s --home=%s", chainID, dipcliHome), "")
	tests.ExecuteT(t, fmt.Sprintf("dipcli config trust-node true --home=%s", dipcliHome), "")

	fooAccAddress := executeGetAccAddress(t, fmt.Sprintf("dipcli keys show foo -a --home=%s", dipcliHome))
	executeWrite(t, fmt.Sprintf("dipd add-genesis-account %s %d%s --home=%s", fooAccAddress, DefaultGenAccountAmount, sdk.NativeTokenName, dipdHome), DefaultKeyPass)

	fooPubkey := executeGetAccAddress(t, fmt.Sprintf("dipd tendermint show-validator --home=%s", dipdHome)) //TODO refact executeGetAccAddress
	executeWrite(t, fmt.Sprintf("dipd gentx --amount 1000000000000pdip --commission-rate 0.10 --commission-max-rate 0.20 --commission-max-change-rate 0.10 --pubkey %s --name foo --home=%s --home-client=%s", fooPubkey, dipdHome, dipcliHome), DefaultKeyPass)
	tests.ExecuteT(t, fmt.Sprintf("dipd collect-gentxs --home=%s", dipdHome), "")

	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)

	p2p2Addr, _, err = server.FreeTCPAddr()
	require.NoError(t, err)

	return
}

func executeWrite(t *testing.T, cmdStr string, writes ...string) (exitSuccess bool) {
	if strings.Contains(cmdStr, "--from") && strings.Contains(cmdStr, "--fee") {
		cmdStr += " --commit"
	}

	exitSuccess, _, _ = executeWriteRetStreams(t, cmdStr, writes...)
	return
}

func executeWriteRetStreams(t *testing.T, cmdStr string, writes ...string) (bool, string, string) {
	proc := tests.GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := proc.StdinPipe.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}

	stdout, stderr, err := proc.ReadAll()
	if err != nil {
		fmt.Println("Err on proc.ReadAll()", err, cmdStr)
	}

	if len(stdout) > 0 {
		t.Log("Stdout:", string(stdout))
	}

	if len(stderr) > 0 {
		t.Log("Stderr:", string(stderr))
	}

	proc.Wait()
	return proc.ExitState.Success(), string(stdout), string(stderr)
}

func executeWriteCheckErr(t *testing.T, cmdStr string, writes ...string) {
	require.True(t, executeWrite(t, cmdStr, writes...))
}

func executeInit(t *testing.T, cmdStr string) (chainID string) {
	_, stderr := tests.ExecuteT(t, cmdStr, DefaultKeyPass)

	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(stderr), &initRes)
	require.NoError(t, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)

	return
}

func executeGetAccAddress(t *testing.T, cmdStr string) (accAddress string) {
	stdout, _ := tests.ExecuteT(t, cmdStr, "")

	accAddress = string([]byte(stdout))
	return
}

func executeGetAccount(t *testing.T, cmdStr string) (acc auth.BaseAccount) {
	out, _ := tests.ExecuteT(t, cmdStr, "")

	var res map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &res)
	require.NoError(t, err, "out %v, err %v", out, err)

	cdc := MakeCodec()

	err = cdc.UnmarshalJSON([]byte(out), &acc)
	require.NoError(t, err, "acc %v, err %v", out, err)

	return
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	params.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	crisis.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		"tendermint/PubKeySecp256k1", nil)
	return cdc
}
