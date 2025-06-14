package akashi_test

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/seiyab/akashi"
)

func TestDiff_String(t *testing.T) {
	t.Run("word", func(t *testing.T) {
		runTest(t, "hello", "world", strings.Join([]string{
			`- "hello"`,
			`+ "world"`,
		}, "\n"))
	})
}

func TestDiff_Primitive(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		runTest(t, 1, 2, strings.Join([]string{
			`- 1`,
			`+ 2`,
		}, "\n"))
	})

	t.Run("float", func(t *testing.T) {
		runTest(t, 1.0, 2.0, strings.Join([]string{
			`- 1.000000`,
			`+ 2.000000`,
		}, "\n"))
	})

	t.Run("bool", func(t *testing.T) {
		runTest(t, true, false, strings.Join([]string{
			`- true`,
			`+ false`,
		}, "\n"))
	})
}

func TestDiff_Func(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		f1 := func(*testing.T) {}
		f2 := func(*testing.T) {}
		// NOTE: reflect.DeepEqual cannot compare functions.
		runTest(t, f1, f2, strings.Join([]string{
			fmt.Sprintf(`- func(*testing.T) { ... } at [%p]`, f1),
			fmt.Sprintf(`+ func(*testing.T) { ... } at [%p]`, f2),
		}, "\n"))
	})

	type f = func()
	t.Run("both nil", func(t *testing.T) {
		runTest(t, f(nil), f(nil), "")
	})

	t.Run("left nil", func(t *testing.T) {
		t.Skip()
		expected := strings.Join([]string{
			`- func(t *testing.T) { ... }`,
			`+ func()(nil)`,
		}, "\n")
		runTest(t, TestDiff_Primitive, f(nil), expected)
	})
}

func TestDiff_Pointer(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		runTest(t, (*int)(nil), (*int)(nil), "")
	})

	t.Run("same", func(t *testing.T) {
		runTest(t, refInt(100), refInt(100), "")
	})

	t.Run("different", func(t *testing.T) {
		expected := strings.Join([]string{
			`- &100`,
			`+ &200`,
		}, "\n")
		runTest(t, refInt(100), refInt(200), expected)
	})
}

func refInt(x int) *int {
	return &x
}

func TestDiff_Interface(t *testing.T) {
	type S struct {
		X fmt.Stringer
		Y io.Closer
	}

	t.Run("same", func(t *testing.T) {
	})

	t.Run("different Stringer", func(t *testing.T) {
		runTest(t,
			S{X: I(1)},
			S{X: I(2)},
			strings.Join([]string{
				`  akashi_test.S{`,
				`-   X: akashi_test.I("1"),`,
				`+   X: akashi_test.I("2"),`,
				`    Y: io.Closer(nil),`,
				`  }`,
			}, "\n"))
	})

	t.Run("different", func(t *testing.T) {
		runTest(t,
			S{X: I(1), Y: C(1)},
			S{X: I(1), Y: C(2)},
			strings.Join([]string{
				`  akashi_test.S{`,
				`    X: akashi_test.I("1"),`,
				`-   Y: 1,`,
				`+   Y: 2,`,
				`  }`,
			}, "\n"))
	})

	t.Run("nil", func(t *testing.T) {
		runTest(t,
			S{X: I(1)},
			S{X: nil},
			strings.Join([]string{
				`  akashi_test.S{`,
				`-   X: akashi_test.I("1"),`,
				`+   X: fmt.Stringer(nil),`,
				`    Y: io.Closer(nil),`,
				`  }`,
			}, "\n"))
	})
}

type I int

func (v I) String() string { return strconv.Itoa(int(v)) }

type C int

func (v C) Close() error { return nil }

func TestDiff_Chan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		type s struct {
			c chan bool
		}
		c := make(chan bool)
		runTest(t, s{c: nil}, s{c: c}, strings.Join([]string{
			`  akashi_test.s{`,
			`-   c: chan bool(nil),`,
			fmt.Sprintf(`+   c: chan bool at [%p],`, c),
			`  }`,
		}, "\n"))
	})

	t.Run("different", func(t *testing.T) {
		c1 := make(<-chan string)
		c2 := make(<-chan string)
		expected := strings.Join([]string{
			fmt.Sprintf(`- <-chan string at [%p]`, c1),
			fmt.Sprintf(`+ <-chan string at [%p]`, c2),
		}, "\n")
		runTest(t, c1, c2, expected)
	})
}

func TestDiff_TypeMismatch(t *testing.T) {
	t.Run("struct vs struct", func(t *testing.T) {
		type s struct {
			i int
		}
		type u struct {
			i int
		}
		runTest(t, s{1}, u{1}, strings.Join([]string{
			`- akashi_test.s{`,
			`-   i: 1,`,
			`- }`,
			`+ akashi_test.u{`,
			`+   i: 1,`,
			`+ }`,
		}, "\n"))
	})

	t.Run("string vs map", func(t *testing.T) {
		runTest(t, "hello", map[string]int{}, strings.Join([]string{
			`- "hello"`,
			`+ map[string]int{`,
			`+ }`,
		}, "\n"))
	})
}

func runTest(t *testing.T, left, right interface{}, want string, opts ...akashi.Option) {
	t.Helper()
	d := akashi.DiffString(left, right, opts...)
	if d != want {
		t.Errorf("expected %q, got %q", want, d)
		p := akashi.DiffString(want, d)
		t.Log("\n" + p)
	}
}
