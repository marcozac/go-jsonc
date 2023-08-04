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
)

// ErrInvalidUTF8 is returned by Sanitize if the data is not valid UTF-8.
var ErrInvalidUTF8 = errors.New("jsonc: invalid UTF-8")

const (
	_ byte = 1 << iota
	_isString
	_isCommentLine
	_isCommentBlock
	_checkNext
)

// Sanitize removes all comments from JSONC data.
// It returns an error if the data is not valid UTF-8.
func Sanitize(data []byte) ([]byte, error) {
	if !utf8.Valid(data) {
		return nil, ErrInvalidUTF8
	}
	return sanitize(data), nil
}

func sanitize(data []byte) []byte {
	var state byte
	return bytes.Map(func(r rune) rune {
		checkNext := state&_checkNext != 0
		state &^= _checkNext
		switch r {
		case '\n':
			state &^= _isCommentLine
		case '"':
			if state&_isString != 0 {
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
