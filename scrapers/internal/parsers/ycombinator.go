package parsers

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/groot34/job-aggregator/scraper/internal/models"
)

type YCombinatorParser struct{}

func (p *YCombinatorParser) Name() string {
	return "YCombinator"
}

func (p *YCombinatorParser) Parse(arg string) ([]models.Job, error) {
	fmt.Println("🔌 Fetching jobs from Y Combinator (using headless browser)...")

	targetURL := "https://www.ycombinator.com/jobs/role/software-engineer"

	// Create chromedp context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	// JavaScript to extract job data
	jsCode := `
	(function() {
		const jobs = [];

		// Get all job links
		const jobLinks = Array.from(document.querySelectorAll('a[href*="/companies/"][href*="/jobs/"]'));

		jobLinks.forEach(link => {
			try {
				const href = link.getAttribute('href');
				const title = link.textContent.trim();

				// Walk up parents to find the largest reasonable container
				let parent = link.closest('div');
				if (!parent) return;
				// Try to get a wider container if possible
				for (let i = 0; i < 4 && parent && parent.parentElement; i++) {
					const p = parent.parentElement;
					if (p && p.tagName === 'DIV') parent = p;
				}

				const parentText = parent.textContent || '';

				// Try to find company link in same parent
				const companyLink = parent.querySelector('a[href*="/companies/"]:not([href*="/jobs/"])');
				const company = companyLink ? companyLink.textContent.trim() : '';

				// Look for location patterns
				let location = '';
				const locationMatch = parentText.match(/(Remote|[A-Z][a-z]+,\s*[A-Z]{2}|San Francisco|New York|London|Bangalore|Bengaluru|Hybrid|On[\s-]?site)/i);
				if (locationMatch) location = locationMatch[0];

				// Look for salary patterns
				let salary = '';
				const salaryMatch = parentText.match(/\$[\d]+K\s*-\s*\$[\d]+K/);
				if (salaryMatch) salary = salaryMatch[0];

				// Look for a posted-date string: "(about X ago)", "X hours ago", "X days ago"
				let dateText = '';
				const dateMatch = parentText.match(/\(?(\b(?:about\s+)?\d+\s+(?:hour|hours|hr|hrs|day|days|week|weeks|month|months|min|minute|minutes)\s+ago)\)?/i);
				if (dateMatch) dateText = dateMatch[1];

				// Also try explicit time elements
				if (!dateText) {
					const timeEl = parent.querySelector('time');
					if (timeEl) {
						if (timeEl.getAttribute('datetime')) {
							dateText = timeEl.getAttribute('datetime');
						} else if (timeEl.textContent) {
							dateText = timeEl.textContent.trim();
						}
					}
				}

				if (title && href && company) {
					jobs.push({
						title: title,
						company: company,
						url: href.startsWith('http') ? href : 'https://www.ycombinator.com' + href,
						location: location,
						salary: salary,
						dateText: dateText,
						fullText: parentText.slice(0, 500)
					});
				}
			} catch(e) {
				console.error('Error parsing job:', e);
			}
		});

		return jobs;
	})();
	`

	var jobsData []map[string]interface{}

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible(`a[href*="/companies/"]`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(jsCode, &jobsData),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scrape YC: %v", err)
	}

	// Convert to Job models
	var jobs []models.Job

	ycDateFallback := regexp.MustCompile(`(?i)\(?((?:about\s+)?\d+\s+(?:hour|hours|hr|hrs|day|days|week|weeks|month|months|min|minute|minutes)\s+ago)\)?`)

	for _, data := range jobsData {
		title, _ := data["title"].(string)
		company, _ := data["company"].(string)
		url, _ := data["url"].(string)
		location, _ := data["location"].(string)
		salary, _ := data["salary"].(string)
		dateText, _ := data["dateText"].(string)
		fullText, _ := data["fullText"].(string)

		if title == "" || company == "" || url == "" {
			continue
		}

		if IsClosedListing(title+" "+company+" "+fullText) {
			continue
		}

		// If JS didn't capture a dateText, try a last-ditch parse from fullText
		// (catches the common YC pattern "CompanyName (S21)•blurb(about 19 hours ago)")
		if dateText == "" && fullText != "" {
			if m := ycDateFallback.FindStringSubmatch(fullText); m != nil {
				dateText = m[1]
			}
		}

		// Parse the posted date. Fall back to scrape time if the format is unknown.
		var postedAt time.Time
		if dateText != "" {
			// Try ISO date first (time[datetime])
			if t, err := time.Parse(time.RFC3339, dateText); err == nil {
				postedAt = t
			} else if t, err := time.Parse("2006-01-02", dateText); err == nil {
				postedAt = t
			} else {
				postedAt = ParseRelativeDate(dateText)
			}
		} else {
			postedAt = time.Now()
		}

		// Extract batch and create tags
		var tags []string
		tags = append(tags, "startup", "yc")
		if strings.Contains(company, "(") && strings.Contains(company, ")") {
			batchStart := strings.Index(company, "(")
			batchEnd := strings.Index(company, ")")
			if batchStart < batchEnd {
				batch := strings.TrimSpace(company[batchStart+1 : batchEnd])
				// Only add batch codes that look like "W24", "S21", etc.
				if len(batch) <= 4 {
					tags = append(tags, "YC-"+batch)
				}
			}
		}

		// Generate unique ID
		jobID := "yc-" + strings.ReplaceAll(url, "https://www.ycombinator.com/companies/", "")
		jobID = strings.ReplaceAll(jobID, "/", "-")

		// Check if remote
		lowerLocation := strings.ToLower(location)
		isRemote := strings.Contains(lowerLocation, "remote")

		// Build description from the page text snippet (truncate long)
		description := ""
		if fullText != "" {
			// strip leading whitespace
			description = strings.TrimSpace(fullText)
			if len(description) > 800 {
				description = description[:800]
			}
		}

		job := models.Job{
			ID:          jobID,
			Title:       strings.TrimSpace(title),
			Company:     strings.TrimSpace(company),
			Location:    strings.TrimSpace(location),
			Description: description,
			URL:         url,
			Source:      "YCombinator",
			PostedAt:    postedAt,
			ScrapedAt:   time.Now(),
			Remote:      isRemote,
			Salary:      salary,
			Tags:        tags,
		}

		jobs = append(jobs, job)
	}

	fmt.Printf("✅ Found %d jobs from Y Combinator\n", len(jobs))
	return jobs, nil
}
