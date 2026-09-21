package parsers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/groot34/job-aggregator/scraper/internal/models"
)

type FreshersworldParser struct{}

func (p *FreshersworldParser) Name() string {
	return "Freshersworld"
}

func (p *FreshersworldParser) Parse(arg string) ([]models.Job, error) {
	fmt.Println("🔌 Fetching jobs from Freshersworld...")

	// Target URL for Freshers
	targetURL := "https://www.freshersworld.com/jobs"

	var jobs []models.Job

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Select individual job cards
	// Note: Selectors might change, this is a best-guess based on common structure.
	// In a real scenario, we'd inspect the DOM.
	// Assuming class "col-md-12 col-lg-12 col-xs-12 padding-none job-container" or similar
	c.OnHTML(".job-container", func(e *colly.HTMLElement) {
		title := e.ChildText(".latest-jobs-title")
		company := e.ChildText(".latest-jobs-company") // sometimes just text in a span
		location := e.ChildText(".job-location")
		desc := e.ChildText(".job-desc")
		relURL := e.ChildAttr("a[href]", "href")
		// Try multiple common selectors for posted date
		dateText := e.ChildText(".job-posted-date")
		if dateText == "" {
			dateText = e.ChildText(".posted-date")
		}
		if dateText == "" {
			dateText = e.ChildText("span.date")
		}
		if dateText == "" {
			// last-ditch: scan the whole card text for a relative date string
			dateText = e.Text
		}

		// Skip if essential data missing
		if title == "" || relURL == "" {
			return
		}

		if IsClosedListing(e.Text) {
			return
		}

		// Freshersworld often puts "Company Name" in a specific structure, fallback if dry
		if company == "" {
			company = "Unknown"
		}

		// Freshersworld date parsing is spotty, best-effort parse relative otherwise use scrape time
		postedAt := ParseRelativeDate(dateText)

		lowerLoc := strings.ToLower(location)
		isRemote := strings.Contains(lowerLoc, "remote") ||
			strings.Contains(lowerLoc, "work from home") ||
			strings.Contains(lowerLoc, "hybrid") ||
			strings.Contains(lowerLoc, "wfh")

		// Make relative links absolute
		absURL := relURL
		if !strings.HasPrefix(absURL, "http") {
			absURL = "https://www.freshersworld.com" + absURL
		}

		job := models.Job{
			ID:          "fw-" + getIDFromURL(relURL),
			Title:       strings.TrimSpace(title),
			Company:     strings.TrimSpace(company),
			Location:    strings.TrimSpace(location),
			Description: strings.TrimSpace(desc),
			URL:         absURL,
			Source:      "Freshersworld",
			PostedAt:    postedAt,
			ScrapedAt:   time.Now(),
			Remote:      isRemote,
			Tags:        []string{"freshersworld", "fresher", "india"},
		}

		jobs = append(jobs, job)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Visit the target
	err := c.Visit(targetURL)
	if err != nil {
		return nil, err
	}

	fmt.Printf("✅ Found %d jobs from Freshersworld\n", len(jobs))
	return jobs, nil
}
