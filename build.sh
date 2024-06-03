#!/bin/bash

cd ./whatstlunch-server
go build -o ../whatstlunch ./cmd/main.go

cd ../whatstlunch-front
bun install
bunx astro build

rm -rf ../public
mv ./dist ../public
