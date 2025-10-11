package tag

import "fmt"

type Tag struct {
	ID    int64  `json:"id,omitzero"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (t Tag) String() string {
	return fmt.Sprintf("Tag{ID: %d, Name: %s, Color: %s}", t.ID, t.Name, t.Color)
}

func (t Tag) GetID() int64 {
	return t.ID
}
