#!/usr/bin/env bash

# Run all benchmarks with different build tags.
set -e

# Set the path to the temporary directory for the benchmark results.
# The path is relative to the root of the repository.
benchmarkDir="$(dirname "${BASH_SOURCE[0]}")/../tmp/benchmarks"

# Create the temporary directory for the benchmark results
# if it does not exist and add a .gitignore file to it.
mkdir -p "$benchmarkDir"
echo "*" >"$benchmarkDir/.gitignore"

# Set the number of times to run each benchmark.
# The default value is 10.
# It may be overridden by passing a value as the first argument.
count=${1:-10}

# Run the benchmarks for the standard library.
echo "
Running benchmarks for the standard library"
go test -run='^$' -bench=. -benchmem -count "$count" |
  tee "$benchmarkDir/standard_library.txt"

# Run Unmarshal benchmarks
buildTags=('' 'jsoniter' 'go_json')
for t in "${buildTags[@]}"; do
  f="$t"
  [ -n "$f" ] || f="standard_library"

  [ "$f" != "standard_library" ] && echo "
Running Unmarshal benchmarks with build tag: '$t'" &&
    go test -run='^$' -bench=BenchmarkUnmarshal -benchmem -count "$count" -tags="$t" |
    tee "$benchmarkDir/$f.txt"

  echo "
Running Unmarshal benchmarks for uncommented JSON with build tag: '$t'"
  go test -run='^$' -bench=BenchmarkUnmarshal -benchmem -count "$count" -tags="$t,uncommented_test" |
    tee "$benchmarkDir/$f"_uncommented.txt
done
