# jsonc - JSON with comments for Go

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/marcozac/go-jsonc)
![License](https://img.shields.io/github/license/marcozac/go-jsonc?color=blue)
[![CI](https://github.com/marcozac/go-jsonc/actions/workflows/ci.yml/badge.svg)](https://github.com/marcozac/go-jsonc/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/marcozac/go-jsonc/branch/main/graph/badge.svg?token=JYj7gCZauN)](https://codecov.io/gh/marcozac/go-jsonc)
[![Go Report Card](https://goreportcard.com/badge/github.com/marcozac/go-jsonc)](https://goreportcard.com/report/github.com/marcozac/go-jsonc)

`jsonc` is a light and dependency-free package for working with JSON with comments data built on top of `encoding/json`.
It allows to remove comments converting to valid JSON-encoded data and to unmarshal JSON with comments into Go values.

The dependencies listed in [go.mod](/go.mod) are only used for testing and benchmarking or to support [alternative libraries](#alternative-libraries).

## Features

- Full support for comment lines and block comments
- Preserve the content of strings that contain comment characters
- Sanitize JSON with comments data by removing comments
- Unmarshal JSON with comments into Go values

## Installation

Install the `jsonc` package:

```bash
go get github.com/marcozac/go-jsonc
```

## Usage

### Sanitize - Remove comments from JSON data

`Sanitize` removes all comments from JSON data, returning valid JSON-encoded byte slice that is compatible with standard library's json.Unmarshal.

It works with comment lines and block comments anywhere in the JSONC data, preserving the content of strings that contain comment characters.

#### Example

```go
package main

import (
    "encoding/json"

    "github.com/marcozac/go-jsonc"
)

func main() {
    invalidData := []byte(`{
        // a comment
        "foo": "bar" /* a comment in a weird place */,
        /*
            a block comment
        */
        "hello": "world" // another comment
    }`)

    // Remove comments from JSONC
    data, err := jsonc.Sanitize(invalidData)
    if err != nil {
        ...
    }

    var v struct{
      Foo   string
      Hello string
    }

    // Unmarshal using any other library
    if err := json.Unmarshal(data, &v); err != nil {
        ...
    }
}
```

### Unmarshal - Parse JSON with comments into a Go value

`Unmarshal` replicates the behavior of the standard library's json.Unmarshal function, with the addition of support for comments.

It is optimized to avoid calling [`Sanitize`](#sanitize---remove-comments-from-json-data) unless it detects comments in the data.
This avoids the overhead of removing comments when they are not present, improving performance on small data sets.

It first checks if the data contains comment characters as `//` or `/*` using [`HasCommentRunes`](https://pkg.go.dev/github.com/marcozac/go-jsonc#HasCommentRunes).
If no comment characters are found, it directly unmarshals the data.

Only if comments are detected it calls [`Sanitize`](#sanitize---remove-comments-from-json-data) before unmarshaling to remove them.
So, `Unmarshal` tries to skip unnecessary work when possible, but currently it is not possible to detect false positives as `//` or `/*` inside strings.

Since the comment detection is based on a simple rune check, it is not recommended to use `Unmarshal` on large data sets unless you are not sure whether they contain comments.
Indeed, `HasCommentRunes` needs to checks every single byte before to return `false` and may drastically slow down the process.

In this case, it is more efficient to call [`Sanitize`](#sanitize---remove-comments-from-json-data) before to unmarshal the data.

#### Example

```go
package main

import "github.com/marcozac/go-jsonc"

func main() {
    invalidData := []byte(`{
        // a comment
        "foo": "bar"
    }`)

    var v struct{ Foo string }

    err := jsonc.Unmarshal(invalidData, &v)
    if err != nil {
    ...
    }
}
```

## Alternative libraries

By default, `jsonc` uses the standard library's `encoding/json` to unmarshal JSON data and has no external dependencies.

It is possible to use build tags to use alternative libraries instead of the standard library's `encoding/json`:

| Tag          | Library                                                              |
| ------------ | -------------------------------------------------------------------- |
| none or both | standard library                                                     |
| jsoniter     | [`github.com/json-iterator/go`](https://github.com/json-iterator/go) |
| go_json      | [`github.com/goccy/go-json`](https://github.com/goccy/go-json)       |

## Benchmarks

This library aims to have performance comparable to the standard library's `encoding/json`.
Unfortunately, comments removal is not free and it is not possible to avoid the overhead of removing comments when they are present.

Currently `jsonc` performs worse than the standard library's `encoding/json` on small data sets about 27% on data with comments in strings and 16% on data without comments.
On medium data sets, the performance gap is increased to about 30% on data with comments in strings and reduced to 12% on data without comments.

However, using one of the [alternative libraries](#alternative-libraries), it is possible to achieve better performance than the standard library's `encoding/json` even considering the overhead of removing comments.

See [benchmarks](/benchmarks) for the full results.

The benchmarks are run on a MacBook Pro (16-inch, 2021), Apple M1 Max, 32 GB RAM.

## Contributing

:heart: Contributions are ~~needed~~ welcome!

Please open an issue or submit a pull request if you would like to contribute.

To submit a pull request:

- Fork this repository
- Create a new branch
- Make changes and commit
- Push to your fork and submit a pull request

## License

This project is licensed under the Apache 2.0 license. See [LICENSE](/LICENSE) file for details.
