package akashi_test

import (
	"fmt"
	"time"

	"github.com/seiyab/akashi"
)

type User struct {
	ID           string
	Name         string
	RegisteredAt time.Time
}

func Example() {
	s := akashi.DiffString(
		User{
			ID:           "123",
			Name:         "Alice",
			RegisteredAt: time.Date(2025, 2, 3, 23, 3, 15, 0, time.UTC),
		},
		User{
			ID:           "123",
			Name:         "Alice",
			RegisteredAt: time.Date(2024, 12, 19, 5, 45, 50, 0, jst),
		},
	)

	fmt.Println(s)

	// Output:
	//   akashi_test.User{
	//     ID: "123",
	//     Name: "Alice",
	// -   RegisteredAt: time.Time("2025-02-03 23:03:15 +0000 UTC"),
	// +   RegisteredAt: time.Time("2024-12-19 05:45:50 +0900 JST"),
	//   }
}
