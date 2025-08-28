package notification

import (
	"fmt"
	"maps"
	"strings"
)

type Generic struct {
	Base
	GenericDetails
	TypeName string
}

func (n Generic) Type() string {
	return n.TypeName
}

func (n Generic) String() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("%s: %q", "type", n.TypeName))

	for k, v := range orderedByKey(n.GenericDetails) {
		if k == "type" {
			continue
		}

		buf.WriteString(", ")

		str, ok := v.(string)
		if ok {
			buf.WriteString(fmt.Sprintf("%s: %q", k, str))
		} else {
			buf.WriteString(fmt.Sprintf("%s: %v", k, v))
		}
	}

	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), buf.String())
}

func (n *Generic) UnmarshalJSON(data []byte) error {
	details := GenericDetails{}
	base, err := unmarshalTo(data, &details)
	if err != nil {
		return err
	}

	typeNameAny, ok := details["type"]
	if !ok {
		return fmt.Errorf("notification does not have type attribute")
	}

	typeName, ok := typeNameAny.(string)
	if !ok {
		return fmt.Errorf("type attribute is not a string")
	}

	delete(details, "id")
	delete(details, "name")
	delete(details, "active")
	delete(details, "isDefault")
	delete(details, "applyExisting")
	delete(details, "userId")
	delete(details, "type")

	*n = Generic{
		Base:           base,
		GenericDetails: details,
		TypeName:       typeName,
	}

	return nil
}

func (n Generic) MarshalJSON() ([]byte, error) {
	details := maps.Clone(n.GenericDetails)
	details["type"] = n.TypeName

	return marshalJSON(n.Base, details)
}

type GenericDetails map[string]any

func (n GenericDetails) Type() string {
	tAny, ok := n["type"]
	if !ok {
		return ""
	}

	t, ok := tAny.(string)
	if !ok {
		return ""
	}

	return t
}
