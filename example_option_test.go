package akashi_test

import (
	"fmt"
	"time"

	"github.com/seiyab/akashi"
)

func ExampleWithFormat() {
	s := akashi.DiffString(
		time.Date(2025, 2, 3, 23, 3, 15, 0, time.UTC),
		time.Date(2024, 12, 19, 5, 45, 50, 0, time.UTC),
		akashi.WithFormat(
			func(t time.Time) string {
				return t.Format(time.RFC3339)
			},
		),
	)

	fmt.Println(s)

	// Output:
	// - time.Time("2025-02-03T23:03:15Z")
	// + time.Time("2024-12-19T05:45:50Z")
}
