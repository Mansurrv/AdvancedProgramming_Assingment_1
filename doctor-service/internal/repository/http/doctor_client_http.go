package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type DoctorClientHTTP struct {
	BaseURL string
}

func NewDoctorClientHTTP(baseURL string) *DoctorClientHTTP {
	return &DoctorClientHTTP{BaseURL: baseURL}
}

func (c *DoctorClientHTTP) DoctorExists(doctorID string) (bool, error) {
	url := fmt.Sprintf("%s/doctors/%s", c.BaseURL, doctorID)
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("cannot reach doctor service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected response from doctor service: %s", string(body))
	}

	return true, nil
}
