package domain

// Show -
type Show struct {
	Name    string               `json:"name,omitempty"`
	Seasons map[string][]Episode `json:"seasons,omitempty"`
	Scraped bool
	Multi   bool
}

// Shows -
type Shows map[string]Show
