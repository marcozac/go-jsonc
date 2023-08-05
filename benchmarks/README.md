# Benchmark results

The tables below show the performance of [`Unmarshal`](#unmarshal---parse-json-with-comments-into-a-go-value) compared to the standard library's `encoding/json` and other alternative libraries on small and medium data sets.

They are formatted as follows:

| Data set      | s/op                                        | B/op | allocs/op |
| ------------- | ------------------------------------------- | ---- | --------- |
| Set reference | result (Δ% on reference / reference result) | same | same      |

See the files in this directory for the full report.

### Standard library

The tables below show the performance of [`Unmarshal`](#unmarshal---parse-json-with-comments-into-a-go-value) compared to the standard library's `encoding/json` on small and medium data sets.

| **Small data set**                                                                     | s/op                      | B/op                        | allocs/op              |
| -------------------------------------------------------------------------------------- | ------------------------- | --------------------------- | ---------------------- |
| [With comments](../testdata/small.json)                                                | 2.536µ                    | 1.344Ki                     | 22.00                  |
| [Without comments](../testdata/small_uncommented.json) (comment characters in strings) | 2.425µ (+27.17% / 1.907µ) | 1.219Ki (+14.71% / 1.062Ki) | 22.00 (+4.76% / 21.00) |
| [Without comment characters](../testdata/small_no_comment_runes.json)                  | 2.306µ (+16.11% / 1.986µ) | 1.062Ki (~% / 1.062Ki)      | 21.00 (~% / 21.00)     |

| **Medium data set**                                                                    | s/op                      | B/op                        | allocs/op                |
| -------------------------------------------------------------------------------------- | ------------------------- | --------------------------- | ------------------------ |
| [With comments](../testdata/small.json)                                                | 301.2µ                    | 324.7Ki                     | 1.067k                   |
| [Without comments](../testdata/small_uncommented.json) (comment characters in strings) | 202.3µ (+30.86% / 154.6µ) | 148.7Ki (+60.41% / 92.70Ki) | 1.067k (+0.09% / 1.066k) |
| [Without comment characters](../testdata/small_no_comment_runes.json)                  | 170.6µ (+11.63% / 152.8µ) | 92.70Ki (~% / 92.70Ki)      | 1.066k (~% / 1.066k)     |

### With [`github.com/json-iterator/go`](https://github.com/json-iterator/go)

| **Small data set**                                                                     | s/op                      | B/op                    | allocs/op              |
| -------------------------------------------------------------------------------------- | ------------------------- | ----------------------- | ---------------------- |
| [With comments](../testdata/small.json)                                                | 1.632µ                    | 944.0                   | 14.00                  |
| [Without comments](../testdata/small_uncommented.json) (comment characters in strings) | 1.702µ (+11.94% / 1.521µ) | 816.0 (+24.39% / 656.0) | 14.00 (+7.69% / 13.00) |
| [Without comment characters](../testdata/small_no_comment_runes.json)                  | 1.603µ (~% / 1.598µ)      | 656.0 (~% / 656.0)      | 12.00 (~% / 13.00)     |

| **Medium data set**                                                                    | s/op                      | B/op                        | allocs/op                |
| -------------------------------------------------------------------------------------- | ------------------------- | --------------------------- | ------------------------ |
| [With comments](../testdata/small.json)                                                | 245.0µ                    | 407.8Ki                     | 3.484k                   |
| [Without comments](../testdata/small_uncommented.json) (comment characters in strings) | 142.4µ (+42.25% / 100.1µ) | 231.8Ki (+31.90% / 175.7Ki) | 3.484k (+0.06% / 3.482k) |
| [Without comment characters](../testdata/small_no_comment_runes.json)                  | 113.1µ (+17.45% / 96.32µ) | 175.7Ki (+0.01% / 175.7Ki)  | 3.482k (~% / 3.482k)     |

### [`github.com/goccy/go-json`](https://github.com/goccy/go-json)

| **Small data set**                                                                     | s/op                      | B/op                    | allocs/op               |
| -------------------------------------------------------------------------------------- | ------------------------- | ----------------------- | ----------------------- |
| [With comments](../testdata/small.json)                                                | 1.794µ                    | 1.047Ki                 | 10.00                   |
| [Without comments](../testdata/small_uncommented.json) (comment characters in strings) | 1.797µ (+15.38% / 1.557µ) | 928.0 (+20.83% / 768.0) | 10.00 (+11.11% / 9.000) |
| [Without comment characters](../testdata/small_no_comment_runes.json)                  | 1.705µ (+3.30% / 1.651µ)  | 768.0 (~% / 768.0)      | 9.00 (~% / 9.000)       |

| **Medium data set**                                                                    | s/op                      | B/op                        | allocs/op              |
| -------------------------------------------------------------------------------------- | ------------------------- | --------------------------- | ---------------------- |
| [With comments](../testdata/small.json)                                                | 213.1µ                    | 434.9Ki                     | 77.00                  |
| [Without comments](../testdata/small_uncommented.json) (comment characters in strings) | 101.4µ (+83.61% / 55.24µ) | 250.4Ki (+28.94% / 194.2Ki) | 73.00 (+2.82% / 71.00) |
| [Without comment characters](../testdata/small_no_comment_runes.json)                  | 72.60µ (+37.97% / 52.62µ) | 194.2Ki (+0.02% / 194.1Ki)  | 71.00 (~% / 71.00)     |
