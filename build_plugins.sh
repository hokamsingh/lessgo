#!/bin/bash

SRC_DIR="app/src"
OUT_DIR="app/plugins"

mkdir -p "$OUT_DIR"

for file in "$SRC_DIR"/*_module.go; do
  if [ -f "$file" ]; then
    moduleName=$(basename "$file" .go)
    go build -buildmode=plugin -o "$OUT_DIR/${moduleName}_module.so" "$file"
  fi
done
