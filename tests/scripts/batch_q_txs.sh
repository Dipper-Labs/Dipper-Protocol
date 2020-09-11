#!/usr/bin/env bash

for tx in `cat txs`
do
dipcli q tx $tx
if [ $? -eq 1 ]; then
echo $tx >>badTxs
fi
done
