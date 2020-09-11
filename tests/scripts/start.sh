#!/usr/bin/env bash

export PATH=$PATH:/root/go/bin
pkill dipd
pkill dipcli

sleep 1
nohup dipd start &
sleep 0.1
nohup dipcli rest-server > dipcli.out 2>&1 &
