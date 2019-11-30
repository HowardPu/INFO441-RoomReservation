package models

// User represents the user information
type User struct {
	ID   int64  `json:"userID"`
	Name string `json:"userName"`
	Type string `json:"userType"`
}
