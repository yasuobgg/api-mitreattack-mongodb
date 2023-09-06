package models

type APT struct {
	ID              string `json:"ID"`
	Name            string `json:"Name"`
	AssociatedGroup string `json:"AssociatedGroups"`
	Description     string `json:"Description"`
	Timestamp       int64  `json:"Timestamp"`
}