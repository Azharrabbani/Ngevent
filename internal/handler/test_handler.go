package handler

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ngevent/internal/dto"
	"ngevent/internal/service"
	"ngevent/internal/utils"

	"github.com/gofiber/fiber/v2"
)

const reportOutputDir = "./reports"

func DijkstraAccuracyTestHandler(c *fiber.Ctx) error {
	log.Println("[Handler] Memulai pengujian akurasi & kinerja Dijkstra...")

	report, err := service.RunDijkstraAccuracyTest()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	timestamp := time.Now().Format("20060102_150405")

	jsonPath, err := utils.SaveAccuracyReportJSON(report, reportOutputDir, timestamp)
	if err != nil {
		log.Printf("[WARN] Gagal menyimpan laporan JSON: %v", err)
	}

	detailCSVPath, summaryCSVPath, err := utils.SaveAccuracyReportCSV(report, reportOutputDir, timestamp)
	if err != nil {
		log.Printf("[WARN] Gagal menyimpan laporan CSV: %v", err)
	}

	log.Printf("[Handler] Laporan disimpan -> json: %s | detail_csv: %s | summary_csv: %s",
		jsonPath, detailCSVPath, summaryCSVPath)

	jsonFile := filepath.Base(jsonPath)
	detailCSVFile := filepath.Base(detailCSVPath)
	summaryCSVFile := filepath.Base(summaryCSVPath)

	return c.JSON(fiber.Map{
		"message":              "Pengujian akurasi & kinerja Dijkstra selesai",
		"total_users":          report.TotalUsers,
		"total_events":         report.TotalEvents,
		"overall_mape_percent": report.OverallMAPEPercent,
		"per_user_summary":     report.Results,
		"files": fiber.Map{
			"json_download_url":        "/reports/download/" + jsonFile,
			"csv_detail_download_url":  "/reports/download/" + detailCSVFile,
			"csv_summary_download_url": "/reports/download/" + summaryCSVFile,
		},
	})
}

func DownloadReportHandler(c *fiber.Ctx) error {
	filename := c.Params("filename")

	if filename == "" || strings.ContainsAny(filename, "/\\") || strings.Contains(filename, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "nama file tidak valid",
		})
	}

	fullPath := filepath.Join(reportOutputDir, filename)

	if _, err := os.Stat(fullPath); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file laporan tidak ditemukan: " + filename,
		})
	}

	return c.Download(fullPath, filename)
}

func SingleDijkstraAccuracyTestHandler(c *fiber.Ctx) error {
	indexParam := c.Params("index")

	index, err := strconv.Atoi(indexParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "parameter index harus berupa angka, contoh: /test/dijkstra-accuracy/single/0",
		})
	}

	log.Printf("[Handler] Memulai pengujian akurasi Dijkstra untuk index=%d", index)

	result, err := service.RunSingleUserAccuracyTest(index)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func MAPESummaryHandler(c *fiber.Ctx) error {
	var req dto.MAPESummaryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "body tidak valid, contoh: {\"mape_values\": [\"4.12\", \"3.87\"]}",
		})
	}

	result, err := service.CalculateAverageMAPE(req.MAPEValues)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
