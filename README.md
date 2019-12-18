# Dipper-Protocol
this is basic finance protocol


# Initialize configuration files and genesis file
dpd init dpd-official --chain-id dpd-chain\
NOTE: If you have run the tutorial before, you can start from scratch with a\
dpd unsafe-reset-all\
or by deleting both of the home folders\
rm -rf ~/.dp*

# Add both accounts, with coins to the genesis file
dpcli keys add alice\
dpcli keys add bob\
dpd add-genesis-account $(dpcli keys show alice -a) 10000000000000000stake,10000000000000000dpc,10000000000000000eth,10000000000000000dai\
dpd add-genesis-account $(dpcli keys show bob -a) 10000000000000000stake,10000000000000000dpc,10000000000000000eth,10000000000000000dai

# create validator
dpd gentx 
  --amount 1000000stake 
  --commission-rate "0.10" 
  --commission-max-rate "0.20" 
  --commission-max-change-rate "0.10" 
  --pubkey $(dpd tendermint show-validator) 
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

# You can withdraw/deposit/borrow/repay coin which bank supports.
dpcli tx dipperProtocol set-oracleprice dipperBank eth 150000000 --from bob\
dpcli tx dipperProtocol bank-withdraw 12000000eth eth --from bob\
dpcli tx dipperProtocol bank-withdraw 150000000dai dai --from alice\
dpcli tx dipperProtocol bank-deposit 1000000eth eth --from bob\
dpcli tx dipperProtocol bank-borrow 120000000dai dai --from bob\
dpcli tx dipperProtocol bank-repay 120000000dai dai --from bob
