package types

import "time"

type Issue struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Status    string     `json:"Status"`
	CreatedAt time.Time  `json:"created_at"`
	SubIssues []SubIssue `json:"sub_issues,omitempty"`
	Comments  []Comment  `json:"comments,omitempty"`
}
