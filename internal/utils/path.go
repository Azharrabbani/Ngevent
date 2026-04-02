package utils

import (
	"fmt"
	"math"
	"ngevent/internal/model"
)

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Earth radis(km)
	const R = 6371

	lat1 = toRadians(lat1)
	lon1 = toRadians(lon1)
	lat2 = toRadians(lat2)
	lon2 = toRadians(lon2)

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Pow(math.Sin(dLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

func BuildGraph(User model.Location, events []model.Location) *model.Graph {
	graph := make(model.Graph)

	// user node
	graph[User.Name] = []model.Edge{}

	for _, event := range events {
		graph[event.Name] = []model.Edge{}
	}

	// edges user -> event
	for _, event := range events {
		dist, err := GetDistanceOSRM(User.Lat, User.Lon, event.Lat, event.Lon)
		if err != nil {
			fmt.Println("OSRM ERROR:", err)
			continue
		}

		graph[User.Name] = append(graph[User.Name], model.Edge{
			To:     event.Name,
			Weight: dist,
		})
	}

	return &graph
}
