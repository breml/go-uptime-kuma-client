package maintenance_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/maintenance"
)

func TestDateTime_UnmarshalJSON(t *testing.T) {
	plus2 := time.FixedZone("", 2*60*60)
	plus10 := time.FixedZone("", 10*60*60)

	tests := []struct {
		name    string
		input   string
		want    maintenance.DateTime
		wantErr []string
	}{
		{
			name:  "null",
			input: `null`,
		},
		{
			name:  "empty string",
			input: `""`,
		},
		{
			name:  "zero instant older clients wrote",
			input: `"0001-01-01T00:00:00Z"`,
			want:  maintenance.NewDateTime(time.Time{}),
		},
		{
			name:  "rfc3339 utc",
			input: `"2026-06-29T13:39:00Z"`,
			want:  maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
		},
		{
			name:  "rfc3339 with offset",
			input: `"2026-06-29T13:39:07+02:00"`,
			want:  maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 7, 0, plus2)),
		},
		{
			name:  "rfc3339 with fractional seconds",
			input: `"2026-06-29T13:39:07.123456789Z"`,
			want:  maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 7, 123456789, time.UTC)),
		},
		{
			name:  "rfc3339 with milliseconds as toISOString emits",
			input: `"2026-06-29T03:39:00.000Z"`,
			want:  maintenance.NewDateTime(time.Date(2026, 6, 29, 3, 39, 0, 0, time.UTC)),
		},
		{
			name:  "space separated with utc offset",
			input: `"2030-02-01 08:00:00Z"`,
			want:  maintenance.NewDateTime(time.Date(2030, 2, 1, 8, 0, 0, 0, time.UTC)),
		},
		{
			name:  "space separated with offset as python str of an aware datetime",
			input: `"2030-02-01 08:00:00+00:00"`,
			want:  maintenance.NewDateTime(time.Date(2030, 2, 1, 8, 0, 0, 0, time.UTC)),
		},
		{
			name:  "minutes with utc offset",
			input: `"2026-06-29T13:39Z"`,
			want:  maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
		},
		{
			name:  "minutes with offset",
			input: `"2026-06-29T13:39+10:00"`,
			want:  maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, plus10)),
		},
		{
			name:  "ui datetime-local without seconds",
			input: `"2026-06-29T13:39"`,
			want:  maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
		},
		{
			name:  "without offset with seconds",
			input: `"2026-06-29T13:39:07"`,
			want:  maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 7, 0, time.UTC)),
		},
		{
			name:  "without offset fractional seconds dropped",
			input: `"2026-06-29T13:39:07.5"`,
			want:  maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 7, 0, time.UTC)),
		},
		{
			name:  "without offset zero seconds and a fraction dropped",
			input: `"2026-06-29T13:39:00.5"`,
			want:  maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
		},
		{
			name:  "space separated with seconds as mariadb returns it",
			input: `"2026-06-29 13:39:07"`,
			want:  maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 7, 0, time.UTC)),
		},
		{
			name:  "space separated without seconds",
			input: `"2026-06-29 13:39"`,
			want:  maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
		},
		{
			name:  "date only",
			input: `"2026-06-29"`,
			want:  maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:  "floating year one",
			input: `"0001-01-01T00:00"`,
			want:  maintenance.NewFloatingDateTime(time.Time{}),
		},
		{
			name:  "date only year one",
			input: `"0001-01-01"`,
			want:  maintenance.NewFloatingDateTime(time.Time{}),
		},
		{
			name:    "garbage",
			input:   `"next tuesday"`,
			wantErr: []string{`"next tuesday"`},
		},
		{
			name:    "too short",
			input:   `"2026"`,
			wantErr: []string{`"2026"`},
		},
		{
			name:    "day out of range",
			input:   `"2026-06-31T13:39"`,
			wantErr: []string{`"2026-06-31T13:39"`, "day out of range"},
		},
		{
			name:    "hour out of range",
			input:   `"2026-06-29T24:00"`,
			wantErr: []string{`"2026-06-29T24:00"`, "hour out of range"},
		},
		{
			name:    "day out of range with offset",
			input:   `"2026-06-31T13:39:00Z"`,
			wantErr: []string{`"2026-06-31T13:39:00Z"`, "day out of range"},
		},
		// The server accepts whatever JavaScript's new Date(...) accepts, which
		// is more than the rows below. These formats are deliberately out of
		// scope for DateTime.
		{
			name:    "offset without colon",
			input:   `"2030-02-01T08:00:00+0000"`,
			wantErr: []string{`"2030-02-01T08:00:00+0000"`},
		},
		{
			name:    "lowercase separator",
			input:   `"2030-02-01t08:00:00Z"`,
			wantErr: []string{`"2030-02-01t08:00:00Z"`},
		},
		{
			name:    "lowercase utc designator",
			input:   `"2030-02-01T08:00:00z"`,
			wantErr: []string{`"2030-02-01T08:00:00z"`},
		},
		{
			name:    "slashes",
			input:   `"2030/02/01 08:00"`,
			wantErr: []string{`"2030/02/01 08:00"`},
		},
		{
			name:    "not a string",
			input:   `20260629`,
			wantErr: []string{"20260629"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := maintenance.NewFloatingDateTime(time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC))

			err := json.Unmarshal([]byte(tc.input), &got)
			if len(tc.wantErr) > 0 {
				require.Error(t, err)

				for _, want := range tc.wantErr {
					require.ErrorContains(t, err, want)
				}

				return
			}

			require.NoError(t, err)
			require.True(t, tc.want.Equal(got), "want %v, got %v", tc.want, got)
			require.Equal(t, tc.want.Floating(), got.Floating())
			require.Equal(t, tc.want.IsZero(), got.IsZero())

			_, wantOffset := tc.want.Time().Zone()
			_, gotOffset := got.Time().Zone()
			require.Equal(t, wantOffset, gotOffset)
		})
	}
}

