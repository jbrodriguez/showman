package domain

// Episode -
type Episode struct {
	Name     string `json:"name,omitempty"`
	Series   string `json:"series,omitempty"`
	Season   string `json:"season,omitempty"`
	Episode  string `json:"episode,omitempty"`
	Location string `json:"location,omitempty"`
	// Files    []string `json:"files,omitempty"`
}
