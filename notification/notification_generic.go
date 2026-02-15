package notification

import (
	"errors"
	"fmt"
	"maps"
	"strings"
)

// Generic represents a generic notification.
type Generic struct {
	Base
	GenericDetails

	TypeName string
}

// Type returns the notification type.
func (n Generic) Type() string {
	return n.TypeName
}

// String returns a string representation of the notification.
func (n Generic) String() string {
	buf := strings.Builder{}
	fmt.Fprintf(&buf, "%s: %q", "type", n.TypeName)

	for k, v := range orderedByKey(n.GenericDetails) {
		if k == "type" {
			continue
		}

		buf.WriteString(", ")

		str, ok := v.(string)
		if ok {
			fmt.Fprintf(&buf, "%s: %q", k, str)
		} else {
			fmt.Fprintf(&buf, "%s: %v", k, v)
		}
	}

	return fmt.Sprintf("%s, %s", formatNotification(n.Base, false), buf.String())
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (n *Generic) UnmarshalJSON(data []byte) error {
	details := GenericDetails{}
	base, err := unmarshalTo(data, &details)
	if err != nil {
		return err
	}

	typeNameAny, ok := details["type"]
	if !ok {
		return errors.New("notification does not have type attribute")
	}

	typeName, ok := typeNameAny.(string)
	if !ok {
		return errors.New("type attribute is not a string")
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

// MarshalJSON marshals a notification into a JSON byte slice.
func (n Generic) MarshalJSON() ([]byte, error) {
	details := maps.Clone(n.GenericDetails)
	details["type"] = n.TypeName

	return marshalJSON(n.Base, details)
}

// GenericDetails represents generic notification configuration details.
type GenericDetails map[string]any

// Type returns the notification type.
func (n GenericDetails) Type() string {
	if n == nil {
		return ""
	}

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
