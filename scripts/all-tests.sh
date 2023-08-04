#!/usr/bin/env bash

# Run all tests with different build tags and the race detector enabled.
set -e

buildTags=('' 'jsoniter' 'go_json' 'jsoniter,go_json')
for t in "${buildTags[@]}"; do
  echo "
Running tests with build tag: $t"
  go test -v -race ./... -tags="$t"
done
