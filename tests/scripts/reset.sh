#!/usr/bin/env bash

pkill dipd
pkill dipcli

sleep 1

dipd unsafe-reset-all

ssh n2 bash reset.sh
ssh n3 bash reset.sh
ssh n4 bash reset.sh
