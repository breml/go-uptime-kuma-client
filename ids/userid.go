// Package ids holds identifier types whose wire format differs between Uptime
// Kuma versions.
package ids

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// UserID is the ID of the account a resource belongs to, as the server reports
// it in its list events.
//
// An Uptime Kuma 2.x server numbers its users and sends the ID as a JSON
// number; a 3.x server, which manages its users with better-auth, gives them
// opaque string IDs and sends those. The type keeps whichever the server sent,
// so a list from either server decodes. A numeric ID goes back on the wire as
// the number it came as, see MarshalJSON, and Int64 returns it as the integer
// a 2.x server understands.
type UserID string

// Int64 returns the ID as an integer, or zero for an ID that is not a number,
// which is every ID an Uptime Kuma 3.x server hands out.
func (id UserID) Int64() int64 {
	value, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil {
		return 0
	}

	return value
}

// IsNumeric reports whether the ID is a number, which is what a 2.x server
// sends.
func (id UserID) IsNumeric() bool {
	_, err := strconv.ParseInt(string(id), 10, 64)

	return err == nil
}

// MarshalJSON writes a numeric ID as a JSON number and any other as a string,
// so an ID round-trips in the form the server sent it. The zero value is
// written as the number 0, which is what the integer field it replaces wrote.
func (id UserID) MarshalJSON() ([]byte, error) {
	if id == "" {
		return []byte("0"), nil
	}

	if id.IsNumeric() {
		return []byte(id), nil
	}

	data, err := json.Marshal(string(id))
	if err != nil {
		return nil, fmt.Errorf("marshal user id: %w", err)
	}

	return data, nil
}

// UnmarshalJSON reads a JSON number or string, and null as the zero value.
func (id *UserID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*id = ""

		return nil
	}

	if len(data) > 0 && data[0] == '"' {
		var value string

		err := json.Unmarshal(data, &value)
		if err != nil {
			return fmt.Errorf("unmarshal user id: %w", err)
		}

		*id = UserID(value)

		return nil
	}

	var value json.Number

	err := json.Unmarshal(data, &value)
	if err != nil {
		return fmt.Errorf("unmarshal user id: %w", err)
	}

	*id = UserID(value)

	return nil
}
