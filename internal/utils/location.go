package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"ngevent/internal/dto"
	"time"
)

func ReverseGeocode(lat, lon string) (*dto.ReverseResponse, error) {

	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json",
		lat,
		lon,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "ngevent-app/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.ReverseResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func SearchLocation(query string) (*[]dto.SearchResponse, error) {
	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?q=%s&format=json&addressdetails=1&limit=10",
		neturl.QueryEscape(query),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "ngevent-app/1.0 (azharr078@gmail.com)")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nominatim returned status %d", resp.StatusCode)
	}

	var result []dto.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("no location found")
	}

	return &result, nil
}
