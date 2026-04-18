package parsers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/groot34/job-aggregator/scraper/internal/models"
)

type IndeedParser struct{}

func (p *IndeedParser) Name() string {
	return "Indeed"
}

func (p *IndeedParser) Parse(arg string) ([]models.Job, error) {
	fmt.Println("🔌 Fetching jobs from Indeed...")

	targetURL := "https://in.indeed.com/jobs?q=software+engineer&l=India"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("window-size", "1920,1080"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var jobsData []map[string]interface{}
	js := `(() => {
		const cards = Array.from(document.querySelectorAll('a.jcs-JobTitle[data-jk]'));
		return cards.map(card => {
			const row = card.closest('tr');
			const title = card.innerText.trim();
			const url = card.href;
			const company = row?.querySelector('[data-testid="company-name"]')?.innerText.trim() || '';
			const location = row?.querySelector('[data-testid="text-location"]')?.innerText.trim() || '';
			const snippet = row?.querySelector('.job-snippet')?.innerText.trim() || '';
			const salary = row?.querySelector('.salary-snippet-container')?.innerText.trim() || '';
			const dateText = row?.querySelector('span.date')?.innerText.trim() || '';
			return {title, url, company, location, snippet, salary, dateText};
		});
	})()`

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible(`a.jcs-JobTitle[data-jk]`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(js, &jobsData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape Indeed: %v", err)
	}

	var jobs []models.Job
	for _, data := range jobsData {
		title, _ := data["title"].(string)
		url, _ := data["url"].(string)
		company, _ := data["company"].(string)
		location, _ := data["location"].(string)
		snippet, _ := data["snippet"].(string)
		salary, _ := data["salary"].(string)
		dateText, _ := data["dateText"].(string)

		if title == "" || url == "" {
			continue
		}

		job := models.Job{
			ID:          "in-" + getIDFromURL(url),
			Title:       strings.TrimSpace(title),
			Company:     strings.TrimSpace(company),
			Location:    strings.TrimSpace(location),
			Description: strings.TrimSpace(snippet),
			URL:         url,
			Source:      "Indeed",
			PostedAt:    parseRelativeDate(dateText),
			ScrapedAt:   time.Now(),
			Remote:      strings.Contains(strings.ToLower(location), "remote") || strings.Contains(strings.ToLower(location), "work from home"),
			Salary:      strings.TrimSpace(salary),
			Tags:        []string{"indeed"},
		}

		jobs = append(jobs, job)
	}

	fmt.Printf("✅ Found %d jobs from Indeed\n", len(jobs))
	return jobs, nil
}

func parseRelativeDate(text string) time.Time {
	text = strings.ToLower(strings.TrimSpace(text))
	now := time.Now()

	if strings.Contains(text, "today") || strings.Contains(text, "just posted") || strings.Contains(text, "hours ago") {
		return now
	}

	if strings.Contains(text, "30+") {
		return now.AddDate(0, 0, -30)
	}

	number := 0
	for _, part := range strings.Fields(text) {
		if n, err := fmt.Sscanf(part, "%d", &number); err == nil && n == 1 {
			break
		}
	}

	if strings.Contains(text, "day") && number > 0 {
		return now.AddDate(0, 0, -number)
	}

	return now
}
