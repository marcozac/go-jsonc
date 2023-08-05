# jsonc - JSON with comments for Go

## Description

jsonc is a Go package for working with JSONC files.

JSONC is an extension to JSON that allows comments. This package provides functions to:

- Sanitize JSONC data by removing comments
- Unmarshal JSONC into Go values

## Installation

Install the jsonc package:

`go get github.com/marcozac/go-jsonc`

## Usage

### Sanitize JSON with comments

Remove comments making the data compatible to JSON.

```go
package main

import (
    "encoding/json"

    "github.com/marcozac/go-jsonc"
)

func main() {
    invalidData := []byte(`{
        // a comment
        "foo": "bar"
    }`)

    // Remove comments from JSONC
    data, err := jsonc.Sanitize(invalidData)
    if err != nil {
        ...
    }

    var v struct{ Foo string }

    // Unmarshal using any other library
    err = json.Unmarshal(data, &v)
    if err != nil {
        ...
    }
}
```

### Unmarshal JSON with comments data into a Go value

The Unmarshal function takes JSONC data and a pointer to the target Go value to unmarshal into.

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

## License

This project is licensed under the Apache 2.0 license. See LICENSE file for details.
