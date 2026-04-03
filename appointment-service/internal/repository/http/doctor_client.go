package http

import (
	"fmt"
	"net/http"
	"time"
)

type DoctorClient struct {
	baseURL string
	client  *http.Client
}

func NewDoctorClient(baseURL string) *DoctorClient {
	return &DoctorClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *DoctorClient) GetDoctor(id string) (bool, error) {
	url := fmt.Sprintf("%s/doctors/%s", c.baseURL, id)

	resp, err := c.client.Get(url)
	if err != nil {
		return false, fmt.Errorf("doctor service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("doctor service error: %d", resp.StatusCode)
	}

	return true, nil
}
