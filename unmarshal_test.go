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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	for _, tt := range unmarshalTestTargets {
		tt := tt
		name := runtime.FuncForPC(reflect.ValueOf(tt).Pointer()).Name()
		t.Run(name[strings.LastIndex(name, ".")+10:], func(t *testing.T) {
			t.Parallel()
			tt(t)
		})
	}
}

var unmarshalTestTargets = [...]func(t *testing.T){
	UnmarshalSmall,
	UnmarshalMedium,
}

func UnmarshalSmall(t *testing.T) {
	unmarshalTest(t, _small, Small{})
}

func UnmarshalMedium(t *testing.T) {
	unmarshalTest(t, _medium, Medium{})
}

func unmarshalTest[T DataType](t *testing.T, data []byte, dt T) {
	t.Helper()
	for _, tt := range []struct {
		Name string
		Func func(t require.TestingT, data []byte, dt T)
	}{
		{"OK", UnmarshalOK[T]},
		{"Error", UnmarshalError[T]},
	} {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			tt.Func(t, data, dt)
		})
	}
}

func UnmarshalOK[T DataType](t require.TestingT, data []byte, dt T) {
	j := dt
	assert.NoError(t, Unmarshal(data, &j), "unmarshal failed")
	FieldsValue(t, j)
}

func UnmarshalError[T DataType](t require.TestingT, data []byte, dt T) {
	j := dt
	assert.ErrorIs(t, Unmarshal(append(data, _invalidChar...), &j), ErrInvalidUTF8, "invalid UTF-8 was not detected")
}

func BenchmarkUnmarshal(b *testing.B) {
	b.Run("Small", func(b *testing.B) {
		benchmarkUnmarshal(b, _small, Small{})
	})
	b.Run("SmallUncommented", func(b *testing.B) { // Check skip sanitization
		benchmarkUnmarshal(b, _smallUncommented, Small{})
	})
	b.Run("Medium", func(b *testing.B) {
		benchmarkUnmarshal(b, _medium, Medium{})
	})
	b.Run("MediumUncommented", func(b *testing.B) { // Check skip sanitization
		benchmarkUnmarshal(b, _mediumUncommented, Medium{})
	})
}

func benchmarkUnmarshal[T DataType](b *testing.B, data []byte, dt T) {
	b.Helper()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			UnmarshalOK(b, data, dt)
		}
	})
}
