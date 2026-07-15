package analytics

import (
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils"
	"time"
)

func HaversinePerformance(user model.Location, events []model.Location) *dto.HaversinePerf {
	start := time.Now()
	results := make(map[string]float64)
	for _, e := range events {
		results[e.Name] = utils.Haversine(user.Lat, user.Lon, e.Lat, e.Lon)
	}
	return &dto.HaversinePerf{
		Results: results,
		Time:    time.Since(start),
	}
}

func DijkstraPerformance(user model.Location, events []model.Location) *dto.DijkstraPerf {
	start := time.Now()

	minLat, minLon, maxLat, maxLon := utils.BoundingBox(user, events, 2.0)
	osmNodes, ways, err := utils.FetchRoadGraph(minLat, minLon, maxLat, maxLon)
	if err != nil {
		fmt.Println("Overpass ERROR:", err)
		osmNodes = make(map[int64]*model.OSMNode)
		ways = nil
	}

	graph, snapCoords := utils.BuildGraphFromOSM(user, events, osmNodes, ways)
	distMap, prevMap := utils.Dijkstra(*graph, user.Name, events[0].Name)

	return &dto.DijkstraPerf{
		DistMap:    distMap,
		PrevMap:    prevMap,
		OSMNodes:   osmNodes,
		SnapCoords: snapCoords, 
		UserLat:    user.Lat,   
		UserLon:    user.Lon,   
		Time:       time.Since(start),
	}
}
