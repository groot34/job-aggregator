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

	// Search URLs: General software engineer, remote internships, fresher software engineer, backend developer
	targetURLs := []string{
		"https://www.linkedin.com/jobs/search?keywords=software%20engineer&location=India&geoId=102713980&trk=public_jobs_jobs-search-bar_search-submit&position=1&pageNum=0",
		"https://www.linkedin.com/jobs/search?keywords=remote&location=India&geoId=102713980&f_TPR=r86400&f_WT=2&f_E=1,2&position=1&pageNum=0",
		"https://www.linkedin.com/jobs/search?keywords=remote%20internships&geoId=103644278&f_TPR=r86400&position=1&pageNum=0",
		"https://www.linkedin.com/jobs/search?keywords=fresher%20software%20engineer&geoId=114806696&f_TPR=r86400&position=1&pageNum=0",
		"https://www.linkedin.com/jobs/search?keywords=backend%20developer&geoId=114806696&f_TPR=r86400&position=1&pageNum=0",
		"https://www.linkedin.com/jobs/search?keywords=remote%20software%20engineer&location=India&geoId=102713980&f_TPR=r86400&f_WT=2&position=1&pageNum=0",
	}

	var jobs []models.Job
	seenJobIDs := make(map[string]bool)

	c := colly.NewCollector(
		// LinkedIn is sensitive to User-Agents. Use a standard one.
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	// LinkedIn public job cards usually define this structure
	c.OnHTML("ul.jobs-search__results-list li", func(e *colly.HTMLElement) {
		title := e.ChildText("h3.base-search-card__title")
		company := e.ChildText("h4.base-search-card__subtitle")
		location := e.ChildText("span.job-search-card__location")
		link := e.ChildAttr("a.base-card__full-link", "href")
		dateStr := e.ChildAttr("time", "datetime")
		dateText := e.ChildText("time")

		if title == "" || link == "" {
			return
		}

		jobID := "li-" + getIDFromURL(link)
		if seenJobIDs[jobID] {
			return
		}
		seenJobIDs[jobID] = true

		postedAt := time.Now()
		if dateStr != "" {
			// Format: 2023-10-25 or RFC3339
			if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				postedAt = t
			} else if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
				postedAt = t
			}
		} else if dateText != "" {
			// Fallback: parse e.g. "2 days ago", "1 hour ago" from the visible text
			postedAt = ParseRelativeDate(dateText)
		}

		lowerLocation := strings.ToLower(location)
		lowerTitle := strings.ToLower(title)
		isRemote := strings.Contains(lowerLocation, "remote") ||
			strings.Contains(lowerTitle, "remote") ||
			strings.Contains(lowerLocation, "hybrid") ||
			strings.Contains(lowerTitle, "hybrid")

		job := models.Job{
			ID:          jobID,
			Title:       strings.TrimSpace(title),
			Company:     strings.TrimSpace(company),
			Location:    strings.TrimSpace(location),
			URL:         link,
			Source:      "LinkedIn",
			PostedAt:    postedAt,
			ScrapedAt:   time.Now(),
			Description: "Click to apply on LinkedIn to view full description.",
			Remote:      isRemote,
			Tags:        []string{"linkedin"},
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

	for _, targetURL := range targetURLs {
		err := c.Visit(targetURL)
		if err != nil {
			fmt.Printf("❌ LinkedIn Scrape Error for %s: %v\n", targetURL, err)
		}
	}

	fmt.Printf("✅ Found %d jobs from LinkedIn\n", len(jobs))
	return jobs, nil
}
