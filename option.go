package akashi

import (
	"reflect"
)

// Option is an option for DiffString.
type Option func(*differ) *differ

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
