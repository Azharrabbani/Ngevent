package analytics

import (
	"fmt"
	"math"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils"
)

func HaversineAccuracy(perf *dto.HaversinePerf, events []model.Location) (string, float64, dto.Haversine) {
	minDist := math.Inf(1)
	totalError := 0.0
	var nearest dto.Haversine

	for _, e := range events {
		dist := perf.Results[e.Name]

		if dist < minDist {
			minDist = dist
			nearest = dto.Haversine{
				Name:     e.Name,
				Distance: fmt.Sprintf("%.2f km", dist),
			}
		}
	}

	_ = totalError

	accuracy := 100.0
	return fmt.Sprintf("%.4f s", perf.Time.Seconds()), accuracy, nearest
}

func DijkstraAccuracy(perf *dto.DijkstraPerf, havResults map[string]float64, events []model.Location) (string, float64, dto.Dijkstra, []dto.PathPoint) {
	minDist := math.Inf(1)
	totalError := 0.0
	var nearest dto.Dijkstra
	var path []dto.PathPoint

	for _, e := range events {
		dij := perf.DistMap[e.Name]
		hav := havResults[e.Name]

		if dij < minDist {
			minDist = dij
			nearest = dto.Dijkstra{
				Name:     e.Name,
				Distance: fmt.Sprintf("%.2f km", dij),
			}
		}

		if dij != 0 {
			totalError += math.Abs(hav-dij) / dij * 100
		}
	}

	// Build path for the nearest event
	for _, e := range events {
		if e.Name == nearest.Name {
			path = utils.BuildRoadPathWithCoords(
				perf.PrevMap,
				"user",
				e.Name,
				perf.OSMNodes,
				0, 0, // user lat/lon not stored here — passed via caller
				e.Lat,
				e.Lon,
			)
			break
		}
	}

	avgError := totalError / float64(len(events))
	accuracy := 100 - avgError

	return fmt.Sprintf("%.4f s", perf.Time.Seconds()), accuracy, nearest, path
}
