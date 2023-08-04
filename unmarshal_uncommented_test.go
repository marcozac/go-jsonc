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

//go:build uncommented_test
// +build uncommented_test

package jsonc

import (
	"testing"

	"github.com/marcozac/go-jsonc/internal/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file does not contain real benchmarks, but it is used to compare the
// performance of the Unmarshal function with the standard json.Unmarshal
// function (or jsoniter / go-json / ...) and how fast the Unmarshal function
// skips the sanitization of the data.

func BenchmarkUnmarshal(b *testing.B) {
	b.Run("SmallUncommented", func(b *testing.B) {
		benchmarkUnmarshal(b, _smallUncommented, Small{})
	})
	b.Run("MediumUncommented", func(b *testing.B) {
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

func UnmarshalOK[T DataType](t require.TestingT, data []byte, dt T) {
	j := dt
	assert.NoError(t, json.Unmarshal(data, &j), "unmarshal failed")
	FieldsValue(t, j)
}
