package service

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	testdata "ngevent/internal"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils"
)

// osrmRequestDelay diberikan antar-call ke OSRM demo server publik untuk
// menghindari rate-limit/429.
const osrmRequestDelay = 3 * time.Second


const boundingBoxPaddingKm = 5.0


func RunDijkstraAccuracyTest() (*dto.AccuracyReport, error) {
	users, events := testdata.GenerateJabodetabekDataset()

	report := &dto.AccuracyReport{
		GeneratedAt: time.Now().Format(time.RFC3339),
		TotalUsers:  len(users),
		TotalEvents: len(events),
	}

	var overallErrors []float64

	for _, user := range users {
		log.Printf("========================================================")
		log.Printf("[Accuracy Test] Titik awal: %s (%.6f, %.6f)", user.Name, user.Lat, user.Lon)
		log.Printf("========================================================")

		result, err := runAccuracyTestForUser(user, events)
		if err != nil {
			log.Printf("[WARN] %v — titik awal ini dilewati", err)
			continue
		}

		log.Printf("--------------------------------------------------------")
		log.Printf("[MAPE] Titik awal %s: %s%% (dari %d pasangan valid, total waktu %ss)",
			user.Name, result.MAPEPercent, len(result.Cases), result.TotalExecTime)
		log.Printf("--------------------------------------------------------")

		report.Results = append(report.Results, *result)
		overallErrors = append(overallErrors, extractPercentErrors(result)...)
	}

	report.OverallMAPEPercent = average(overallErrors)

	log.Printf("========================================================")
	log.Printf("[Accuracy Test] SELESAI. Overall MAPE pooled (%d pasangan valid): %.2f%%",
		len(overallErrors), report.OverallMAPEPercent)
	log.Printf("========================================================")

	return report, nil
}


func runAccuracyTestForUser(user model.Location, events []model.Location) (*dto.UserMAPEResult, error) {
	minLat, minLon, maxLat, maxLon := utils.BoundingBox(user, events, boundingBoxPaddingKm)

	time.Sleep(osrmRequestDelay)
	osmNodes, ways, err := utils.FetchRoadGraph(minLat, minLon, maxLat, maxLon)
	if err != nil {
		return nil, fmt.Errorf("gagal fetch data OSM untuk %s: %w", user.Name, err)
	}
	log.Printf("[OSM] nodes=%d ways=%d", len(osmNodes), len(ways))

	graph, snapCoords := utils.BuildGraphFromOSM(user, events, osmNodes, ways)

	result := &dto.UserMAPEResult{UserName: user.Name}
	var userErrors []float64
	userStart := time.Now()

	for _, event := range events {
		algoStart := time.Now()
		_, prevMap := utils.Dijkstra(*graph, user.Name, event.Name)
		dijkstraDist, _ := utils.BuildRoadPathWithCoords(
			prevMap, user.Name, event.Name,
			osmNodes, snapCoords,
			user.Lat, user.Lon, event.Lat, event.Lon,
		)
		algoElapsed := time.Since(algoStart)
		execMs := float64(algoElapsed.Seconds())

		osrmDist, err := utils.GetRouteOSRM(user.Lat, user.Lon, event.Lat, event.Lon)	
		if err != nil {
			log.Printf("[WARN] OSRM gagal untuk %s -> %s: %v — pasangan ini dilewati", user.Name, event.Name, err)
			continue
		}

		dist := osrmDist.Routes[0].Distance / 1000
		var pctErr float64
		if dist != 0 {
			pctErr = math.Abs(dist-dijkstraDist) / dist * 100
		}

		caseResult := dto.TestCaseResult{
			UserName:  user.Name,
			EventName: event.Name,

			DijkstraDistanceKM: dijkstraDist,
			OSRMDistanceKM:     dist,
			PercentageErrorPct: pctErr,
			DijkstraExecTimeMs: execMs,

			DijkstraDistance: fmt.Sprintf("%.2fkm", dijkstraDist),
			OSRMDistance:     fmt.Sprintf("%.2fkm", dist),
			PercentageError:  fmt.Sprintf("%.2f%%", pctErr),
			DijkstraExecTime: fmt.Sprintf("%.2fs", execMs),
		}

		log.Printf("  [%s -> %s] Dijkstra=%s | OSRM=%s | Error=%s | Waktu Algoritma=%s",
			user.Name, event.Name, caseResult.DijkstraDistance, caseResult.OSRMDistance,
			caseResult.PercentageError, caseResult.DijkstraExecTime)

		result.Cases = append(result.Cases, caseResult)
		userErrors = append(userErrors, pctErr)
	}

	totalExecTimeMs := float64(time.Since(userStart).Seconds())
	mapeVal := average(userErrors)

	result.MAPEPercentValue = mapeVal
	result.TotalExecTimeMsValue = totalExecTimeMs
	result.MAPEPercent = fmt.Sprintf("%.2f%%", mapeVal)
	result.TotalExecTime = fmt.Sprintf("%.2fs", totalExecTimeMs)

	return result, nil
}

