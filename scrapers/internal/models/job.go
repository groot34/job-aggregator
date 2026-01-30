package models

import "time"

// Job represents a standardized job listing
type Job struct {
	ID          string    `json:"externalId"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Source      string    `json:"source"` // e.g., "LinkedIn", "Indeed"
	PostedAt    time.Time `json:"postedAt"`
	ScrapedAt   time.Time `json:"scrapedAt"`
	Remote      bool      `json:"remote"`
	Salary      string    `json:"salary,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}
