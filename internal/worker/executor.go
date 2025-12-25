package worker

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/theweird-kid/blaze/internal/models"
)

func executeHTTP(ctx context.Context, job *models.Job, run *models.JobRun) error {
	timeout := time.Duration(job.HTTP.TimeoutSec) * time.Second
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequestWithContext(
		ctx,
		job.HTTP.Method,
		job.HTTP.URL,
		bytes.NewReader(job.HTTP.Body),
	)
	if err != nil {
		return err
	}

	// Idempotency Key
	req.Header.Set("Idempotency-Key", run.IdempotencyKey)

	for k, v := range job.HTTP.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("server error %d", resp.StatusCode)
	}

	return nil
}
