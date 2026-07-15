package utils

import (
	"errors"
	"fmt"
	"log"
	"math"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"strings"
	"time"
)

type Item struct {
	Node     string
	Priority float64
}

type MinHeap struct{ data []Item }

func NewMinHeap() *MinHeap  { return &MinHeap{} }
func (h *MinHeap) Len() int { return len(h.data) }

func (h *MinHeap) Push(item Item) {
	h.data = append(h.data, item)
	h.siftUp(len(h.data) - 1)
}

func (h *MinHeap) Pop() Item {
	min := h.data[0]
	last := len(h.data) - 1
	h.data[0] = h.data[last]
	h.data = h.data[:last]
	if len(h.data) > 0 {
		h.siftDown(0)
	}
	return min
}

func (h *MinHeap) siftUp(i int) {
	for i > 0 {
		p := (i - 1) / 2
		if h.data[i].Priority >= h.data[p].Priority {
			break
		}
		h.data[i], h.data[p] = h.data[p], h.data[i]
		i = p
	}
}

func (h *MinHeap) siftDown(i int) {
	n := len(h.data)
	for {
		s := i
		l, r := 2*i+1, 2*i+2
		if l < n && h.data[l].Priority < h.data[s].Priority {
			s = l
		}
		if r < n && h.data[r].Priority < h.data[s].Priority {
			s = r
		}
		if s == i {
			break
		}
		h.data[i], h.data[s] = h.data[s], h.data[i]
		i = s
	}
}

func Dijkstra(graph model.Graph, start, target string) (map[string]float64, map[string]string) {
	dist := make(map[string]float64, len(graph))
	prev := make(map[string]string, len(graph))
	visited := make(map[string]bool, len(graph))

	for node := range graph {
		dist[node] = math.Inf(1)
	}
	dist[start] = 0

	pq := NewMinHeap()
	pq.Push(Item{Node: start, Priority: 0})

	for pq.Len() > 0 {
		cur := pq.Pop()

		if cur.Node == target {
			break
		}

		if visited[cur.Node] {
			continue
		}
		visited[cur.Node] = true

		for _, edge := range graph[cur.Node] {
			if visited[edge.To] {
				continue
			}
			if nd := dist[cur.Node] + edge.Weight; nd < dist[edge.To] {
				dist[edge.To] = nd
				prev[edge.To] = cur.Node
				pq.Push(Item{Node: edge.To, Priority: nd})
			}
		}
	}
	return dist, prev
}

func BuildRoadPathWithCoords(
	prev map[string]string,
	start, target string,
	osmNodes map[int64]*model.OSMNode,
	snapCoords map[string][2]float64,
	userLat, userLon float64,
	eventLat, eventLon float64,
) (float64, []dto.PathPoint) {
	// 1. Reconstruct the node sequence from prevMap
	var raw []string
	seen := make(map[string]bool)
	cur := target
	for cur != "" && !seen[cur] {
		seen[cur] = true
		raw = append([]string{cur}, raw...)
		next := prev[cur]
		if next == cur {
			break
		}
		cur = next
	}

	if len(raw) == 0 || raw[0] != start {
		fallbackDist := Haversine(userLat, userLon, eventLat, eventLon)
		return fallbackDist, []dto.PathPoint{
			{Name: "user", Lat: userLat, Lon: userLon},
			{Name: target, Lat: eventLat, Lon: eventLon},
		}
	}

	// 2. Resolve each node to coordinates
	type point struct {
		key      string
		name     string
		lat, lon float64
	}

	var points []point
	for _, node := range raw {
		switch {
		case node == start:
			points = append(points, point{key: node, name: "user", lat: userLat, lon: userLon})

		case node == target:
			points = append(points, point{key: node, name: target, lat: eventLat, lon: eventLon})

		case strings.HasPrefix(node, "osm:"):
			var id int64
			fmt.Sscanf(node, "osm:%d", &id)
			if n, ok := osmNodes[id]; ok {
				points = append(points, point{key: node, name: n.StreetName, lat: n.Lat, lon: n.Lon})
			}

		case strings.HasPrefix(node, "snap:"):
			if coords, ok := snapCoords[node]; ok {
				points = append(points, point{key: node, name: "", lat: coords[0], lon: coords[1]})
			}
		}
	}

	// 3. Calculate the actual physical distance from the resolved coordinates
	var realDist float64
	for i := 1; i < len(points); i++ {
		realDist += Haversine(points[i-1].lat, points[i-1].lon, points[i].lat, points[i].lon)
	}

	// 4. Suppress repeated street names, but all coordinates remain included
	var result []dto.PathPoint
	lastStreet := ""
	for _, p := range points {
		display := p.name
		if strings.HasPrefix(p.key, "osm:") || strings.HasPrefix(p.key, "snap:") {
			if p.name == "" || p.name == lastStreet {
				display = ""
			} else {
				lastStreet = p.name
			}
		}
		result = append(result, dto.PathPoint{Name: display, Lat: p.lat, Lon: p.lon})
	}

	return realDist, result
}

