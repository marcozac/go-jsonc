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
// It returns an error if the data is not valid UTF-8.
func Sanitize(data []byte) ([]byte, error) {
	if !utf8.Valid(data) {
		return nil, ErrInvalidUTF8
	}
	return sanitize(data), nil
}

const (
	_ byte = 1 << iota
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
// It uses the standard library by default, but can be configured to use the
// jsoniter or go-json library instead by using build tags.
//
//	| tag           | library   |
//	|---------------|-----------|
//	| go_json	| go-json   |
//	| jsoniter	| jsoniter  |
//	| none or both	| standard  |
//
// If the data contains any '/' character, it is assumed to be not sanitized,
// and [Sanitize] will be called before unmarshaling returning an error if the
// data is not valid UTF-8.
//
// To improve performance, if the data does not contain any '/' character,
// it is assumed to be sanitized, and [Sanitize] will not be called, adding
// just a small overhead to the unmarshaling process if the data does not
// contain comments.
//
// Caveat: if the data contains a string that looks like a comment, for
// example: {"url": "http://example.com"}, Unmarshal calls [Sanitize] anyway,
// even if the data does not contain any comment.
func Unmarshal(data []byte, v any) error {
	if bytes.ContainsRune(data, '/') {
		var err error
		data, err = Sanitize(data)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(data, v)
}
