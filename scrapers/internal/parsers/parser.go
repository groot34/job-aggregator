package parsers

import (
	"github.com/groot34/job-aggregator/scraper/internal/models"
)

// Parser is the interface that all site-specific scrapers must implement
type Parser interface {
	Parse(url string) ([]models.Job, error)
	Name() string
}
