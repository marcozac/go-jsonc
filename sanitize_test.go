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

//go:build !uncommented_test
// +build !uncommented_test

package jsonc

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/marcozac/go-jsonc/internal/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitize(t *testing.T) {
	t.Parallel()
	for _, tt := range sanitizeTestTargets {
		tt := tt
		name := runtime.FuncForPC(reflect.ValueOf(tt).Pointer()).Name()
		t.Run(name[strings.LastIndex(name, ".")+9:], func(t *testing.T) {
			t.Parallel()
			tt(t)
		})
	}
}

var sanitizeTestTargets = [...]func(t *testing.T){
	SanitizeSmall,
	SanitizeMedium,
}

func SanitizeSmall(t *testing.T) {
	sanitizeTest(t, _small, Small{})
}

func SanitizeMedium(t *testing.T) {
	sanitizeTest(t, _medium, Medium{})
}

func sanitizeTest[T DataType](t *testing.T, data []byte, dt T) {
	t.Helper()
	for _, tt := range []struct {
		Name string
		Func func(t require.TestingT, data []byte, dt T)
	}{
		{"OK", SanitizeOK[T]},
		{"Error", SanitizeError[T]},
	} {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			tt.Func(t, data, dt)
		})
	}
}

func SanitizeOK[T DataType](t require.TestingT, data []byte, dt T) {
	s, err := Sanitize(data)
	require.NoError(t, err, "sanitize failed")
	j := dt
	assert.NoError(t, json.Unmarshal(s, &j), "sanitized JSON is invalid")
	FieldsValue(t, j)
}

func SanitizeError[T DataType](t require.TestingT, data []byte, dt T) {
	_, err := Sanitize(append(data, _invalidChar...))
	assert.ErrorIs(t, err, ErrInvalidUTF8, "invalid UTF-8 was not detected")
}

func BenchmarkSanitize(b *testing.B) {
	b.Run("Small", func(b *testing.B) {
		benchmarkSanitize(b, _small, Small{})
	})
	b.Run("SmallUncommented", func(b *testing.B) {
		benchmarkSanitize(b, _smallUncommented, Small{})
	})
	b.Run("Medium", func(b *testing.B) {
		benchmarkSanitize(b, _medium, Medium{})
	})
	b.Run("MediumUncommented", func(b *testing.B) {
		benchmarkSanitize(b, _mediumUncommented, Medium{})
	})
}

func benchmarkSanitize[T DataType](b *testing.B, data []byte, dt T) {
	b.Helper()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			SanitizeOK(b, data, dt)
		}
	})
}
