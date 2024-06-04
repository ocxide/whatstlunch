#!/bin/bash

rm -rf ./dist/
mkdir -p ./dist

cd ./whatstlunch-server
go build -o ../dist/whatstlunch ./cmd/main.go

cd ../whatstlunch-front
bun install
bunx astro build

mv ./dist ../dist/public/
