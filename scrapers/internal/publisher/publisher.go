package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/groot34/job-aggregator/scraper/internal/models"
)

// PublishJobs sends a batch of jobs to the backend API
func PublishJobs(jobs []models.Job) error {
	apiUrl := os.Getenv("BACKEND_API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:5000/api/jobs/batch"
	}

	// Transform to JSON
	payload, err := json.Marshal(jobs)
	if err != nil {
		return fmt.Errorf("failed to marshal jobs: %v", err)
	}

	// Send POST request
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send request to backend: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("backend returned status: %d", resp.StatusCode)
	}

	fmt.Printf("📤 Sent %d jobs to backend successfully\n", len(jobs))
	return nil
}
