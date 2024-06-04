#!/bin/bash

rm ./dist/whatstlunch
rm -rf ./dist/public

mkdir -p ./dist

cd ./whatstlunch-server
go build -o ../dist/whatstlunch ./cmd/main.go

cd ../whatstlunch-front
bun install
bunx astro build

mv ./dist ../dist/public/
