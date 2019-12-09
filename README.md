# Dipper-Protocol
this is basic finance protocol


# Initialize configuration files and genesis file
dpd init dpd-official --chain-id dpd-chain

# Add both accounts, with coins to the genesis file
dpd add-genesis-account $(dpcli keys show alice -a) 10000000000000000stake,10000000000000000dpc\
dpd add-genesis-account $(dpcli keys show bob -a) 1000000000000stake,10000000000000000dpc

# create validator
dpd gentx \
  --amount 1000000stake \
  --commission-rate "0.10" \
  --commission-max-rate "0.20" \
  --commission-max-change-rate "0.10" \
  --pubkey $(dpd tendermint show-validator) \
  --name alice

# collect gentx
dpd collect-gentxs


# Configure your CLI to eliminate need for chain-id flag
dpcli config chain-id dpd-chain\
dpcli config output json\
dpcli config indent true\
dpcli config trust-node true\
dpd start --log_level "*:debug" --trace\
curl http://127.0.0.1:26657/status
