package main

// Episode -
type Episode struct {
	Name     string   `json:"name,omitempty"`
	Season   string   `json:"season,omitempty"`
	Episode  string   `json:"episode,omitempty"`
	Location string   `json:"location,omitempty"`
	Files    []string `json:"files,omitempty"`
}

// Show -
type Show struct {
	Name    string                `json:"name,omitempty"`
	Seasons map[string][]*Episode `json:"seasons,omitempty"`
	Scraped bool
	Multi   bool
}

// Shows -
type Shows map[string]*Show
