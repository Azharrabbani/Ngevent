package utils

import (
	"container/heap"
	"fmt"
	"math"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"strings"
)

type Item struct {
	Node     string
	Priority float64
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]

	return item
}

func Dijkstra(graph model.Graph, start string) (map[string]float64, map[string]string) {
	dist := make(map[string]float64)
	prev := make(map[string]string)

	for node := range graph {
		dist[node] = math.Inf(1)
	}
	dist[start] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{Node: start, Priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item)

		for _, edge := range graph[current.Node] {
			newDist := dist[current.Node] + edge.Weight

			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				prev[edge.To] = current.Node

				heap.Push(pq, &Item{
					Node:     edge.To,
					Priority: newDist,
				})
			}
		}
	}

	return dist, prev
}

func BuildRoadPathWithCoords(prev map[string]string, start, target string, osmNodes map[int64]*model.OSMNode, userLat, userLon, eventLat, eventLon float64) []dto.PathPoint {
	// 1. Reconstruct raw node from prev map
	var raw []string
	cur := target
	for cur != "" {
		raw = append([]string{cur}, raw...)
		cur = prev[cur]
	}

	if len(raw) == 0 || raw[0] != start {
		return []dto.PathPoint {
			{Name: start, Lat: userLat, Lon: userLon},
			{Name: target, Lat: eventLat, Lon: eventLon},
		}
	}

	// 2. Resolve each node to name and coordinate
	var resolved []dto.PathPoint
	for _, node := range raw {
		if strings.HasPrefix(node, "osm:") {
			var id int64
			fmt.Sscanf(node, "osm:%d", &id)
			if n, ok := osmNodes[id]; ok && n.StreetName != "" {
				resolved = append(resolved, dto.PathPoint{
					Name: n.StreetName,
					Lat: n.Lat,
					Lon: n.Lon,
				})
			}
		} else if node == start {
			resolved = append(resolved, dto.PathPoint{Name: "user", Lat: userLat, Lon: userLon})
		} else if node == target {
			resolved = append(resolved, dto.PathPoint{Name: target, Lat: eventLat, Lon: eventLon})
		}
	}

	// 3. Deduplicate display name, but keep the coordinate
	lastName := ""
	for i := range resolved {
		if resolved[i].Name != "" && resolved[i].Name != lastName {
			lastName = resolved[i].Name
		} else if resolved[i].Name == lastName {
			resolved[i].Name = "" // Street name same, suppress repeat label
		}
	}

	return resolved
}

func ComputePathToEvent(userLat, userLon float64, eventName string, eventLat, eventLon float64) (string, []dto.PathPoint) {
	user := model.Location{
		Name: "user",
		Lat:  userLat,
		Lon:  userLon,
	}

	event := model.Location{
		Name: eventName,
		Lat:  eventLat,
		Lon:  eventLon,
	}

	events := []model.Location{event}

	// Bounding box around user and event with 2km padding
	minLat, minLon, maxLat, maxLon := BoundingBox(user, events, 2.0)

	osmNodes, ways, err := FetchRoadGraph(minLat, minLon, maxLat, maxLon)
	if err != nil {
		fmt.Printf("ComputePathToEvent: Overpass ERROR: %v\n", err)

		// Fallback straight-line distance
		hav := Haversine(userLat, userLon, eventLat, eventLon)
		return fmt.Sprintf("%.2f km", hav), []dto.PathPoint{
			{Name: "user",      Lat: userLat,  Lon: userLon},
			{Name: eventName,   Lat: eventLat, Lon: eventLon},
		}
	}

	graph := BuildGraphFromOSM(user, events, osmNodes, ways)
	distMap, prevMap := Dijkstra(*graph, user.Name)

	dist := distMap[eventName]
	distStr := fmt.Sprintf("%.2f km", dist)

	path := BuildRoadPathWithCoords(prevMap, user.Name, eventName, osmNodes, userLat, userLon, eventLat, eventLon)

	return distStr, path
}
