// Copyright 2023 Marco Zaccaro. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonc

import (
	"bytes"
	"errors"
	"unicode/utf8"

	"github.com/marcozac/go-jsonc/internal/json"
)

// ErrInvalidUTF8 is returned by Sanitize if the data is not valid UTF-8.
var ErrInvalidUTF8 = errors.New("jsonc: invalid UTF-8")

// Sanitize removes all comments from JSONC data.
// It returns [ErrInvalidUTF8] if the data is not valid UTF-8.
//
// NOTE: it does not checks whether the data is valid JSON or not.
func Sanitize(data []byte) ([]byte, error) {
	if !utf8.Valid(data) {
		return nil, ErrInvalidUTF8
	}
	return sanitize(data), nil
}

const (
	_hasCommentRunes byte = 1 << iota
	_isString
	_isCommentLine
	_isCommentBlock
	_checkNext
)

func sanitize(data []byte) []byte {
	var state byte
	return bytes.Map(func(r rune) rune {
		checkNext := state&_checkNext != 0
		state &^= _checkNext
		switch r {
		case '\n':
			state &^= _isCommentLine
		case '\\':
			if state&_isString != 0 {
				state |= _checkNext
			}
		case '"':
			if state&_isString != 0 {
				if checkNext { // escaped quote
					break // switch => write rune
				}
				state &^= _isString
			} else if state&(_isCommentLine|_isCommentBlock) == 0 {
				state |= _isString
			}
		case '/':
			if state&_isString != 0 {
				break // switch => write rune
			}
			if state&_isCommentBlock != 0 {
				if checkNext {
					state &^= _isCommentBlock
				} else {
					state |= _isCommentLine
				}
			} else {
				if checkNext {
					state |= _isCommentLine
				} else {
					state |= _checkNext
				}
			}
			return -1 // mark rune for skip
		case '*':
			if state&_isString != 0 {
				break // switch => write rune
			}
			if checkNext {
				state |= _isCommentBlock
			} else if state&_isCommentBlock != 0 {
				state |= _checkNext
			}
			return -1 // mark rune for skip
		}
		if state&(_isCommentLine|_isCommentBlock) != 0 {
			return -1 // mark rune for skip
		}
		return r
	}, data)
}

// Unmarshal parses the JSONC-encoded data and stores the result in the value
// pointed by v removing all comments from the data (if any).
//
// It uses [HasCommentRunes] to check whether the data contains any comment.
// Note that this operation is as expensive as the larger the data. On small
// data sets it just adds a small overhead to the unmarshaling process, but
// on large data sets it may have a significant impact on performance. In such
// cases, it may be more efficient to call [Sanitize] and then the standard
// (or any other) library directly.
//
// If the data contains comment runes, it calls [Sanitize] to remove them and
// returns [ErrInvalidUTF8] if the data is not valid UTF-8.
//
// Any error is reported from [json.Unmarshal] as is.
//
// It uses the standard library for unmarshaling by default, but can be
// configured to use the jsoniter or go-json library instead by using build
// tags.
//
//	| tag           | library                       |
//	|---------------|-------------------------------|
//	| none or both	| standard library              |
//	| go_json	| "github.com/goccy/go-json"    |
//	| jsoniter	| "github.com/json-iterator/go" |
//
// Example:
//
//	data := []byte(`{/* comment */"name": "John", "age": 30}`)
//	type T struct {
//		Name string
//		Age  int
//	}
//	var t T
//	err := jsonc.Unmarshal(data, &t)
//	...
func Unmarshal(data []byte, v any) error {
	if HasCommentRunes(data) {
		var err error
		data, err = Sanitize(data)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(data, v)
}

// HasCommentRunes returns true if the data contains any comment rune.
// It checks whether the data contains any '/' character, and if so, it looks
// whether the previous one is a '/' or the next one is a '/' or a '*'.
// If not, it returns false.
//
// Caveat: if the data contains a string that looks like a comment as
// '{"url": "http://example.com"}', HasCommentRunes returns true.
//
// For example, it returns true for the following data:
//
//	{
//		// comment
//		"key": "value"
//	}
//
// or
//
//	{
//		/* comment
//		"key": "value"
//		*/
//		"foo": "bar"
//	}
//
// But also for:
//
//	{ "key": "value // comment" }
func HasCommentRunes(data []byte) bool {
	var state byte
	bytes.IndexFunc(data, func(r rune) bool {
		if state&_checkNext != 0 {
			if r == '/' || r == '*' {
				state |= _hasCommentRunes
				return true
			}
			state &^= _checkNext
		}
		if r == '/' {
			state |= _checkNext
		}
		return false
	})
	return state&_hasCommentRunes != 0
}
