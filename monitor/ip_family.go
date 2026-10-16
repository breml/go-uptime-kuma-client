package monitor

// IPFamily represents the IP protocol version a monitor uses to reach its
// target. The HTTP, keyword and json-query monitors carry it in
// IPFamilyDetails, the Globalping monitor in GlobalpingDetails.
type IPFamily string

// IP families. The empty value selects the IP family automatically.
const (
	IPFamilyAuto IPFamily = ""
	IPFamilyIPv4 IPFamily = "ipv4"
	IPFamilyIPv6 IPFamily = "ipv6"
)

// String returns the string representation of the IP family.
func (f IPFamily) String() string {
	return string(f)
}

// wireValue returns the IP family as it is sent to the server. IPFamilyAuto is
// sent as null, the value the server and its web interface use for the
// automatic selection.
func (f IPFamily) wireValue() any {
	if f == IPFamilyAuto {
		return nil
	}

	return f
}

// IPFamilyDetails contains the IP family setting embedded by the HTTP,
// HTTPKeyword and HTTPJSONQuery monitors, the only types the server applies it
// to besides Globalping, which carries its own GlobalpingDetails.IPFamily. It
// is an add-on to these monitors, not a monitor type, and therefore has no
// Type method.
type IPFamilyDetails struct {
	// IPFamily restricts the connection to IPv4 or IPv6. The server stores the
	// automatic selection as null, which unmarshals to IPFamilyAuto.
	IPFamily IPFamily `json:"ipFamily"`
}
