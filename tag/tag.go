package tag

import "fmt"

// Tag represents a tag that can be applied to monitors.
type Tag struct {
	ID    int64  `json:"id,omitzero"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (t Tag) String() string {
	return fmt.Sprintf("Tag{ID: %d, Name: %s, Color: %s}", t.ID, t.Name, t.Color)
}

// GetID returns the tag's unique identifier.
func (t Tag) GetID() int64 {
	return t.ID
}

// MonitorTag represents the association between a monitor and a tag.
type MonitorTag struct {
	ID        int64  `json:"id"`         // Association ID
	TagID     int64  `json:"tag_id"`     // Tag ID
	MonitorID int64  `json:"monitor_id"` // Monitor ID
	Value     string `json:"value"`      // Optional tag value for this monitor
	Name      string `json:"name"`       // Tag name (from joined tag table)
	Color     string `json:"color"`      // Tag color (from joined tag table, not customizable per association)
}

func (mt MonitorTag) String() string {
	return fmt.Sprintf("MonitorTag{ID: %d, TagID: %d, MonitorID: %d, Value: %s, Name: %s, Color: %s}",
		mt.ID, mt.TagID, mt.MonitorID, mt.Value, mt.Name, mt.Color)
}

// TagWithMonitors extends Tag with monitor associations.
//
//revive:disable:exported
type TagWithMonitors struct {
	Tag

	Monitors []int64 `json:"monitors"` // List of associated monitor IDs
}

//revive:enable:exported

// MonitorTags represents all tags for a specific monitor.
type MonitorTags struct {
	MonitorID int64        `json:"monitor_id"`
	Tags      []MonitorTag `json:"tags"`
}
