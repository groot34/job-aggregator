package parsers

import (
	"context"
	"fmt"
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
		// Try multiple possible selectors
		const jobElements = document.querySelectorAll('[class*="JobsList_tableRow"], [class*="ycdc-card"], a[href*="/jobs/"]');
		
		// Get all job links
		const jobLinks = Array.from(document.querySelectorAll('a[href*="/companies/"][href*="/jobs/"]'));
		
		jobLinks.forEach(link => {
			try {
				const href = link.getAttribute('href');
				const title = link.textContent.trim();
				
				// Find parent container
				let parent = link.closest('div');
				if (!parent) return;
				
				// Try to find company link in same parent
				const companyLink = parent.querySelector('a[href*="/companies/"]:not([href*="/jobs/"])');
				const company = companyLink ? companyLink.textContent.trim() : '';
				
				// Extract all text from parent for location/salary
				const parentText = parent.textContent;
				
				// Look for location patterns
				let location = '';
				const locationMatch = parentText.match(/(Remote|[A-Z][a-z]+,\s*[A-Z]{2}|San Francisco|New York|London|Bangalore)/i);
				if (locationMatch) location = locationMatch[0];
				
				// Look for salary patterns
				let salary = '';
				const salaryMatch = parentText.match(/\$[\d]+K\s*-\s*\$[\d]+K/);
				if (salaryMatch) salary = salaryMatch[0];
				
				if (title && href && company) {
					jobs.push({
						title: title,
						company: company,
						url: href.startsWith('http') ? href : 'https://www.ycombinator.com' + href,
						location: location,
						salary: salary
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

	for _, data := range jobsData {
		title, _ := data["title"].(string)
		company, _ := data["company"].(string)
		url, _ := data["url"].(string)
		location, _ := data["location"].(string)
		salary, _ := data["salary"].(string)

		if title == "" || company == "" || url == "" {
			continue
		}

		// Extract batch and create tags
		var tags []string
		tags = append(tags, "startup", "yc")
		if strings.Contains(company, "(") && strings.Contains(company, ")") {
			batchStart := strings.Index(company, "(")
			batchEnd := strings.Index(company, ")")
			if batchStart < batchEnd {
				batch := strings.TrimSpace(company[batchStart+1 : batchEnd])
				tags = append(tags, "YC-"+batch)
			}
		}

		// Generate unique ID
		jobID := "yc-" + strings.ReplaceAll(url, "https://www.ycombinator.com/companies/", "")
		jobID = strings.ReplaceAll(jobID, "/", "-")

		// Check if remote
		isRemote := strings.Contains(strings.ToLower(location), "remote")

		job := models.Job{
			ID:          jobID,
			Title:       strings.TrimSpace(title),
			Company:     strings.TrimSpace(company),
			Location:    strings.TrimSpace(location),
			Description: "",
			URL:         url,
			Source:      "YCombinator",
			PostedAt:    time.Now(),
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
