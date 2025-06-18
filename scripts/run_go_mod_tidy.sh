#!/bin/bash

for file in "$@"; do
  dir=$(dirname "$file")
  (cd "$dir" && go mod tidy)
done
