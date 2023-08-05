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

package jsonc_test

import (
	"fmt"

	"github.com/marcozac/go-jsonc"
)

func ExampleUnmarshal() {
	var v interface{}

	data := []byte(`{/* comment */"foo": "bar"}`)

	err := jsonc.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)

	// Output:
	// map[foo:bar]
}

func ExampleUnmarshal_sanitizeError() {
	var v interface{}

	invalid := []byte(`{/* comment */"foo": "invalid utf8"}`)
	invalid = append(invalid, []byte("\xa5")...)

	err := jsonc.Unmarshal(invalid, &v)
	fmt.Println(err)

	// Output:
	// jsonc: invalid UTF-8
}
