package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ngevent/internal/dto"
)

func ReverseGeocode(lat, lon string) (*dto.ReverseResponse, error) {

	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json",
		lat,
		lon,
	)

	client := &http.Client{}

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
