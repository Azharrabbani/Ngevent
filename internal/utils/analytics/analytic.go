package analytics

import (
	"log"
	"math"
	"ngevent/internal/model"
	"ngevent/internal/utils"
)

func ComputeAnalytic(user model.Location, events []model.Location) {
    // Haversine
    havPerf := HaversinePerformance(user, events)
    havTime, _, havNearest := HaversineAccuracy(havPerf, events)

    // Dijkstra
    dijPerf := DijkstraPerformance(user, events)
    dijTime, _, dijNearest, _ := DijkstraAccuracy(dijPerf, havPerf.Results, events)

    for _, e := range events {
        havDist := havPerf.Results[e.Name]
        dijDist := dijPerf.DistMap[e.Name]

        osrmDist, err := utils.GetDistanceOSRM(user.Lat, user.Lon, e.Lat, e.Lon)
        if err != nil {
            log.Printf("[ROUTE] OSRM failed for %s: %v\n", e.Name, err)
            continue
        }

        havAccuracy := 0.0
        if osrmDist > 0 {
            havAccuracy = (1 - math.Abs(osrmDist-havDist)/osrmDist) * 100
        }

        dijAccuracy := 0.0
        if osrmDist > 0 {
            dijAccuracy = (1 - math.Abs(osrmDist-dijDist)/osrmDist) * 100
        }

        log.Printf("[ANALYTIC] ============================================\n")
        log.Printf("[ANALYTIC] Event         : %s\n", e.Name)
        log.Printf("[ANALYTIC] Ground Truth  : %.2f km (OSRM)\n", osrmDist)
        log.Printf("[ANALYTIC] --------------------------------------------\n")
        log.Printf("[ANALYTIC] Haversine     : %.2f km | time: %s | akurasi: %.2f%%\n",
            havDist, havTime, havAccuracy)
        log.Printf("[ANALYTIC] Dijkstra      : %.2f km | time: %s | akurasi: %.2f%%\n",
            dijDist, dijTime, dijAccuracy)
        log.Printf("[ANALYTIC] Detour factor : %.2fx (jalan %.1f%% lebih panjang dari garis lurus)\n",
            osrmDist/havDist, (osrmDist/havDist-1)*100)
        log.Printf("[ANALYTIC] ============================================\n")
    }

    _ = havNearest
    _ = dijNearest
}
