package maintenance

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	floatingLayoutDate    = "2006-01-02"
	floatingLayoutMinutes = "2006-01-02T15:04"
	floatingLayoutSeconds = "2006-01-02T15:04:05"
	offsetLayoutMinutes   = "2006-01-02T15:04Z07:00"
	offsetLayoutSeconds   = time.RFC3339Nano
)

// dateTimeKind tells what a DateTime holds.
type dateTimeKind uint8

const (
	kindAbsent dateTimeKind = iota
	kindOffset
	kindFloating
)

// DateTime is a date and time of a maintenance window as Uptime Kuma stores it.
//
// Uptime Kuma checks a date only with JavaScript's new Date(...). A window
// created in its UI carries the value of the browser's datetime-local input,
// such as "2026-06-29T13:39", with no seconds and no UTC offset, while a window
// written by NewSingleMaintenance carries RFC 3339. A server backed by SQLite
// keeps a date as it was sent.
//
// A server backed by MariaDB behaves differently, as observed on Uptime Kuma
// 2.5.6 with MariaDB 10.11: the database rejects a date with an offset, so a
// single window with RFC 3339 dates cannot be created there, and a floating
// date is returned as sent until the server restarts and as
// "2006-01-02 15:04:05" afterwards. Floating therefore does not tell that a
// window was created in the UI.
//
// A floating value, see Floating, carries no UTC offset: its digits are the
// wall-clock time in the maintenance window's Timezone. Time returns those
// digits in UTC; In reads them as wall-clock time in a location. A floating
// value carries whole seconds only: MariaDB drops a fraction and the server's
// date library misreads one, so NewFloatingDateTime and decoding drop it.
// Floating values are the only form both database backends accept and the
// server reads the same way everywhere.
//
// A value with an offset is an instant. Uptime Kuma 2.5 computes the status of
// a window with dayjs.tz(date, timezone), which reads the UTC digits of such a
// value as wall-clock time in the window's timezone, so only floating values,
// or values of a window whose timezone is UTC, mean the same to Uptime Kuma as
// they do here.
//
// Decoding and encoding keep the kind, floating or with offset, not the bytes:
// "2026-06-29 13:39:00" is written back as "2026-06-29T13:39" and
// "2026-06-29T13:39:00.000Z" as "2026-06-29T13:39:00Z". A value with an offset
// is encoded exactly as time.Time encodes it.
//
// The zero value is the absent date, decoded from and encoded as JSON null or
// empty text. It is distinct from any present date, including the zero
// instant "0001-01-01T00:00:00Z", which is a present value with an offset. The
// absent value is never floating.
//
// Compare values with Equal, not ==: two values with an offset holding the
// same instant can differ in their location, such as "Z" and "+00:00".
type DateTime struct {
	t    time.Time
	kind dateTimeKind
}

// NewDateTime returns t as a present DateTime with an offset, encoded as
// time.Time encodes it, also for the zero time.Time. The monotonic clock
// reading of t is dropped.
func NewDateTime(t time.Time) DateTime {
	return DateTime{t: t.Round(0), kind: kindOffset}
}

// NewFloatingDateTime returns the wall-clock digits of t in t's own location,
// truncated to whole seconds, as a present floating DateTime: the maintenance
// window's Timezone gives them their meaning. To get the digits of an instant
// in the window's timezone, pass t.In(loc).
func NewFloatingDateTime(t time.Time) DateTime {
	return DateTime{
		t:    time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC),
		kind: kindFloating,
	}
}

// Time returns the instant of a value with an offset. For a floating value it
// returns the wall-clock digits held in UTC, which are the instant only when
// the window's timezone is UTC; use In for the instant. For the absent value
// it returns the zero time.Time.
func (d DateTime) Time() time.Time {
	return d.t
}

// In reads a floating value's digits as wall-clock time in loc and converts a
// value with an offset to loc; the latter is not how Uptime Kuma's status
// computation reads such a value, see DateTime. A nil loc is UTC. The absent
// value gives the zero time.Time.
func (d DateTime) In(loc *time.Location) time.Time {
	if loc == nil {
		loc = time.UTC
	}

	if d.kind == kindOffset {
		return d.t.In(loc)
	}

	if d.kind == kindFloating {
		return time.Date(d.t.Year(), d.t.Month(), d.t.Day(), d.t.Hour(), d.t.Minute(), d.t.Second(), 0, loc)
	}

	return time.Time{}
}

// IsZero reports whether d is the absent date. A present date at the zero
// instant, such as "0001-01-01T00:00:00Z", is not zero.
func (d DateTime) IsZero() bool {
	return d.kind == kindAbsent
}

// Floating reports whether d is a present date without UTC offset, see
// DateTime.
func (d DateTime) Floating() bool {
	return d.kind == kindFloating
}