func ComputePathToEvent(userLat, userLon float64, eventName string, eventLat, eventLon float64) (*dto.RouteComputation, error) {
	user := model.Location{Name: "user", Lat: userLat, Lon: userLon}
	event := model.Location{Name: eventName, Lat: eventLat, Lon: eventLon}
	events := []model.Location{event}

	directDist := Haversine(userLat, userLon, eventLat, eventLon)
	padding := math.Max(2.0, directDist*0.3)
	minLat, minLon, maxLat, maxLon := BoundingBox(user, events, padding)

	fetchStart := time.Now()
	osmNodes, ways, err := FetchRoadGraph(minLat, minLon, maxLat, maxLon)
	fetchElapsed := time.Since(fetchStart)

	if err != nil {
		log.Printf("[Route] Overpass ERROR: %v", err)
		return nil, fmt.Errorf(" failed to fetch road map data: %w", err)
	}

	if len(osmNodes) == 0 || len(ways) == 0 {
		log.Printf("[Route] Overpass returned empty data: nodes=%d ways=%d", len(osmNodes), len(ways))
		return nil, errors.New("failed to fetch road map data: empty data returned from Overpass")
	}

	log.Printf("[Route] Overpass: nodes=%d ways=%d padding=%.1fkm fetchTime=%.3fs",
		len(osmNodes), len(ways), padding, fetchElapsed.Seconds())

	graph, snapCoords := BuildGraphFromOSM(user, events, osmNodes, ways)

	log.Printf("[Route] Graph: totalNodes=%d userEdges=%d eventEdges=%d",
		len(*graph), len((*graph)[user.Name]), len((*graph)[eventName]))

	algoStart := time.Now()
	distMap, prevMap := Dijkstra(*graph, user.Name, eventName)
	algoElapsed := time.Since(algoStart)

	cost := distMap[eventName]
	log.Printf("[Route] Dijkstra cost=%.4f (weighted) algoTime=%.4fs", cost, algoElapsed.Seconds())

	if math.IsInf(cost, 1) || cost == 0 {
		log.Printf("[Route] Dijkstra unreachable")
		return nil, errors.New("failed to find route to event location")
	}

	realDist, path := BuildRoadPathWithCoords(
		prevMap, user.Name, eventName,
		osmNodes, snapCoords,
		userLat, userLon, eventLat, eventLon,
	)

	log.Printf("[Analytic] ============================================")
	log.Printf("[Analytic] Event              : %s", eventName)
	log.Printf("[Analytic] Jarak fisik rute   : %.2f km", realDist)
	log.Printf("[Analytic] Jarak garis lurus  : %.2f km", directDist)
	log.Printf("[Analytic] Fetch OSM          : %.3f s (I/O jaringan)", fetchElapsed.Seconds())
	log.Printf("[Analytic] Kalkulasi Dijkstra : %.4f s (murni algoritma)", algoElapsed.Seconds())
	log.Printf("[Analytic] Total nodes graph  : %d", len(*graph))
	log.Printf("[Analytic] Waypoints di rute  : %d titik", len(path))
	log.Printf("[Analytic] ============================================")

	return &dto.RouteComputation{
		Distance:       fmt.Sprintf("%.1f km", realDist),
		Path:           path,
		DijkstraCost:   cost,
		DijkstraTimeMs: float64(algoElapsed.Microseconds()) / 1000.0,
	}, nil
}
