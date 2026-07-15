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
	DistMap    map[string]float64
	PrevMap    map[string]string
	OSMNodes   map[int64]*model.OSMNode
	SnapCoords map[string][2]float64
	UserLat    float64
	UserLon    float64
	Time       time.Duration
}

type RouteTestReq struct {
	FromLat float64 `query:"from_lat"`
	FromLon float64 `query:"from_lon"`
	ToLat   float64 `query:"to_lat"`
	ToLon   float64 `query:"to_lon"`
	ToName  string  `query:"to_name"`
}

type RouteComputation struct {
	Distance       string
	Path           []PathPoint
	DijkstraCost   float64
	DijkstraTimeMs float64
}