// Equal reports whether d and other are of the same kind and hold the same
// instant or, when floating, the same digits. Two absent values are equal.
func (d DateTime) Equal(other DateTime) bool {
	return d.kind == other.kind && d.t.Equal(other.t)
}

// String returns the text MarshalText writes, the empty string for the absent
// value. For a value MarshalText rejects, such as a year outside 0 to 9999 or
// an offset of a day or more, it formats the value in the same layout anyway.
func (d DateTime) String() string {
	if d.kind == kindOffset {
		return d.t.Format(offsetLayoutSeconds)
	}

	if d.kind == kindFloating {
		return d.formatFloating()
	}

	return ""
}

// MarshalJSON encodes the absent value as null and any other value as the
// string MarshalText returns.
func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}

	text, err := d.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("marshal maintenance date to JSON: %w", err)
	}

	// The text holds digits, "-", ":", ".", "T", "Z" and "+" only, none of
	// which JSON escapes.
	data := make([]byte, 0, len(text)+2)
	data = append(data, '"')
	data = append(data, text...)

	return append(data, '"'), nil
}

// UnmarshalJSON decodes null and the empty string as the absent value and any
// other string as UnmarshalText does.
func (d *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = DateTime{}
		return nil
	}

	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("maintenance date %s: not a string: %w", data, err)
	}

	return d.UnmarshalText([]byte(s))
}

// MarshalText encodes the absent value as empty text, a floating value in the
// Uptime Kuma UI's format, "2006-01-02T15:04", or "2006-01-02T15:04:05" when
// the seconds are not zero, and a value with an offset exactly as
// time.Time.MarshalText does, including its errors. A floating year outside 0
// to 9999 is an error too.
func (d DateTime) MarshalText() ([]byte, error) {
	if d.kind == kindOffset {
		text, err := d.t.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("marshal maintenance date: %w", err)
		}

		return text, nil
	}

	if d.kind == kindFloating {
		if year := d.t.Year(); year < 0 || year > 9999 {
			return nil, fmt.Errorf("marshal maintenance date %s: year outside of range [0,9999]", d.String())
		}

		return []byte(d.formatFloating()), nil
	}

	return []byte{}, nil
}

// UnmarshalText decodes empty text as the absent value; RFC 3339 and
// "2006-01-02T15:04Z07:00", also with a space instead of the "T", as a value
// with an offset; and "2006-01-02T15:04", "2006-01-02T15:04:05",
// "2006-01-02 15:04", "2006-01-02 15:04:05" and "2006-01-02" as a floating
// value, whose fractional seconds, if any, are dropped. It always replaces all
// of d.
func (d *DateTime) UnmarshalText(text []byte) error {
	*d = DateTime{}

	s := string(text)
	if s == "" {
		return nil
	}

	layout, kind := layoutOf(s)

	t, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf(
			"maintenance date %q: want RFC 3339 or a date and time without offset such as %q: %w",
			s, floatingLayoutMinutes, err,
		)
	}

	if kind == kindFloating {
		*d = NewFloatingDateTime(t)
		return nil
	}

	*d = DateTime{t: t, kind: kind}

	return nil
}

// formatFloating formats a floating value without checking its year.
func (d DateTime) formatFloating() string {
	if d.t.Second() == 0 {
		return d.t.Format(floatingLayoutMinutes)
	}

	return d.t.Format(floatingLayoutSeconds)
}

// layoutOf picks the one layout s can be in from its shape, so that no string
// is parsed as both floating and with an offset, and a parse error tells what
// is wrong with s in that layout. Seconds in a layout accept a fraction.
func layoutOf(s string) (string, dateTimeKind) {
	if len(s) == len(floatingLayoutDate) {
		return floatingLayoutDate, kindFloating
	}

	timePart := ""
	if len(s) > len(floatingLayoutDate) {
		timePart = s[len(floatingLayoutDate)+1:]
	}

	withSeconds := len(timePart) > len("15:04") && timePart[len("15:04")] == ':'

	var layout string

	kind := kindFloating

	switch {
	case (strings.HasSuffix(timePart, "Z") || strings.ContainsAny(timePart, "+-")) && withSeconds:
		layout, kind = offsetLayoutSeconds, kindOffset

	case strings.HasSuffix(timePart, "Z") || strings.ContainsAny(timePart, "+-"):
		layout, kind = offsetLayoutMinutes, kindOffset

	case withSeconds:
		layout = floatingLayoutSeconds

	default:
		layout = floatingLayoutMinutes
	}

	if len(s) > len(floatingLayoutDate) && s[len(floatingLayoutDate)] == ' ' {
		layout = strings.Replace(layout, "T", " ", 1)
	}

	return layout, kind
}
