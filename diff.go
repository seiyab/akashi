package akashi

import (
	"fmt"
	"reflect"
	"unsafe"
)

// DiffString returns a string that represents the difference between x and y.
func DiffString(x, y interface{}, options ...Option) string {
	d := differ{}
	for _, opt := range options {
		d = *opt(&d)
	}
	t := d.diff(x, y)
	return t.format()
}

type differ struct {
	reflectEqual func(v1, v2 reflect.Value) bool
	formats      formats
}

type formats map[reflect.Type]func(reflect.Value) string

func (d differ) diff(x, y interface{}) root {
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	p := diffProcess{differ: d}
	t := p.diff(v1, v2)
	return root{inner: t}
}

type diffProcess struct {
	differ       differ
	depth        int
	leftVisited  map[visit]bool
	rightVisited map[visit]bool
	pureVisited  map[visit]bool
}

type visit struct {
	ptr unsafe.Pointer
	typ reflect.Type
}

const maxDepth = 500

func (p diffProcess) diff(v1, v2 reflect.Value) (result diffTree) {
	if p.depth > maxDepth {
		return fail{difference: 1, message: "max depth exceeded"}
	}
	p.depth = p.depth + 1

	defer func() {
		if r := recover(); r != nil {
			result = fail{difference: 1, message: fmt.Sprintf("panic: %v", r)}
		}
	}()

	d := p.differ
	if d.reflectEqual != nil {
		if d.reflectEqual(v1, v2) {
			return p.pure(v1)
		}
	} else if lightDeepEqual(v1, v2) {
		return p.pure(v1)
	}
	if !v1.IsValid() || !v2.IsValid() {
		return fail{difference: 1, message: "invalid value"}
	}
	if v1.Type() != v2.Type() {
		return p.eachSide(v1, v2)
	}

	var cyclic bool
	p, cyclic = p.cycle(v1, v2)
	if cyclic {
		return split{
			left:  cycle{},
			right: cycle{},
		}
	}

	diffFunc, ok := diffFuncs[v1.Kind()]
	if !ok {
		return fail{difference: 1, message: fmt.Sprintf("unsopported kind: %s", v1.Kind())}
	}
	t := diffFunc(v1, v2, p)
	if f, ok := d.formats[v1.Type()]; ok {
		t = format2{left: v1, right: v2, original: t, format: f}
	} else if v1.Type().Implements(stringerType) {
		t = format2{left: v1, right: v2, original: t}
	}
	return t
}

func lightDeepEqual(v1 reflect.Value, v2 reflect.Value) bool {
	if v1.Type() != v2.Type() {
		return false
	}
	if v1.CanInterface() && v2.CanInterface() {
		return reflect.DeepEqual(v1.Interface(), v2.Interface())
	}
	if v1.CanAddr() && v2.CanAddr() && v1.Addr().Pointer() == v2.Addr().Pointer() {
		return true
	}
	return false // can't go better until go 1.20
}

func (p diffProcess) cycle(v1 reflect.Value, v2 reflect.Value) (diffProcess, bool) {
	if !hard(v1) && !hard(v2) {
		return p, false
	}
	leftCycle := false
	rightCycle := false
	p = p.clone()
	if hard(v1) && v1.CanAddr() {
		addr := v1.Addr().UnsafePointer()
		vis := visit{ptr: addr, typ: v1.Type()}
		leftCycle = p.leftVisited[vis]
		p.leftVisited[vis] = true
	}
	if hard(v2) && v2.CanAddr() {
		addr := v2.Addr().UnsafePointer()
		vis := visit{ptr: addr, typ: v2.Type()}
		rightCycle = p.rightVisited[vis]
		p.rightVisited[vis] = true
	}

	return p, leftCycle && rightCycle
}

func (p diffProcess) clone() diffProcess {
	return diffProcess{
		differ:       p.differ,
		leftVisited:  cloneVisits(p.leftVisited),
		rightVisited: cloneVisits(p.rightVisited),
		pureVisited:  cloneVisits(p.pureVisited),
	}
}

func hard(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Interface:
		return !v.IsNil()
	case reflect.Struct, reflect.Array:
		return true
	}
	return false
}
