package vault

import "time"

type Document struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Entries   []Entry   `json:"entries"`
}

type Entry struct {
	Alias             string    `json:"alias"`
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	URI               string    `json:"uri,omitempty"`
	Note              string    `json:"note,omitempty"`
	Tags              []string  `json:"tags,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PasswordUpdatedAt time.Time `json:"password_updated_at"`
}

func NewDocument(now time.Time) *Document {
	return &Document{
		Version:   1,
		CreatedAt: now.UTC(),
		UpdatedAt: now.UTC(),
		Entries:   []Entry{},
	}
}
