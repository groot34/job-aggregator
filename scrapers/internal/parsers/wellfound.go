package parsers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/groot34/job-aggregator/scraper/internal/models"
)

type WellfoundParser struct{}

func (p *WellfoundParser) Name() string {
	return "Wellfound"
}

func (p *WellfoundParser) Parse(arg string) ([]models.Job, error) {
	fmt.Println("🔌 Fetching jobs from Wellfound...")

	// Wellfound is very dynamic (React).
	// This simple HTML scraper might fail if they verify browser integrity heavily or render via JS mostly.
	// But public landing pages sometimes have SSR content.
	targetURL := "https://wellfound.com/role/software-engineer"

	var jobs []models.Job

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", "https://google.com")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("div[data-test='JobListItem']", func(e *colly.HTMLElement) {
		title := e.ChildText("h2") // Often h2 or similar
		company := e.ChildText("div[data-test='StartupName']")
		link := e.ChildAttr("a", "href")
		location := e.ChildText("[data-test='Location']")
		if location == "" {
			location = e.ChildText(".job-list-item__location")
		}
		salary := e.ChildText("[data-test='Compensation']")

		if title == "" {
			return
		}

		if IsClosedListing(e.Text) {
			return
		}

		// Wellfound URLs are relative often
		if !strings.HasPrefix(link, "http") {
			link = "https://wellfound.com" + link
		}

		// Posted date extraction — try relative date strings in the card text
		postedAt := ParseRelativeDate(e.Text)

		lowerLoc := strings.ToLower(location)
		isRemote := strings.Contains(lowerLoc, "remote") ||
			strings.Contains(strings.ToLower(title), "remote") ||
			strings.Contains(lowerLoc, "hybrid")

		job := models.Job{
			ID:          "wf-" + getIDFromURL(link),
			Title:       strings.TrimSpace(title),
			Company:     strings.TrimSpace(company),
			Location:    strings.TrimSpace(location),
			URL:         link,
			Source:      "Wellfound",
			PostedAt:    postedAt,
			ScrapedAt:   time.Now(),
			Description: "View on Wellfound",
			Remote:      isRemote,
			Salary:      strings.TrimSpace(salary),
			Tags:        []string{"wellfound", "startup"},
		}

		jobs = append(jobs, job)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(targetURL)
	if err != nil {
		fmt.Printf("❌ Wellfound Scrape Error: %v\n", err)
		return nil, nil
	}

	fmt.Printf("✅ Found %d jobs from Wellfound\n", len(jobs))
	return jobs, nil
}
