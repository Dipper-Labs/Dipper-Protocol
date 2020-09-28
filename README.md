# Dipper-Protocol
the next generation of basic finance protocol.


# Quick start
To build a private blockchain, if you want to jion our testnet or get more detail, you can click [this link](http://docs.dippernetwork.com "DIP").
## 1.0 install
```
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
dipd add-genesis-account $(dipcli keys show alice -a) 10000000000000000pdip
dipd add-genesis-account $(dipcli keys show bob -a) 10000000000000000pdip
```
## 1.3 create validator
```
dipd gentx --amount 1000000pdip --commission-rate "0.10" --commission-max-rate "0.20" --commission-max-change-rate "0.10" --pubkey $(dipd tendermint show-validator) --name alice
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
## 1.6 transfer asset
```
dipcli send --from $(dipcli keys show bob -a)  --to $(dipcli keys show alice -a) --amount 1000000000000pdip
 ```

## 1.7 query account
```
dipcli query account  $(dipcli keys show jack -a)
dipcli query account  $(dipcli keys show alice -a)
```
# Smart contract property 
## 2.1 deploy contract
```
dipcli vm create --code_file=./contract/demo/demo.bc --abi_file=./contract/demo/demo.abi --from $(dipcli keys show -a alice) --args '' --amount=0pdip --gas=1000000
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
dipcli vm call --from $(dipcli keys show -a alice) --contract_addr=dip1jd8jqhnruunhrxh75da02dm7fr29cdkqtq8wmq --abi_file ./contract/demo/demo.abi --method=transfer --args 'dip1dcu73lw9uqkygpde4z4z22f079skta49vxs2r0 10'  --amount 0pdip --gas 2000000
```
## 2.5 call contract method, such as balanceOf
```
dipcli query account dip1gtp5xtfnuqpw3dgaxqdk3n8m6d9t4uvwwqt6ms
```
## 2.6 call contract method, such as query alice money
```
dipcli query vm call $(dipcli keys show -a alice) dip1gcwk24al08lul80aejyq409mjgtqfu9uhgwtw4 balanceOf ./contract/demo/demo.abi --args "dip16g54d2akrlln48j5p7gcv4nucfzdn2zsxe54j4" --amount 0pdip
```

