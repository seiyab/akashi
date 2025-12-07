package akashi

import (
	"reflect"

	"github.com/seiyab/akashi/internal/doc"
)

// Option is an option for DiffString.
type Option func(*differ) *differ

type ColorProfile = doc.ColorProfile

var DefaultColorProfile = ColorProfile{
	Left:    "\x1b[31m",
	Right:   "\x1b[32m",
	Neutral: "",
	Reset:   "\x1b[0m",
}

// WithReflectEqual sets a function to compare reflect.Value.
func WithReflectEqual(fn func(v1, v2 reflect.Value) bool) Option {
	return func(d *differ) *differ {
		d.reflectEqual = fn
		return d
	}
}

// WithFormat sets a function to format a value.
func WithFormat(fn interface{}) Option {
	return func(d *differ) *differ {
		fnValue := reflect.ValueOf(fn)
		fnType := fnValue.Type()
		if fnType.Kind() != reflect.Func || fnType.NumIn() != 1 || fnType.NumOut() != 1 {
			panic("WithFormat: fn must be a function with one input and one output")
		}
		if fnType.Out(0).Kind() != reflect.String {
			panic("WithFormat: fn must return a string")
		}

		if d.formats == nil {
			d.formats = make(map[reflect.Type]func(reflect.Value) string)
		}

		// Get the type from the function's input parameter
		inputType := fnType.In(0)
		d.formats[inputType] = func(v reflect.Value) string {
			return fnValue.Call([]reflect.Value{v})[0].String()
		}
		return d
	}
}

func WithColor(profile ...ColorProfile) Option {
	return func(d *differ) *differ {
		p := DefaultColorProfile
		if len(profile) > 0 {
			p = profile[0]
			if p.Reset == "" {
				p.Reset = DefaultColorProfile.Reset
			}
		}
		d.colorProfile = &p
		return d
	}
}
