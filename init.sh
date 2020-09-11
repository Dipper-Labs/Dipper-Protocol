#!/bin/bash
set -x

passwd="11111111"

kill -9 $(pgrep dipd)
kill -9 $(pgrep dipcli)

rm -rf ~/.dipd
#rm -rf ~/.dipcli

# set moniker and chain-id
dipd init mymoniker --chain-id dip-chain

# set up config for CLI
dipcli config chain-id dip-chain
dipcli config output json
dipcli config indent true
dipcli config trust-node true

# add keys
echo -e "${passwd}\n${passwd}\n" | dipcli keys add alice
echo -e "${passwd}\n${passwd}\n" | dipcli keys add bob
echo -e "${passwd}\n${passwd}\n" | dipcli keys add jack

# add genesis account
vesting_start_time=$(date +%s)
if [ "$(uname)" == "Darwin" ]; then
    vesting_end_time=$(date -v+1d +%s)
else
    vesting_end_time=$(date --date="+1 day" +%s)
fi
dipd add-genesis-account $(dipcli keys show alice -a) 30000000000000000000pdip
dipd add-genesis-account $(dipcli keys show bob -a) 40000000000000000000pdip --vesting-amount 5000000000000000000pdip --vesting-start-time ${vesting_start_time} --vesting-end-time ${vesting_end_time}
dipd add-genesis-account $(dipcli keys show jack -a) 30000000000000000000pdip --vesting-amount 5000000000000000000pdip  --vesting-end-time ${vesting_end_time}

echo "${passwd}" | dipd gentx \
  --amount 1000000000000pdip \
  --commission-rate "0.10" \
  --commission-max-rate "0.20" \
  --commission-max-change-rate "0.10" \
  --pubkey $(dipd tendermint show-validator) \
  --name alice

# collect genesis tx
dipd collect-gentxs

# validate genesis file
dipd validate-genesis

# start the node
dipd start --log_level "*:debug" --trace
