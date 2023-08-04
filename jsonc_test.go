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
	"encoding/json"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "embed"
)

//go:embed testdata/test.json
var _jsonc []byte

var _invalidJsonc = []byte(`{
	// comment
	"foo": "` + "\xa5" + `"
}`)

func TestSanitize(t *testing.T) {
	for _, tt := range sanitizeTests {
		tt := tt
		name := runtime.FuncForPC(reflect.ValueOf(tt).Pointer()).Name()
		t.Run(name[strings.LastIndex(name, ".")+9:], func(t *testing.T) {
			t.Parallel()
			tt(t)
		})
	}
}

var sanitizeTests = [...]func(t *testing.T){
	SanitizeOK,
	SanitizeError,
}

type T struct {
	Foo   string `json:"foo"`
	Baz   string `json:"baz"`
	Hello string `json:"hello"`
	X     string `json:"x,omitempty"`
}

func SanitizeOK(t *testing.T) {
	s, err := Sanitize(_jsonc)
	require.NoError(t, err, "sanitize failed")
	var j T
	assert.NoError(t, json.Unmarshal(s, &j), "sanitized JSON is invalid")
	fieldsValue(t, j)
	if t.Failed() {
		t.Logf("sanitized JSON: \n%s", s)
	}
}

func SanitizeError(t *testing.T) {
	s, err := Sanitize(_invalidJsonc)
	assert.ErrorIs(t, err, ErrInvalidUTF8, "invalid UTF-8 was not detected")
	if t.Failed() {
		t.Logf("sanitized JSON: \n%s", s)
	}
}

func TestUnmarshal(t *testing.T) {
	for _, tt := range unmarshalTests {
		tt := tt
		name := runtime.FuncForPC(reflect.ValueOf(tt).Pointer()).Name()
		t.Run(name[strings.LastIndex(name, ".")+8:], func(t *testing.T) {
			t.Parallel()
			tt(t)
		})
	}
}

var unmarshalTests = [...]func(t *testing.T){
	UnmarshalOK,
	UnmarshalError,
}

func UnmarshalOK(t *testing.T) {
	var j T
	assert.NoError(t, Unmarshal(_jsonc, &j), "unmarshal failed")
	fieldsValue(t, j)
}

func UnmarshalError(t *testing.T) {
	var j T
	assert.ErrorIs(t, Unmarshal(_invalidJsonc, &j), ErrInvalidUTF8, "invalid UTF-8 was not detected")
}

func fieldsValue(t *testing.T, j T) {
	t.Helper()
	assert.Equal(t, "bar // comment inside a string", j.Foo, "foo field is invalid")
	assert.Equal(t, "qux /* comment block inside a string */", j.Baz, "baz field is invalid")
	assert.Equal(t, "world *//* /* // x */", j.Hello, "hello field is invalid")
}

func BenchmarkSanitize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = Sanitize(_jsonc)
	}
}

//go:embed testdata/no_slash.json
var _jsonNoSlash []byte

func BenchmarkUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var j T
		_ = Unmarshal(_jsonc, &j)
	}
}

func BenchmarkUnmarshalNoSlash(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var j T
		_ = Unmarshal(_jsonNoSlash, &j)
	}
}

func BenchmarkStdUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var j T
		_ = json.Unmarshal(_jsonNoSlash, &j)
	}
}
