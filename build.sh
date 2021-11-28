#!/bin/bash

# Pass the -p flag to build for the server.
if [ $# -ge 1 ] && [ "$1" == "-p" ]; then
    mkdir -p ./build/cmd
    (cd ./cmd && env GOOS=linux go build -o ../build/cmd/storybot.exe)
    cp -R ./db ./build/
    cp ./start.sh ./build/
    cp ./env.toml ./build/
    cp ./command*.toml ./build/
else
    (cd ./cmd && go build -o ./storybot.exe)
fi