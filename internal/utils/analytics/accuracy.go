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
	var nearest dto.Haversine

	for _, e := range events {
		if dist := perf.Results[e.Name]; dist < minDist {
			minDist = dist
			nearest = dto.Haversine{
				Name:     e.Name,
				Distance: fmt.Sprintf("%.2f km", dist),
			}
		}
	}

	return fmt.Sprintf("%.4f s", perf.Time.Seconds()), 100.0, nearest
}

func DijkstraAccuracy(perf *dto.DijkstraPerf, havResults map[string]float64, events []model.Location) (string, float64, dto.Dijkstra, []dto.PathPoint) {
	minDist := math.Inf(1)
	totalError := 0.0
	var nearest dto.Dijkstra
	var nearestPath []dto.PathPoint

	for _, e := range events {
		hav := havResults[e.Name]

		
		realDist, path := utils.BuildRoadPathWithCoords(
			perf.PrevMap,
			"user",
			e.Name,
			perf.OSMNodes,
			perf.SnapCoords,
			perf.UserLat,
			perf.UserLon,
			e.Lat,
			e.Lon,
		)

		if realDist < minDist {
			minDist = realDist
			nearest = dto.Dijkstra{
				Name:     e.Name,
				Distance: fmt.Sprintf("%.2f km", realDist),
			}
			nearestPath = path
		}

		if realDist != 0 {
			totalError += math.Abs(hav-realDist) / realDist * 100
		}
	}

	avgError := totalError / float64(len(events))
	return fmt.Sprintf("%.4f s", perf.Time.Seconds()), 100 - avgError, nearest, nearestPath
}
