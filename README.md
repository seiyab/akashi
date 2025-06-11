# 灯 - Akashi

[![Go Reference](https://pkg.go.dev/badge/github.com/seiyab/akashi.svg)](https://pkg.go.dev/github.com/seiyab/akashi)

`akashi` is a Go library that provides a utility to compute and display the differences between two Go values. The core functionality is exposed through the `DiffString` function.

## Features

-   Compares two Go values of any type.
-   Returns a human-readable string representing the differences.
-   Supports the `fmt.Stringer` interface for readable output. For example, `time.Time` values are displayed as `time.Time("2025-02-03 23:03:15 +0000 UTC")`.

## Installation

To install `akashi`, use `go get`:

```bash
go get github.com/seiyab/akashi
```

## Usage

The primary function in this library is `DiffString`. It takes two `interface{}` values and an optional list of `Option`s.

```go
package main

import (
	"fmt"
	"github.com/seiyab/akashi"
)

func main() {
	type User struct {
		ID           string
		Name         string
		RegisteredAt time.Time
	}

	obj1 := User{
		ID:           "123",
		Name:         "Alice",
		RegisteredAt: time.Date(2025, 2, 3, 23, 3, 15, 0, time.UTC),
	}
	obj2 := User{
		ID:           "123",
		Name:         "Alice",
		RegisteredAt: time.Date(2024, 12, 19, 5, 45, 50, 0, jst),
	}

	diff = akashi.DiffString(obj1, obj2)

	fmt.Println(diff)
	// Output:
	//   akashi_test.User{
	//     ID: "123",
	//     Name: "Alice",
	// -   RegisteredAt: time.Time("2025-02-03 23:03:15 +0000 UTC"),
	// +   RegisteredAt: time.Time("2024-12-19 05:45:50 +0900 JST"),
	//   }
}
```

### `DiffString`

```go
func DiffString(x, y interface{}, options ...Option) string
```

`DiffString` returns a string that represents the difference between `x` and `y`.
You can pass `Option`s to customize the diffing process.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [LICENSE](./LICENSE) file.
