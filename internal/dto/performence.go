package dto

import (
	"ngevent/internal/model"
	"time"
)

type HaversinePerf struct {
	Results map[string]float64
	Time    time.Duration
}

type DijkstraPerf struct {
	DistMap  map[string]float64
	PrevMap  map[string]string
	OSMNodes map[int64]*model.OSMNode
	Time     time.Duration
}