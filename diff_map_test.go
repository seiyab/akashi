package akashi_test

import (
	"strings"
	"testing"
)

func TestDiff_Map(t *testing.T) {
	t.Run("map of primitive", func(t *testing.T) {
		type testCase struct {
			name  string
			left  map[string]int
			right map[string]int
			want  string
		}
		for _, tc := range []testCase{
			{
				name:  "identical",
				left:  map[string]int{"a": 1, "b": 2},
				right: map[string]int{"a": 1, "b": 2},
				want:  "",
			},
			{
				name:  "different values",
				left:  map[string]int{"a": 1, "b": 2},
				right: map[string]int{"a": 1, "b": 3},
				want: strings.Join([]string{
					`  map[string]int{`,
					`    "a": 1,`,
					`-   "b": 2,`,
					`+   "b": 3,`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "different keys",
				left:  map[string]int{"a": 1, "b": 2},
				right: map[string]int{"a": 1, "c": 3},
				want: strings.Join([]string{
					`  map[string]int{`,
					`    "a": 1,`,
					`-   "b": 2,`,
					`+   "c": 3,`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "nil maps",
				left:  nil,
				right: map[string]int{"a": 1},
				want: strings.Join([]string{
					`- map[string]int(nil)`,
					`+ map[string]int{`,
					`+   "a": 1,`,
					`+ }`,
				}, "\n"),
			},
			{
				name:  "nil maps (right)",
				left:  map[string]int{"a": 1},
				right: nil,
				want: strings.Join([]string{
					`- map[string]int{`,
					`-   "a": 1,`,
					`- }`,
					`+ map[string]int(nil)`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})

	t.Run("map with int key", func(t *testing.T) {
		type testCase struct {
			name  string
			left  map[int]string
			right map[int]string
			want  string
		}
		for _, tc := range []testCase{
			{
				name:  "identical",
				left:  map[int]string{1: "a", 2: "b"},
				right: map[int]string{1: "a", 2: "b"},
				want:  "",
			},
			{
				name:  "different values",
				left:  map[int]string{1: "a", 2: "b"},
				right: map[int]string{1: "a", 2: "c"},
				want: strings.Join([]string{
					`  map[int]string{`,
					`    1: "a",`,
					`-   2: "b",`,
					`+   2: "c",`,
					`  }`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})

	t.Run("nested map", func(t *testing.T) {
		type testCase struct {
			name  string
			left  map[string]map[string]int
			right map[string]map[string]int
			want  string
		}
		for _, tc := range []testCase{
			{
				name: "different nested values",
				left: map[string]map[string]int{
					"x": {"a": 1, "b": 2},
					"y": {"c": 3},
				},
				right: map[string]map[string]int{
					"x": {"a": 1, "b": 3},
					"y": {"c": 3},
				},
				want: strings.Join([]string{
					`  map[string]map[string]int{`,
					`    "x": map[string]int{`,
					`      "a": 1,`,
					`-     "b": 2,`,
					`+     "b": 3,`,
					`    },`,
					`    "y": map[string]int{`,
					`:`,
					`  }`,
				}, "\n"),
			},
			{
				name: "nil nested map",
				left: map[string]map[string]int{
					"x": {"a": 1},
					"y": nil,
				},
				right: map[string]map[string]int{
					"x": {"a": 1},
					"y": {"b": 2},
				},
				want: strings.Join([]string{
					`  map[string]map[string]int{`,
					`    "x": map[string]int{`,
					`      "a": 1,`,
					`    },`,
					`-   "y": map[string]int(nil),`,
					`+   "y": map[string]int{`,
					`+     "b": 2,`,
					`+   },`,
					`  }`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})

	t.Run("map with interface key", func(t *testing.T) {
		type testCase struct {
			name  string
			left  map[interface{}]int
			right map[interface{}]int
			want  string
		}
		for _, tc := range []testCase{
			{
				name:  "mixed key types",
				left:  map[interface{}]int{1: 10, "x": 20, true: 30},
				right: map[interface{}]int{1: 10, "x": 25, true: 30},
				want: strings.Join([]string{
					`  map[interface {}]int{`,
					`-   "x": 20,`,
					`+   "x": 25,`,
					`    1: 10,`,
					`    true: 30,`,
					`  }`,
				}, "\n"),
			},
			{
				name:  "with nil key",
				left:  map[interface{}]int{nil: 1, "x": 2},
				right: map[interface{}]int{nil: 1, "y": 2},
				want: strings.Join([]string{
					`  map[interface {}]int{`,
					`-   "x": 2,`,
					`+   "y": 2,`,
					`    interface {}(<nil>): 1,`,
					`  }`,
				}, "\n"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				runTest(t, tc.left, tc.right, tc.want)
			})
		}
	})
}
