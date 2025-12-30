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

// formatNotification formats a notification instance as a string representation.
// If includeType is true, it includes the notification type in the output.
//nolint:revive // includeType is not a control coupling flag, it's a meaningful parameter
func formatNotification(s any, includeType bool) string {
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

// unmarshalTo unmarshals notification data into a Base and a detail struct.
// It parses both the base fields and the nested configuration data.
func unmarshalTo(data []byte, detail any) (Base, error) {
	notificationBase := Base{}

	err := json.Unmarshal(data, &notificationBase)
	if err != nil {
		return notificationBase, fmt.Errorf("unmarshal notification: %w", err)
	}

	err = json.Unmarshal([]byte(notificationBase.configStr), detail)
	if err != nil {
		return notificationBase, fmt.Errorf("unmarshal notification config: %w", err)
	}

	return notificationBase, nil
}

// marshalJSON marshals a notification into JSON bytes.
// It combines base notification fields with type-specific details.
func marshalJSON(base Base, details interface{ Type() string }) ([]byte, error) {
	detailData, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("marshal notification details: %w", err)
	}

	rawMap := make(map[string]any, 20)
	rawMap["id"] = base.ID
	rawMap["name"] = base.Name
	rawMap["active"] = base.IsActive
	rawMap["isDefault"] = base.IsDefault
	rawMap["applyExisting"] = base.ApplyExisting
	rawMap["userId"] = base.UserID
	rawMap["type"] = base.typeFromConfigStr

	err = json.Unmarshal(detailData, &rawMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal notification details: %w", err)
	}

	rawMap["type"] = details.Type()

	data, err := json.Marshal(rawMap)
	if err != nil {
		return nil, fmt.Errorf("marshal notification: %w", err)
	}
	return data, nil
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
