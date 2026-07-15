package analytics

import (
	"log"
	"time"

	"ngevent/internal/model"
	"ngevent/internal/utils"
)


const boundingBoxPaddingKm = 5.0


func ComputeAnalytic(user model.Location, events []model.Location) {
	if len(events) == 0 {
		return
	}

	minLat, minLon, maxLat, maxLon := utils.BoundingBox(user, events, boundingBoxPaddingKm)

	fetchStart := time.Now()
	osmNodes, ways, err := utils.FetchRoadGraph(minLat, minLon, maxLat, maxLon)
	fetchElapsed := time.Since(fetchStart)
	if err != nil {
		log.Printf("[ANALYTIC] Gagal fetch data OSM: %v — analytic dibatalkan\n", err)
		return
	}

	graph, snapCoords := utils.BuildGraphFromOSM(user, events, osmNodes, ways)

	for _, e := range events {
		algoStart := time.Now()
		_, prevMap := utils.Dijkstra(*graph, user.Name, e.Name)
		dist, _ := utils.BuildRoadPathWithCoords(
			prevMap, user.Name, e.Name,
			osmNodes, snapCoords,
			user.Lat, user.Lon, e.Lat, e.Lon,
		)
		algoElapsed := time.Since(algoStart)

		log.Printf("[ANALYTIC] ============================================\n")
		log.Printf("[ANALYTIC] Event                    : %s\n", e.Name)
		log.Printf("[ANALYTIC] Jarak rute (Dijkstra)     : %.2f km\n", dist)
		log.Printf("[ANALYTIC] Waktu fetch OSM            : %.3fs (I/O jaringan, sekali untuk semua event)\n", fetchElapsed.Seconds())
		log.Printf("[ANALYTIC] Waktu kalkulasi Dijkstra   : %.4fs (murni algoritma: shortest path + rekonstruksi rute)\n", algoElapsed.Seconds())
		log.Printf("[ANALYTIC] ============================================\n")
	}
}