func TestDateTime_UnmarshalReplacesFloatingTarget(t *testing.T) {
	floating := maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC))

	got := floating

	err := json.Unmarshal([]byte(`"2026-06-29T13:39:00Z"`), &got)
	require.NoError(t, err)
	require.False(t, got.Floating())
	require.True(t, maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)).Equal(got))

	got = floating

	err = json.Unmarshal([]byte(`""`), &got)
	require.NoError(t, err)
	require.True(t, got.IsZero())
	require.False(t, got.Floating())
}

func TestDateTime_MarshalJSON(t *testing.T) {
	zurich, err := time.LoadLocation("Europe/Zurich")
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   maintenance.DateTime
		want    string
		wantErr string
	}{
		{
			name:  "absent",
			input: maintenance.DateTime{},
			want:  `null`,
		},
		{
			name:  "zero time with offset as master wrote it",
			input: maintenance.NewDateTime(time.Time{}),
			want:  `"0001-01-01T00:00:00Z"`,
		},
		{
			name:  "zero time floating",
			input: maintenance.NewFloatingDateTime(time.Time{}),
			want:  `"0001-01-01T00:00"`,
		},
		{
			name:  "rfc3339 utc",
			input: maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
			want:  `"2026-06-29T13:39:00Z"`,
		},
		{
			name:  "rfc3339 with location",
			input: maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, zurich)),
			want:  `"2026-06-29T13:39:00+02:00"`,
		},
		{
			name:  "rfc3339 with fractional seconds",
			input: maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 7, 120000000, time.UTC)),
			want:  `"2026-06-29T13:39:07.12Z"`,
		},
		{
			name:  "floating without seconds",
			input: maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
			want:  `"2026-06-29T13:39"`,
		},
		{
			name:  "floating with seconds",
			input: maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 7, 0, time.UTC)),
			want:  `"2026-06-29T13:39:07"`,
		},
		{
			name:  "floating fraction dropped on construct",
			input: maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 7, 500000000, time.UTC)),
			want:  `"2026-06-29T13:39:07"`,
		},
		{
			name:  "floating zero seconds and a fraction dropped on construct",
			input: maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 500000000, time.UTC)),
			want:  `"2026-06-29T13:39"`,
		},
		{
			name:  "floating keeps the wall clock digits of another location",
			input: maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, zurich)),
			want:  `"2026-06-29T13:39"`,
		},
		{
			name:    "floating year after 9999",
			input:   maintenance.NewFloatingDateTime(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: "year outside of range",
		},
		{
			name:    "floating year before 0",
			input:   maintenance.NewFloatingDateTime(time.Date(-1, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: "year outside of range",
		},
		{
			name:    "rfc3339 year after 9999",
			input:   maintenance.NewDateTime(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: "year outside of range",
		},
		{
			name:    "rfc3339 offset of 25 hours",
			input:   maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.FixedZone("", 25*60*60))),
			wantErr: "timezone hour outside of range",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.input)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			require.JSONEq(t, tc.want, string(got))
		})
	}
}

