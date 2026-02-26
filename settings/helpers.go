package settings

import (
	"fmt"
	"reflect"
	"strings"
)

// formatSettings formats a Settings instance as a string representation.
// It masks the SteamAPIKey field for security purposes before including in the output.
// The function uses reflection to iterate through exported fields and builds a string representation.
func formatSettings(s Settings) string {
	buf := strings.Builder{}

	val := reflect.ValueOf(s)
	typ := reflect.TypeFor[Settings]()

	first := true
	for i := range val.NumField() {
		field := typ.Field(i)
		value := val.Field(i)

		if !field.IsExported() {
			continue
		}

		name := strings.Split(field.Tag.Get("json"), ",")[0]

		var valueStr string
		switch {
		case field.Name == "SteamAPIKey" && value.Kind() == reflect.String && value.String() != "":
			valueStr = "\"***\""

		case value.Kind() == reflect.String:
			valueStr = fmt.Sprintf("%q", value.String())

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
