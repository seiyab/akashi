package akashi_test

import (
	"fmt"

	"github.com/seiyab/akashi"
)

func Example_slice() {
	s := akashi.DiffString(
		[]string{"a", "b", "c", "d"},
		[]string{"a", "c", "d", "e"},
	)

	fmt.Println(s)

	// Output:
	//   []string{
	//     "a",
	// -   "b",
	//     "c",
	//     "d",
	// +   "e",
	//   }
}
