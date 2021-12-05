#!/bin/bash

# Builds the wasm module, the frontend (/webapp) and update github pages (/docs)

echo "Building wasm module, webapp interface and updating github pages (/docs)"

GOOS=js GOARCH=wasm go build -o main.wasm wasm.go && mv main.wasm webapp/public/wasm
npm run --prefix webapp build
rm -rf docs
mv webapp/build ./docs
