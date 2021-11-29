#!/bin/bash

# Pass the -p flag to build for production.
if [ $# -ge 1 ] && [ "$1" == "-p" ]; then
    mkdir -p ./dist/cmd
    (cd ./cmd && env GOOS=linux go build -o ../dist/cmd/storybot)
    cp -R ./db ./dist/
    cp ./prod_start.sh ./dist/start.sh
    cp ./end.sh ./dist/
    cp ./env.toml ./dist/
    cp ./command.toml ./dist/
    cp ./command_data.toml ./dist/
    echo "Production build complete."
    exit 0
fi

# Local development build.
(cd ./cmd && go build -o ./storybot.exe)
echo "Development build complete."