func RunSingleUserAccuracyTest(index int) (*dto.UserMAPEResult, error) {
	users, events := testdata.GenerateJabodetabekDataset()

	if index < 0 || index >= len(users) {
		return nil, fmt.Errorf("index tidak valid: %d (harus 0-%d, total %d titik awal)", index, len(users)-1, len(users))
	}

	user := users[index]

	log.Printf("========================================================")
	log.Printf("[Single Accuracy Test] index=%d titik awal: %s (%.6f, %.6f)", index, user.Name, user.Lat, user.Lon)
	log.Printf("========================================================")

	result, err := runAccuracyTestForUser(user, events)
	if err != nil {
		return nil, err
	}

	log.Printf("--------------------------------------------------------")
	log.Printf("[MAPE] Titik awal %s (index=%d): %s%% (dari %d pasangan valid, total waktu %ss)",
		user.Name, index, result.MAPEPercent, len(result.Cases), result.TotalExecTime)
	log.Printf("--------------------------------------------------------")

	return result, nil
}

// CalculateAverageMAPE menghitung rata-rata dari sekumpulan nilai mape_percent
// (string, tanpa "%") yang dikirim manual oleh tester — biasanya hasil
// copy-paste dari field "mape_percent" pada masing-masing response
// RunSingleUserAccuracyTest setelah semua titik awal (10) selesai diuji.
func CalculateAverageMAPE(rawValues []string) (*dto.MAPESummaryResponse, error) {
	if len(rawValues) == 0 {
		return nil, fmt.Errorf("mape_values tidak boleh kosong, contoh body: {\"mape_values\": [\"4.12\", \"3.87\"]}")
	}

	parsed := make([]float64, 0, len(rawValues))
	for i, raw := range rawValues {
		cleaned := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(raw), "%"))
		val, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return nil, fmt.Errorf("mape_values[%d]=%q bukan angka yang valid: %w", i, raw, err)
		}
		parsed = append(parsed, val)
	}

	avg := average(parsed)

	log.Printf("========================================================")
	log.Printf("[MAPE Summary] Menghitung rata-rata dari %d nilai MAPE yang dikirim manual:", len(parsed))
	for i, v := range parsed {
		log.Printf("  data-uji-%d: %.2f%%", i+1, v)
	}
	log.Printf("[MAPE Summary] Rata-rata MAPE: %.2f%%", avg)
	log.Printf("========================================================")

	return &dto.MAPESummaryResponse{
		TotalDataUji:       len(parsed),
		MAPEValues:         parsed,
		AverageMAPEPercent: avg,
	}, nil
}

func extractPercentErrors(result *dto.UserMAPEResult) []float64 {
	errs := make([]float64, 0, len(result.Cases))
	for _, c := range result.Cases {
		errs = append(errs, c.PercentageErrorPct)
	}
	return errs
}

func average(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}
