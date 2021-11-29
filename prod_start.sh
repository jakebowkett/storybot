#!/bin/bash
sudo setcap CAP_NET_BIND_SERVICE=+eip ./cmd/storybot
nohup ./cmd/storybot -e ./env.toml -c true > out.log 2>&1 &
echo $! > pid.txt