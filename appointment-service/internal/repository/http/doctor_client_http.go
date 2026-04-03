package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"appointment-service/internal/apperr"
)

type DoctorClientHTTP struct {
	BaseURL string
	client  *http.Client
}

const (
	maxDoctorAttempts = 3
	retryBaseDelay    = 150 * time.Millisecond
)

func NewDoctorClientHTTP(baseURL string) *DoctorClientHTTP {
	return &DoctorClientHTTP{
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *DoctorClientHTTP) DoctorExists(doctorID string) (bool, error) {
	url := fmt.Sprintf("%s/doctors/%s", c.BaseURL, doctorID)
	client := c.client
	if client == nil {
		client = http.DefaultClient
	}

	for attempt := 1; attempt <= maxDoctorAttempts; attempt++ {
		resp, err := client.Get(url)
		if err != nil {
			if attempt < maxDoctorAttempts {
				time.Sleep(retryBaseDelay * time.Duration(attempt))
				continue
			}
			log.Printf("doctor service unavailable after %d attempts: doctor_id=%s err=%v", attempt, doctorID, err)
			return false, fmt.Errorf("%w: %v", apperr.ErrDoctorServiceUnavailable, err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}

		if resp.StatusCode == http.StatusOK {
			return true, nil
		}

		if isRetryableStatus(resp.StatusCode) && attempt < maxDoctorAttempts {
			time.Sleep(retryBaseDelay * time.Duration(attempt))
			continue
		}

		if len(body) == 0 {
			log.Printf("doctor service error: doctor_id=%s status=%d", doctorID, resp.StatusCode)
			return false, fmt.Errorf("%w: status %d", apperr.ErrDoctorServiceError, resp.StatusCode)
		}
		log.Printf("doctor service error: doctor_id=%s status=%d body=%s", doctorID, resp.StatusCode, string(body))
		return false, fmt.Errorf("%w: %s", apperr.ErrDoctorServiceError, string(body))
	}

	log.Printf("doctor service unavailable: doctor_id=%s err=retry exhausted", doctorID)
	return false, fmt.Errorf("%w: retry exhausted", apperr.ErrDoctorServiceUnavailable)
}

func isRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}
