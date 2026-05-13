package analytics

import (
	"fmt"
	"log"
	"math"
	"ngevent/internal/dto"
	"ngevent/internal/model"
)

func ComputeAnalytic(user model.Location, events []model.Location) {
	// Haversine
	havPerf := HaversinePerformance(user, events)
	havTime, _, havNearest := HaversineAccuracy(havPerf, events)

	//Dijkstra
	dijPerf := DijkstraPerformance(user, events)
	dijTime, dijAccuracy, dijNearest, _ := DijkstraAccuracy(dijPerf, havPerf.Results, events)

	// Haversine accuracy vs Dijkstra as reference
	totalErrorHav := 0.0
	for _, e := range events {
		hav := havPerf.Results[e.Name]
		dij := dijPerf.DistMap[e.Name]
		if hav != 0 {
			totalErrorHav += math.Abs(hav-dij) / hav * 100
		}
	}
	havAccuracy := 100 - (totalErrorHav / float64(len(events)))

	dijkstra := dto.Dijkstra{
		Name:     dijNearest.Name,
		Distance: dijNearest.Distance,
		Time:     dijTime,
		Accuracy: fmt.Sprintf("%.2f%%", dijAccuracy),
	}

	haversine := dto.Haversine{
		Name:     havNearest.Name,
		Distance: havNearest.Distance,
		Time:     havTime,
		Accuracy: fmt.Sprintf("%.2f%%", havAccuracy),
	}

	log.Printf("[ROUTE] Event     : %s\n", events[0].Name)
	log.Printf("[ROUTE] Haversine : name=%s dist=%s time=%s accuracy=%s\n",
		haversine.Name, haversine.Distance, haversine.Time, haversine.Accuracy)
	log.Printf("[ROUTE] Dijkstra  : name=%s dist=%s time=%s accuracy=%s\n",
		dijkstra.Name, dijkstra.Distance, dijkstra.Time, dijkstra.Accuracy)
}
