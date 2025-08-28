package notification

import (
	"cmp"
	"encoding/json"
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"
)

func formatNotification(s any, withType bool) string {
	buf := strings.Builder{}

	first := true
	if withType {
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

func unmarshalTo(data []byte, detail any) (Base, error) {
	notificationBase := Base{}

	err := json.Unmarshal(data, &notificationBase)
	if err != nil {
		return notificationBase, err
	}

	err = json.Unmarshal([]byte(notificationBase.configStr), detail)
	if err != nil {
		return notificationBase, err
	}

	return notificationBase, nil
}

func marshalJSON(base Base, details interface{ Type() string }) ([]byte, error) {
	baseData, err := json.Marshal(base)
	if err != nil {
		return nil, err
	}

	detailData, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	rawMap := make(map[string]any, 20)
	err = json.Unmarshal(baseData, &rawMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(detailData, &rawMap)
	if err != nil {
		return nil, err
	}

	rawMap["type"] = details.Type()

	return json.Marshal(rawMap)
}

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
