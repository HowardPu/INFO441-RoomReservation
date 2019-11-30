package models

// Search Term that the users will provide to obtain results
type SearchTerm struct {
	Floor    int64  `json:"type,omitempty"`
	Name     string `json:"type,omitempty"`
	Type     string `json:"type,omitempty"`
	Capacity string `json:"type,omitempty"`
}
