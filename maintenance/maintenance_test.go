package maintenance_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/maintenance"
)

// uiCreatedMaintenanceList is the maintenanceList payload of a window created
// in the Uptime Kuma 2.5.5 UI, as reported in issue #401.
const uiCreatedMaintenanceList = `{
  "1": {
    "id": 1, "title": "example", "strategy": "single", "active": true,
    "dateRange": ["2026-06-29T13:39", "2026-07-12T14:39"],
    "timeslotList": [{"startDate": "2026-06-29T13:39", "endDate": "2026-07-12T14:39"}],
    "timezone": "Australia/Brisbane", "timezoneOption": "SAME_AS_SERVER",
    "timezoneOffset": "+10:00", "status": "ended"
  }
}`

func TestMaintenance_UnmarshalUICreatedList(t *testing.T) {
	var list map[string]maintenance.Maintenance

	err := json.Unmarshal([]byte(uiCreatedMaintenanceList), &list)
	require.NoError(t, err)
	require.Len(t, list, 1)

	m := list["1"]
	require.Equal(t, int64(1), m.ID)
	require.Equal(t, "example", m.Title)
	require.Equal(t, "single", m.Strategy)
	require.Equal(t, "Australia/Brisbane", m.Timezone)

	start := maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC))
	end := maintenance.NewFloatingDateTime(time.Date(2026, 7, 12, 14, 39, 0, 0, time.UTC))

	requireDates(t, []*maintenance.DateTime{&start, &end}, m.DateRange)

	require.Len(t, m.TimeslotList, 1)
	require.True(t, start.Equal(m.TimeslotList[0].StartDate), "start: got %v", m.TimeslotList[0].StartDate)
	require.True(t, end.Equal(m.TimeslotList[0].EndDate), "end: got %v", m.TimeslotList[0].EndDate)

	encoded, err := json.Marshal(m)
	require.NoError(t, err)

	var dates struct {
		DateRange    json.RawMessage `json:"dateRange"`
		TimeslotList json.RawMessage `json:"timeslotList"`
	}

	err = json.Unmarshal(encoded, &dates)
	require.NoError(t, err)
	require.JSONEq(t, `["2026-06-29T13:39","2026-07-12T14:39"]`, string(dates.DateRange))
	require.JSONEq(t, `[{"startDate":"2026-06-29T13:39","endDate":"2026-07-12T14:39"}]`, string(dates.TimeslotList))
}

func TestMaintenance_UnmarshalDateVariants(t *testing.T) {
	floating := maintenance.NewFloatingDateTime(time.Date(2026, 6, 29, 13, 39, 0, 0, time.UTC))
	cronStart := maintenance.NewDateTime(time.Date(2026, 6, 29, 3, 39, 0, 0, time.UTC))
	cronEnd := maintenance.NewDateTime(time.Date(2026, 6, 29, 5, 39, 0, 0, time.UTC))

	tests := []struct {
		name          string
		input         string
		wantDateRange []*maintenance.DateTime
		wantTimeslots []maintenance.Timeslot
	}{
		{
			name:          "nil elements",
			input:         `{"dateRange":[null,null]}`,
			wantDateRange: []*maintenance.DateTime{nil, nil},
		},
		{
			name:          "single element",
			input:         `{"dateRange":["2026-06-29T13:39"]}`,
			wantDateRange: []*maintenance.DateTime{&floating},
		},
		{
			name: "cron timeslots as toISOString emits them",
			input: `{"dateRange":[null,null],"timeslotList":[` +
				`{"startDate":"2026-06-29T03:39:00.000Z","endDate":"2026-06-29T05:39:00.000Z"}]}`,
			wantDateRange: []*maintenance.DateTime{nil, nil},
			wantTimeslots: []maintenance.Timeslot{{StartDate: cronStart, EndDate: cronEnd}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var m maintenance.Maintenance

			err := json.Unmarshal([]byte(tc.input), &m)
			require.NoError(t, err)
			requireDates(t, tc.wantDateRange, m.DateRange)
			require.Len(t, m.TimeslotList, len(tc.wantTimeslots))

			for i, want := range tc.wantTimeslots {
				require.True(
					t,
					want.StartDate.Equal(m.TimeslotList[i].StartDate),
					"got %v",
					m.TimeslotList[i].StartDate,
				)
				require.True(t, want.EndDate.Equal(m.TimeslotList[i].EndDate), "got %v", m.TimeslotList[i].EndDate)
				require.False(t, m.TimeslotList[i].StartDate.Floating())
			}
		})
	}
}

func TestHelpers_DateRangeWireFormat(t *testing.T) {
	plus2 := time.FixedZone("", 2*60*60)
	start := time.Date(2030, 2, 1, 8, 0, 0, 0, plus2)
	end := time.Date(2030, 2, 1, 10, 0, 0, 0, plus2)
	timeRange := []maintenance.TimeOfDay{{Hours: 2}, {Hours: 4}}

	tests := []struct {
		name string
		m    *maintenance.Maintenance
		want string
	}{
		{
			name: "single",
			m:    maintenance.NewSingleMaintenance("t", "d", start, end, "UTC"),
			want: `["2030-02-01T08:00:00+02:00","2030-02-01T10:00:00+02:00"]`,
		},
		{
			name: "single with zero start as master sent it",
			m:    maintenance.NewSingleMaintenance("t", "d", time.Time{}, end, "UTC"),
			want: `["0001-01-01T00:00:00Z","2030-02-01T10:00:00+02:00"]`,
		},
		{
			name: "manual",
			m:    maintenance.NewManualMaintenance("t", "d"),
			want: `[null,null]`,
		},
		{
			name: "cron",
			m:    maintenance.NewCronMaintenance("t", "d", "0 2 * * *", 30, "UTC"),
			want: `[null,null]`,
		},
		{
			name: "recurring weekday",
			m:    maintenance.NewRecurringWeekdayMaintenance("t", "d", []int{1}, timeRange, "UTC"),
			want: `[null,null]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.m)
			require.NoError(t, err)

			var wire struct {
				DateRange json.RawMessage `json:"dateRange"`
			}

			err = json.Unmarshal(encoded, &wire)
			require.NoError(t, err)
			require.JSONEq(t, tc.want, string(wire.DateRange))
		})
	}
}

func TestMaintenance_ZeroInstantRoundTrip(t *testing.T) {
	input := `{"dateRange":["0001-01-01T00:00:00Z","2030-02-01T10:00:00Z"]}`

	var m maintenance.Maintenance

	err := json.Unmarshal([]byte(input), &m)
	require.NoError(t, err)
	require.Len(t, m.DateRange, 2)
	require.NotNil(t, m.DateRange[0])
	require.False(t, m.DateRange[0].IsZero())
	require.False(t, m.DateRange[0].Floating())

	encoded, err := json.Marshal(m)
	require.NoError(t, err)

	var wire struct {
		DateRange json.RawMessage `json:"dateRange"`
	}

	err = json.Unmarshal(encoded, &wire)
	require.NoError(t, err)
	require.Equal(t, `["0001-01-01T00:00:00Z","2030-02-01T10:00:00Z"]`, string(wire.DateRange))
}

func requireDates(t *testing.T, want []*maintenance.DateTime, got []*maintenance.DateTime) {
	t.Helper()

	require.Len(t, got, len(want))

	for i := range want {
		if want[i] == nil {
			require.Nil(t, got[i], "element %d", i)
			continue
		}

		require.NotNil(t, got[i], "element %d", i)
		require.True(t, want[i].Equal(*got[i]), "element %d: want %v, got %v", i, want[i], got[i])
	}
}
