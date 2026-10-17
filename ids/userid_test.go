package ids_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/ids"
)

// TestUserIDRoundTrip covers the two wire forms: the number a 2.x server
// sends and the string a 3.x server sends, each written back as it came.
func TestUserIDRoundTrip(t *testing.T) {
	tests := map[string]struct {
		wire    string
		want    ids.UserID
		wantInt int64
		numeric bool
	}{
		"number": {wire: `1`, want: "1", wantInt: 1, numeric: true},
		"string": {wire: `"AbC123xyz"`, want: "AbC123xyz", wantInt: 0, numeric: false},
		"null":   {wire: `null`, want: "", wantInt: 0, numeric: false},
		"numeric string stays a string in memory but a number on the wire": {
			wire: `"7"`, want: "7", wantInt: 7, numeric: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var id ids.UserID

			require.NoError(t, json.Unmarshal([]byte(test.wire), &id))
			require.Equal(t, test.want, id)
			require.Equal(t, test.wantInt, id.Int64())
			require.Equal(t, test.numeric, id.IsNumeric())

			out, err := json.Marshal(id)
			require.NoError(t, err)

			switch {
			case test.wire == "null":
				require.Equal(t, "0", string(out), "the zero value writes what the integer field it replaces wrote")

			case test.numeric:
				require.Equal(t, string(test.want), string(out))

			default:
				require.Equal(t, test.wire, string(out))
			}
		})
	}
}

// TestUserIDInStruct covers the field in a list entry, which is where the
// type lives.
func TestUserIDInStruct(t *testing.T) {
	type entry struct {
		ID     int64      `json:"id"`
		UserID ids.UserID `json:"userId"`
	}

	var fromV2, fromV3 entry

	require.NoError(t, json.Unmarshal([]byte(`{"id":1,"userId":1}`), &fromV2))
	require.NoError(t, json.Unmarshal([]byte(`{"id":1,"userId":"kVqz8Z"}`), &fromV3))

	require.Equal(t, ids.UserID("1"), fromV2.UserID)
	require.Equal(t, ids.UserID("kVqz8Z"), fromV3.UserID)

	out, err := json.Marshal(fromV2)
	require.NoError(t, err)
	require.JSONEq(t, `{"id":1,"userId":1}`, string(out))

	out, err = json.Marshal(fromV3)
	require.NoError(t, err)
	require.JSONEq(t, `{"id":1,"userId":"kVqz8Z"}`, string(out))
}
