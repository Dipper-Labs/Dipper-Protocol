# Dipper-Protocol
this is basic finance protocol


# Initialize configuration files and genesis file
dipd init dipd-official --chain-id dipd-chain\
NOTE: If you have run the tutorial before, you can start from scratch with a\
dipd unsafe-reset-all\
or by deleting both of the home folders\
rm -rf ~/.dip*

# Add both accounts, with coins to the genesis file
dipcli keys add alice\
dipcli keys add bob\
dipd add-genesis-account $(dipcli keys show alice -a) 10000000000000000stake,10000000000000000dpc,10000000000000000eth,10000000000000000dai\
dipd add-genesis-account $(dipcli keys show bob -a) 10000000000000000stake,10000000000000000dpc,10000000000000000eth,10000000000000000dai


# create validator
dipd gentx 
  --amount 1000000stake 
  --commission-rate "0.10" 
  --commission-max-rate "0.20" 
  --commission-max-change-rate "0.10" 
  --pubkey $(dipd tendermint show-validator) 
  --name alice

# collect gentx
dipd collect-gentxs


# Configure your CLI to eliminate need for chain-id flag
dipcli config chain-id dipd-chain\
dipcli config output json\
dipcli config indent true\
dipcli config trust-node true\
dipd start --log_level "*:debug" --trace\
curl http://127.0.0.1:26657/status

# You can withdraw/deposit/borrow/repay coin which bank supports.
dipcli tx dipperBank set-oracleprice dipperBank eth 150000000 --from bob\
dipcli tx dipperBank bank-withdraw 12000000eth eth --from bob\
dipcli tx dipperBank bank-withdraw 150000000dai dai --from alice\
dipcli tx dipperBank bank-deposit 1000000eth eth --from bob\
dipcli tx dipperBank bank-borrow 120000000dai dai --from bob\
dipcli tx dipperBank bank-repay 120000000dai dai --from bob

#deploy contract
dipcli vm create --code_file=./demo/demo.bc \
--from $(dipcli keys show -a alice) --amount=0pdip \
--gas=1000000
##query txhash
dipcli query tx <txhash>
##query contract code
dipcli query vm code <contract address>
##call contract method <transfer>
dipcli vm call --from $(dipcli keys show -a alice) \
--contract_addr dip1gtp5xtfnuqpw3dgaxqdk3n8m6d9t4uvwwqt6ms \
--method transfer  \
--args  "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002" \
--amount 1000000pdip \
--abi_file ./demo/demo.abi
##call contract method <blanceof> 
dipcli query account dip1gtp5xtfnuqpw3dgaxqdk3n8m6d9t4uvwwqt6ms
##call contract method <query>
dipcli keys parse $(dipcli keys show -a alice)
dipcli query vm call $(dipcli keys show -a alice) dip1gtp5xtfnuqpw3dgaxqdk3n8m6d9t4uvwwqt6ms balanceOf "000000000000000000000000DB8822D044FE1C13AA32AF72F27A113E849FC27E" 0pdip ./demo/demo.abi
