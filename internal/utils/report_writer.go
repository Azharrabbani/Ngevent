package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"ngevent/internal/dto"
)


func SaveAccuracyReportJSON(report *dto.AccuracyReport, dir, timestamp string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori %s: %w", dir, err)
	}

	path := filepath.Join(dir, fmt.Sprintf("dijkstra_accuracy_%s.json", timestamp))

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("gagal marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("gagal menulis file JSON: %w", err)
	}

	return path, nil
}


func SaveAccuracyReportCSV(report *dto.AccuracyReport, dir, timestamp string) (detailPath, summaryPath string, err error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", fmt.Errorf("gagal membuat direktori %s: %w", dir, err)
	}

	detailPath = filepath.Join(dir, fmt.Sprintf("dijkstra_accuracy_detail_%s.csv", timestamp))
	summaryPath = filepath.Join(dir, fmt.Sprintf("dijkstra_accuracy_summary_%s.csv", timestamp))

	if err := writeDetailCSV(report, detailPath); err != nil {
		return "", "", err
	}
	if err := writeSummaryCSV(report, summaryPath); err != nil {
		return "", "", err
	}

	return detailPath, summaryPath, nil
}

func writeDetailCSV(report *dto.AccuracyReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("gagal membuat file detail CSV: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	_ = w.Write([]string{
		"user_name", "event_name", "dijkstra_distance", "osrm_distance",
		"percentage_error", "dijkstra_exec_time",
	})

	for _, ur := range report.Results {
		for _, c := range ur.Cases {
			_ = w.Write([]string{
				c.UserName,
				c.EventName,
				c.DijkstraDistance,
				c.OSRMDistance,
				c.PercentageError,
				c.DijkstraExecTime,
			})
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("gagal menulis detail CSV: %w", err)
	}
	return nil
}

func writeSummaryCSV(report *dto.AccuracyReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("gagal membuat file summary CSV: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	_ = w.Write([]string{"user_name", "jumlah_pasangan_valid", "mape_percent", "total_exec_time"})

	totalCases := 0
	for _, ur := range report.Results {
		_ = w.Write([]string{
			ur.UserName,
			fmt.Sprintf("%d", len(ur.Cases)),
			ur.MAPEPercent,
			ur.TotalExecTime,
		})
		totalCases += len(ur.Cases)
	}

	_ = w.Write([]string{"OVERALL", fmt.Sprintf("%d", totalCases), fmt.Sprintf("%.2f%%", report.OverallMAPEPercent), ""})

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("gagal menulis summary CSV: %w", err)
	}
	return nil
}
