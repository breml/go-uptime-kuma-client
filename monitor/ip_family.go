package monitor

// IPFamily represents the IP protocol version a monitor uses to reach its
// target.
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

// IPFamilyDetails contains the IP family setting of the HTTP, keyword and
// json-query monitors. The server applies it to these monitor types only.
type IPFamilyDetails struct {
	// IPFamily restricts the connection to IPv4 or IPv6. The server stores the
	// automatic selection as null, which unmarshals to IPFamilyAuto.
	IPFamily IPFamily `json:"ipFamily"`
}

// ipFamilyValue returns the IP family as it is sent to the server.
// IPFamilyAuto is sent as null, the value the server and its web interface
// use for the automatic selection.
func (d IPFamilyDetails) ipFamilyValue() any {
	if d.IPFamily == IPFamilyAuto {
		return nil
	}

	return d.IPFamily
}
