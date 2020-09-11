#!/usr/bin/env bash

echo txhash:account_address:total:amount:fee

txhashs_file=$1
for txh in `cat $txhashs_file`
do

gas=`dipcli q tx $txh | grep gas_used | awk -F '"' '{print $4}'`
fee=$((gas*1000))
acc=`dipcli q tx $txh | grep from | awk -F '"' '{print $4}'`
amount=`dipcli q account $acc | grep amount | awk -F '"' '{print $4}'`
total=$((amount+fee))
echo $txh:$acc:$total:$amount:$fee

done
