#!/bin/bash

set -e

go build -o shout cmd/main.go
mkdir -p ~/.local/bin
ln -sf "$(pwd)/shout" ~/.local/bin/llm
