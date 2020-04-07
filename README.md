# Dipper-Protocol
this is basic finance protocol

# Quick start
## 1.0 install
```
statik -src=client/lcd/swagger-ui -dest=client/lcd -f -m
cd Dipper-Protocol
make install
```
## 1.1 Initialize configuration files and genesis file
```
dipd init dipd-official --chain-id dipd-chain
NOTE: If you have run the tutorial before, you can start from scratch with 
dipd unsafe-reset-all
or by deleting both of the home folders
rm -rf ~/.dip*
```
## 1.2 Add both accounts, with coins to the genesis file
```
dipcli keys add alice
dipcli keys add bob
dipd add-genesis-account $(dipcli keys show alice -a) 10000000000000000stake,10000000000000000pdip
dipd add-genesis-account $(dipcli keys show bob -a) 10000000000000000stake,10000000000000000pdip
```
## 1.3 create validator
```
dipd gentx --amount 1000000stake --commission-rate "0.10" --commission-max-rate "0.20" --commission-max-change-rate "0.10" --pubkey $(dipd tendermint show-validator) --name alice
```
## 1.4 collect gentx
```
dipd collect-gentxs
```
## 1.5 Configure your CLI to eliminate need for chain-id flag
```
dipcli config chain-id dipd-chain
dipcli config output json
dipcli config indent true
dipcli config trust-node true
dipd start --log_level "*:debug" --trace
curl http://127.0.0.1:26657/status
```
# Smart contract property 
## 2.1 deploy contract
```
dipcli vm create --code_file=./contract/demo/demo.bc --from $(dipcli keys show -a alice) --amount=0pdip --gas=1000000
```
## 2.2 query txhash
```
dipcli query tx <txhash>
```
## 2.3 query contract code
```
dipcli query vm code <contract address>
```
## 2.4 call contract method <transfer>
```
dipcli vm call --from $(dipcli keys show -a alice) --contract_addr dip1gtp5xtfnuqpw3dgaxqdk3n8m6d9t4uvwwqt6ms --method transfer --args "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002" --amount 1000000pdip --abi_file ./contract/demo/demo.abi
```
## 2.5 call contract method, such as balanceOf
```
dipcli query account dip1gtp5xtfnuqpw3dgaxqdk3n8m6d9t4uvwwqt6ms
```
## 2.6 call contract method, such as query alice money
```
dipcli keys parse $(dipcli keys show -a alice)
dipcli query vm call $(dipcli keys show -a alice) dip1gtp5xtfnuqpw3dgaxqdk3n8m6d9t4uvwwqt6ms balanceOf "000000000000000000000000DB8822D044FE1C13AA32AF72F27A113E849FC27E" 0pdip ./contract/demo/demo.abi
```
