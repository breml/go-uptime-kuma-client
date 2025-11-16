package statuspage

type Incident struct {
	ID      int64  `json:"id,omitempty"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Style   string `json:"style"`
	Pin     bool   `json:"pin"`
}
