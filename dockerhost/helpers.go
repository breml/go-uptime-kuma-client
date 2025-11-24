package dockerhost

import (
	"fmt"
	"reflect"
	"strings"
)

func formatDockerHost(d DockerHost) string {
	buf := strings.Builder{}

	val := reflect.ValueOf(d)
	typ := reflect.TypeOf(d)

	first := true
	for i := 0; i < val.NumField(); i++ {
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
