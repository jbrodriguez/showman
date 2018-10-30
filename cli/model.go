package main

type Episode struct {
	Name     string   `json:"name,omitempty"`
	Season   string   `json:"season,omitempty"`
	Episode  string   `json:"episode,omitempty"`
	Location string   `json:"location,omitempty"`
	Files    []string `json:"files,omitempty"`
}

type Show struct {
	Name    string                `json:"name,omitempty"`
	Seasons map[string][]*Episode `json:"seasons,omitempty"`
}

type Shows map[string]*Show