func TestDateTime_RoundTripKeepsKind(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: `"2026-06-29T13:39"`, want: `"2026-06-29T13:39"`},
		{input: `"2026-06-29T13:39:07"`, want: `"2026-06-29T13:39:07"`},
		{input: `"2026-06-29T13:39:07.5"`, want: `"2026-06-29T13:39:07"`},
		{input: `"2026-06-29T13:39:00"`, want: `"2026-06-29T13:39"`},
		{input: `"2026-06-29 13:39:00"`, want: `"2026-06-29T13:39"`},
		{input: `"2026-06-29"`, want: `"2026-06-29T00:00"`},
		{input: `"0001-01-01T00:00"`, want: `"0001-01-01T00:00"`},
		{input: `"0001-01-01T00:00:00Z"`, want: `"0001-01-01T00:00:00Z"`},
		{input: `"2026-06-29T13:39:00Z"`, want: `"2026-06-29T13:39:00Z"`},
		{input: `"2026-06-29T13:39:00.000Z"`, want: `"2026-06-29T13:39:00Z"`},
		{input: `"2030-02-01 08:00:00Z"`, want: `"2030-02-01T08:00:00Z"`},
		{input: `"2030-02-01 08:00:00+00:00"`, want: `"2030-02-01T08:00:00Z"`},
		{input: `"2026-06-29T13:39+10:00"`, want: `"2026-06-29T13:39:00+10:00"`},
		{input: `"2026-06-29T13:39:07+10:00"`, want: `"2026-06-29T13:39:07+10:00"`},
		{input: `"2026-06-29T13:39:07.123456789-05:30"`, want: `"2026-06-29T13:39:07.123456789-05:30"`},
		{input: `null`, want: `null`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			var first maintenance.DateTime

			err := json.Unmarshal([]byte(tc.input), &first)
			require.NoError(t, err)

			encoded, err := json.Marshal(first)
			require.NoError(t, err)
			require.JSONEq(t, tc.want, string(encoded))

			var second maintenance.DateTime

			err = json.Unmarshal(encoded, &second)
			require.NoError(t, err)
			require.True(t, first.Equal(second), "first %v, second %v", first, second)
		})
	}
}

// TestDateTime_MapKeyTextEncoding pins that a floating map key is encoded
// through MarshalText in the floating form and decoded back as floating.
func TestDateTime_MapKeyTextEncoding(t *testing.T) {
	key := maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC))

	encoded, err := json.Marshal(map[maintenance.DateTime]int{key: 1})
	require.NoError(t, err)
	require.JSONEq(t, `{"2026-06-29T13:39":1}`, string(encoded))

	var decoded map[maintenance.DateTime]int

	err = json.Unmarshal(encoded, &decoded)
	require.NoError(t, err)
	require.Len(t, decoded, 1)

	for got := range decoded {
		require.True(t, key.Equal(got), "got %v", got)
	}
}

func TestDateTime_TextRoundTrip(t *testing.T) {
	minus5 := time.FixedZone("", -5*60*60)

	for _, input := range []maintenance.DateTime{
		{},
		maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
		maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, minus5)),
	} {
		text, err := input.MarshalText()
		require.NoError(t, err)

		var got maintenance.DateTime

		err = got.UnmarshalText(text)
		require.NoError(t, err)
		require.True(t, input.Equal(got), "want %v, got %v", input, got)

		_, wantOffset := input.Time().Zone()
		_, gotOffset := got.Time().Zone()
		require.Equal(t, wantOffset, gotOffset)
	}

	text, err := maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, minus5)).MarshalText()
	require.NoError(t, err)
	require.Equal(t, "2026-06-29T13:39:00-05:00", string(text))
}

func TestNewFloatingDateTime(t *testing.T) {
	brisbane, err := time.LoadLocation("Australia/Brisbane")
	require.NoError(t, err)

	got := maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 7, 5, brisbane))

	require.True(t, got.Floating())
	require.Equal(t, time.Date(2026, 6, 29, 13, 39, 7, 0, time.UTC), got.Time())
}

func TestNewFloatingDateTime_AlwaysPresent(t *testing.T) {
	tests := map[string]time.Time{
		"zero time":        {},
		"year one non-utc": time.Date(1, 1, 1, 0, 0, 0, 0, time.FixedZone("", 2*60*60)),
	}

	for name, in := range tests {
		t.Run(name, func(t *testing.T) {
			got := maintenance.NewFloatingDateTime(in)

			require.False(t, got.IsZero())
			require.True(t, got.Floating())
			require.False(t, maintenance.DateTime{}.Equal(got))
			require.Equal(t, "0001-01-01T00:00", got.String())
		})
	}
}

func TestNewDateTime(t *testing.T) {
	in := time.Date(2026, 6, 29, 13, 39, 7, 0, time.FixedZone("", 10*60*60))

	got := maintenance.NewDateTime(in)

	require.False(t, got.Floating())
	require.False(t, got.IsZero())
	require.True(t, in.Equal(got.Time()))

	zero := maintenance.NewDateTime(time.Time{})
	require.False(t, zero.IsZero())
	require.False(t, zero.Floating())
	require.False(t, maintenance.DateTime{}.Equal(zero))
}

