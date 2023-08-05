#!/usr/bin/env bash

# Run all benchmarks with different build tags.
set -e

# Set the path to the temporary directory for the benchmark results.
# The path is relative to the root of the repository.
benchmarkDir="$(dirname "${BASH_SOURCE[0]}")/../tmp/benchmark"

# Create the temporary directory for the benchmark results
# if it does not exist and add a .gitignore file to it.
mkdir -p "$benchmarkDir"
echo "*" >"$benchmarkDir/.gitignore"

# Set the number of times to run each benchmark.
# The default value is 10.
# It may be overridden by passing a value as the first argument.
count=${1:-10}

buildTags=('' 'jsoniter' 'go_json')
for t in "${buildTags[@]}"; do
  f="$t"
  [ -n "$f" ] || f="standard_library"

  echo "
Running benchmarks with build tag: '$t'"
  go test -run='^$' -bench=. -benchmem -count "$count" -tags="$t" | tee "$benchmarkDir/$f.txt"

  echo "
Running benchmarks for uncommented JSON with build tag: '$t'"
  go test -run='^$' -bench=. -benchmem -count "$count" -tags="$t,uncommented_test" | tee "$benchmarkDir/$f"_uncommented.txt
done
