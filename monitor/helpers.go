package monitor

import (
	"cmp"
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"
)

// formatMonitor formats a monitor instance as a string representation.
// If includeType is true, it includes the monitor type in the output.
func formatMonitor( //nolint:revive // includeType is not a control coupling flag, it's a meaningful parameter
	s any,
	includeType bool,
) string {
	buf := strings.Builder{}

	first := true
	if includeType {
		typer, ok := s.(interface{ Type() string })
		if ok {
			buf.WriteString("type: " + typer.Type())
			first = false
		}
	}

	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// handle pointer
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := range val.NumField() {
		field := typ.Field(i)
		value := val.Field(i)

		if !field.IsExported() {
			continue
		}

		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")

		if name == "" || name == "-" {
			name = field.Name
		}

		// Render pointer fields by dereferencing so optional values such as
		// *string appear as their underlying value (or <nil>) instead of as
		// a pointer address.
		var valueStr string
		switch {
		case value.Kind() == reflect.String:
			valueStr = fmt.Sprintf("%q", value.String())

		case value.Kind() == reflect.Pointer && value.IsNil():
			valueStr = "<nil>"

		case value.Kind() == reflect.Pointer && value.Elem().Kind() == reflect.String:
			valueStr = fmt.Sprintf("%q", value.Elem().String())

		case value.Kind() == reflect.Pointer:
			valueStr = fmt.Sprintf("%v", value.Elem().Interface())

		default:
			valueStr = fmt.Sprintf("%v", value.Interface())
		}

		if !first {
			buf.WriteString(", ")
		}

		first = false

		_, _ = fmt.Fprintf(&buf, "%s: %s", name, valueStr)
	}

	return buf.String()
}

// orderedByKey returns an iterator over a map's key-value pairs in sorted order.
// Keys must be of a comparable ordered type.
func orderedByKey[K cmp.Ordered, E any](m map[K]E) iter.Seq2[K, E] {
	return func(yield func(K, E) bool) {
		keys := make([]K, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}

		slices.Sort(keys)

		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}
