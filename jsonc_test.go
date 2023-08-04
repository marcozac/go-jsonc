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

func SanitizeOK(t *testing.T) {
	s, err := Sanitize(_jsonc)
	require.NoError(t, err, "sanitize failed")
	type J struct {
		Foo   string `json:"foo"`
		Baz   string `json:"baz"`
		Hello string `json:"hello"`
		X     string `json:"x,omitempty"`
	}
	var j J
	assert.NoError(t, json.Unmarshal(s, &j), "sanitized JSON is invalid")
	assert.Equal(t, "bar // comment inside a string", j.Foo, "foo field is invalid")
	assert.Equal(t, "qux /* comment block inside a string */", j.Baz, "baz field is invalid")
	assert.Equal(t, "world *//* /* // x */", j.Hello, "hello field is invalid")
	if t.Failed() {
		t.Logf("sanitized JSON: \n%s", s)
	}
}

func SanitizeError(t *testing.T) {
	s, err := Sanitize([]byte(`{"foo": "` + "\xa5" + `"}`))
	assert.ErrorIs(t, err, ErrInvalidUTF8, "invalid UTF-8 was not detected")
	if t.Failed() {
		t.Logf("sanitized JSON: \n%s", s)
	}
}

func BenchmarkSanitize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = Sanitize(_jsonc)
	}
}
