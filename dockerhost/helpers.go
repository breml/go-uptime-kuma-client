package dockerhost

import (
	"fmt"
	"reflect"
	"strings"
)

// formatDockerHost formats a DockerHost instance as a string representation.
// It uses reflection to iterate through exported fields and builds a
// comma-separated string of field names and values for display.
func formatDockerHost(d DockerHost) string {
	buf := strings.Builder{}

	val := reflect.ValueOf(d)
	typ := reflect.TypeFor[DockerHost]()

	first := true
	for i := range val.NumField() {
		field := typ.Field(i)
		value := val.Field(i)

		if !field.IsExported() {
			continue
		}

		name := strings.Split(field.Tag.Get("json"), ",")[0]

		var valueStr string
		if value.Kind() == reflect.String {
			valueStr = fmt.Sprintf("%q", value.String())
		} else {
			valueStr = fmt.Sprintf("%v", value.Interface())
		}

		if !first {
			buf.WriteString(", ")
		}

		first = false

		_, _ = buf.WriteString(fmt.Sprintf("%s: %s", name, valueStr))
	}

	return buf.String()
}
