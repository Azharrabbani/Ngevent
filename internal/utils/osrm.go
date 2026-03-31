package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OSRMResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Legs     []struct {
			Steps []struct {
				Name     string  `json:"name"`
				Distance float64 `json:"distance"`
				Duration float64 `json:"duration"`
			} `json:"steps"`
		} `json:"legs"`
	} `json:"routes"`
}

func GetDistanceOSRM(lat1, lon1, lat2, lon2 float64) (float64, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false",
		lon1, lat1, lon2, lat2,
	)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result OSRMResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Routes[0].Distance / 1000, nil
}

func GetRouteOSRM(lat1, lon1, lat2, lon2 float64) (*OSRMResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false&steps=true",
		lon1, lat1, lon2, lat2,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result OSRMResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func ExtractPath(route *OSRMResponse, eventName string) []string {
	var path []string

	path = append(path, "user")

	for _, leg := range route.Routes[0].Legs {
		for _, step := range leg.Steps {
			if step.Name != "" {
				path = append(path, step.Name)
			}
		}
	}

	path = append(path, eventName)

	return path
}