func TestNewDateTime_DropsMonotonicReading(t *testing.T) {
	now := time.Now()
	require.Contains(t, now.String(), "m=")

	got := maintenance.NewDateTime(now)
	require.NotContains(t, got.Time().String(), "m=")
	require.True(t, now.Equal(got.Time()))
}

func TestDateTime_In(t *testing.T) {
	brisbane, err := time.LoadLocation("Australia/Brisbane")
	require.NoError(t, err)

	floating := maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC))
	instant := maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC))

	tests := []struct {
		name         string
		in           maintenance.DateTime
		loc          *time.Location
		want         time.Time
		wantLocation *time.Location
	}{
		{
			name:         "floating in brisbane",
			in:           floating,
			loc:          brisbane,
			want:         time.Date(2026, 6, 29, 3, 39, 0, 0, time.UTC),
			wantLocation: brisbane,
		},
		{
			name:         "floating in nil location",
			in:           floating,
			want:         time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC),
			wantLocation: time.UTC,
		},
		{
			name:         "offset in brisbane",
			in:           instant,
			loc:          brisbane,
			want:         time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC),
			wantLocation: brisbane,
		},
		{
			name:         "offset in nil location",
			in:           instant,
			want:         time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC),
			wantLocation: time.UTC,
		},
		{
			name:         "floating year one in brisbane",
			in:           maintenance.NewFloatingDateTime(time.Time{}),
			loc:          brisbane,
			want:         time.Date(1, 1, 1, 0, 0, 0, 0, brisbane),
			wantLocation: brisbane,
		},
		{
			name:         "offset year one in brisbane",
			in:           maintenance.NewDateTime(time.Time{}),
			loc:          brisbane,
			want:         time.Time{},
			wantLocation: brisbane,
		},
		{
			name:         "absent",
			in:           maintenance.DateTime{},
			loc:          brisbane,
			want:         time.Time{},
			wantLocation: time.UTC,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.In(tc.loc)

			require.True(t, tc.want.Equal(got), "want %v, got %v", tc.want, got)
			require.Equal(t, tc.wantLocation, got.Location())
		})
	}

	require.Equal(t, 23, instant.In(brisbane).Hour())
	require.Equal(t, 13, floating.In(brisbane).Hour())
}

func TestDateTime_String(t *testing.T) {
	tests := []struct {
		name string
		in   maintenance.DateTime
		want string
	}{
		{name: "absent", in: maintenance.DateTime{}, want: ""},
		{
			name: "floating",
			in:   maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)),
			want: "2026-06-29T13:39",
		},
		{
			name: "offset",
			in:   maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.FixedZone("", 10*60*60))),
			want: "2026-06-29T13:39:00+10:00",
		},
		{
			name: "zero time with offset",
			in:   maintenance.NewDateTime(time.Time{}),
			want: "0001-01-01T00:00:00Z",
		},
		{
			name: "floating year after 9999",
			in:   maintenance.NewFloatingDateTime(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)),
			want: "10000-01-01T00:00",
		},
		{
			name: "offset of 25 hours",
			in:   maintenance.NewDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.FixedZone("", 25*60*60))),
			want: "2026-06-29T13:39:00+25:00",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, tc.in.String())
		})
	}
}

func TestDateTime_Equal(t *testing.T) {
	digits := time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC)

	floating := maintenance.NewFloatingDateTime(digits)
	instant := maintenance.NewDateTime(digits)

	require.True(t, floating.Equal(maintenance.NewFloatingDateTime(digits)))
	require.True(t, instant.Equal(maintenance.NewDateTime(digits.In(time.FixedZone("", 2*60*60)))))
	require.False(t, floating.Equal(instant))
	require.False(t, instant.Equal(floating))
	require.False(t, floating.Equal(maintenance.DateTime{}))

	var absent maintenance.DateTime

	require.True(t, absent.Equal(maintenance.DateTime{}))
}

func TestDateTime_EqualNotDoubleEquals(t *testing.T) {
	var utcZ, utcOffset maintenance.DateTime

	require.NoError(t, json.Unmarshal([]byte(`"2030-02-01T08:00:00Z"`), &utcZ))
	require.NoError(t, json.Unmarshal([]byte(`"2030-02-01T08:00:00+00:00"`), &utcOffset))

	require.True(t, utcZ.Equal(utcOffset))
	require.NotEqual(t, utcZ, utcOffset, "the values differ in their location")
}
