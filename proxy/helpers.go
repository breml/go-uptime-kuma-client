package proxy

import (
	"fmt"
	"reflect"
	"strings"
)

func formatProxy(p Proxy) string {
	buf := strings.Builder{}

	val := reflect.ValueOf(p)
	typ := reflect.TypeOf(p)

	first := true
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		if !field.IsExported() {
			continue
		}

		name := strings.Split(field.Tag.Get("json"), ",")[0]

		var valueStr string
		if field.Name == "Password" && value.Kind() == reflect.String && value.String() != "" {
			valueStr = "\"***\""
		} else if value.Kind() == reflect.String {
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
