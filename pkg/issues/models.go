package issues

import "time"

type User struct {
	Login string `json:"login"`
	URL   string `json:"url"`
}

type Comment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	Author    User      `json:"author"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Issue struct {
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	State     string     `json:"state"`
	Author    User       `json:"author"`
	Labels    []string   `json:"labels"`
	Assignees []User     `json:"assignees"`
	Milestone string     `json:"milestone"`
	URL       string     `json:"url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	Comments  []Comment  `json:"comments"`
}

type IndexEntry struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}
