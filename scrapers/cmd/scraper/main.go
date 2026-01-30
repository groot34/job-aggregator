package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/groot34/job-aggregator/scraper/internal/models"
	"github.com/groot34/job-aggregator/scraper/internal/parsers"
	"github.com/groot34/job-aggregator/scraper/internal/publisher"
	"github.com/groot34/job-aggregator/scraper/internal/skills"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using defaults")
	}

	fmt.Println("🚀 Job Scraper Service Started")

	// Registry of parsers
	siteParsers := []parsers.Parser{
		&parsers.YCombinatorParser{},
		&parsers.FreshersworldParser{},
		&parsers.LinkedInParser{},
		&parsers.WellfoundParser{},
		// Add more parsers here
	}

	runScrapers(siteParsers)
}

func runScrapers(siteParsers []parsers.Parser) {
	var wg sync.WaitGroup
	results := make(chan []models.Job, len(siteParsers))

	for _, p := range siteParsers {
		wg.Add(1)
		go func(parser parsers.Parser) {
			defer wg.Done()
			fmt.Printf("🕷️  Starting scraper for: %s\n", parser.Name())

			// For YC and others we don't need a specific URL in Parse() args
			// But for interface consistency we might pass a dummy or specific start URL
			jobs, err := parser.Parse("")
			if err != nil {
				log.Printf("❌ Error scraping %s: %v\n", parser.Name(), err)
				return
			}
			results <- jobs
		}(p)
	}

	// Close channel when all done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allFilteredJobs []models.Job
	for jobs := range results {
		for _, j := range jobs {
			// Skill & Domain Filtering
			// Check if it's a software job
			fullText := j.Title + " " + j.Description + " " + strings.Join(j.Tags, " ")
			if !skills.IsSoftwareJob(fullText) {
				continue
			}

			// Add extracted skills to the job model
			j.Tags = skills.ExtractSkills(fullText)
			allFilteredJobs = append(allFilteredJobs, j)

			// Print for demo
			// fmt.Printf("   Found: %s\n", j.Title)
		}
	}

	if len(allFilteredJobs) > 0 {
		fmt.Printf("📦 Preparing to send %d valid jobs to backend...\n", len(allFilteredJobs))
		if err := publisher.PublishJobs(allFilteredJobs); err != nil {
			log.Printf("❌ Failed to publish jobs: %v\n", err)
		}
	}

	fmt.Printf("\n🏁 Scrape finished. Total valid jobs processed: %d\n", len(allFilteredJobs))

	// Keep alive if needed (e.g. cron mode), for now exit
	time.Sleep(2 * time.Second)
}
