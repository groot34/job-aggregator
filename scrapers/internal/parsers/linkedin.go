package parsers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/groot34/job-aggregator/scraper/internal/models"
)

type LinkedInParser struct{}

func (p *LinkedInParser) Name() string {
	return "LinkedIn"
}

func (p *LinkedInParser) Parse(arg string) ([]models.Job, error) {
	fmt.Println("🔌 Fetching jobs from LinkedIn...")

	// URL: Public job search for Software Engineers in India
	// We can parameterize this later
	targetURL := "https://www.linkedin.com/jobs/search?keywords=software%20engineer&location=India&geoId=102713980&trk=public_jobs_jobs-search-bar_search-submit&position=1&pageNum=0"

	var jobs []models.Job

	c := colly.NewCollector(
		// LinkedIn is sensitive to User-Agents. Use a standard one.
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// LinkedIn public job cards usually define this stricture
	c.OnHTML("ul.jobs-search__results-list li", func(e *colly.HTMLElement) {
		title := e.ChildText("h3.base-search-card__title")
		company := e.ChildText("h4.base-search-card__subtitle")
		location := e.ChildText("span.job-search-card__location")
		link := e.ChildAttr("a.base-card__full-link", "href")
		dateStr := e.ChildAttr("time", "datetime")

		if title == "" || link == "" {
			return
		}

		postedAt := time.Now()
		if dateStr != "" {
			// Format: 2023-10-25
			t, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				postedAt = t
			}
		}

		job := models.Job{
			ID:        "li-" + getIDFromURL(link),
			Title:     strings.TrimSpace(title),
			Company:   strings.TrimSpace(company),
			Location:  strings.TrimSpace(location),
			URL:       link,
			Source:    "LinkedIn",
			PostedAt:  postedAt,
			ScrapedAt: time.Now(),
			// Description is not on the card, would need to visit details page
			// But visiting details page triggers auth wall often.
			Description: "Click to apply on LinkedIn to view full description.",
			Remote:      strings.Contains(strings.ToLower(location), "remote"),
		}

		// Filter out obfuscated data
		if strings.Contains(job.Title, "**") || strings.Contains(job.Company, "**") {
			return
		}

		jobs = append(jobs, job)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(targetURL)
	if err != nil {
		fmt.Printf("❌ LinkedIn Scrape Error: %v\n", err)
		return nil, nil // Return empty, don't crash the whole run
	}

	fmt.Printf("✅ Found %d jobs from LinkedIn\n", len(jobs))
	return jobs, nil
}
