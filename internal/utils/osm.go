package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"ngevent/internal/model"
	"time"
)

type OverpassResponse struct {
	Elements []OverpassElement `json:"elements"`
}

type OverpassElement struct {
	Type  string            `json:"type"`
	ID    int64             `json:"id"`
	Lat   float64           `json:"lat"`
	Lon   float64           `json:"lon"`
	Nodes []int64           `json:"nodes"`
	Tags  map[string]string `json:"tags"`
}

// FetchRoadGraph queries Overpass API for road topology within a bounding box.
func FetchRoadGraph(minLat, minLon, maxLat, maxLon float64) (map[int64]*model.OSMNode, [][]int64, error) {
	query := fmt.Sprintf(`
		[out:json][timeout:30];
		(
		  way["highway"~"^(motorway|trunk|primary|secondary|tertiary|residential|unclassified|living_street)$"]
		     (%f,%f,%f,%f);
		);
		(._;>;);
		out body;
	`, minLat, minLon, maxLat, maxLon)

	apiURL := "https://overpass-api.de/api/interpreter?data=" + url.QueryEscape(query)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("overpass build request failed: %w", err)
	}

	req.Header.Set("User-Agent", "ngevent-app/1.0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 35 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("overpass request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("overpass returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("overpass decode failed: %w", err)
	}

	nodes := make(map[int64]*model.OSMNode)
	var ways [][]int64

	for _, el := range result.Elements {
		switch el.Type {
		case "node":
			nodes[el.ID] = &model.OSMNode{
				ID:  el.ID,
				Lat: el.Lat,
				Lon: el.Lon,
			}
		case "way":
			if len(el.Nodes) > 0 {
				ways = append(ways, el.Nodes)
			}
		}
	}

	// Assign street names from ways onto their nodes
	for _, el := range result.Elements {
		if el.Type != "way" {
			continue
		}
		streetName := el.Tags["name"]
		if streetName == "" {
			continue
		}
		for _, nodeID := range el.Nodes {
			if n, ok := nodes[nodeID]; ok && n.StreetName == "" {
				n.StreetName = streetName
			}
		}
	}

	return nodes, ways, nil
}

// SnapToGraph finds the nearest OSM node ID to a given lat/lon.
func SnapToGraph(lat, lon float64, nodes map[int64]*model.OSMNode) (int64, bool) {
	if len(nodes) == 0 {
		return 0, false
	}

	minDist := math.Inf(1)
	var nearest int64

	for id, node := range nodes {
		d := Haversine(lat, lon, node.Lat, node.Lon)
		if d < minDist {
			minDist = d
			nearest = id
		}
	}

	return nearest, true
}

func BoundingBox(user model.Location, events []model.Location, paddingKm float64) (minLat, minLon, maxLat, maxLon float64) {
	minLat = user.Lat
	maxLat = user.Lat
	minLon = user.Lon
	maxLon = user.Lon

	for _, e := range events {
		if e.Lat < minLat {
			minLat = e.Lat
		}
		if e.Lat > maxLat {
			maxLat = e.Lat
		}
		if e.Lon < minLon {
			minLon = e.Lon
		}
		if e.Lon > maxLon {
			maxLon = e.Lon
		}
	}

	// ~1 degree lat ≈ 111 km
	pad := paddingKm / 111.0
	minLat -= pad
	maxLat += pad
	minLon -= pad
	maxLon += pad

	return
}

// BuildGraphFromOSM constructs a full road graph from OSM nodes/ways,
// then injects user and event nodes snapped to their nearest road nodes.
func BuildGraphFromOSM(
	user model.Location,
	events []model.Location,
	osmNodes map[int64]*model.OSMNode,
	ways [][]int64,
) *model.Graph {
	graph := make(model.Graph)

	// If no OSM data, fall back to Haversine star graph ───
	if len(osmNodes) == 0 {
		graph[user.Name] = []model.Edge{}
		for _, event := range events {
			dist := Haversine(user.Lat, user.Lon, event.Lat, event.Lon)
			graph[user.Name] = append(graph[user.Name], model.Edge{
				To:     event.Name,
				Weight: dist,
			})
			graph[event.Name] = []model.Edge{}
		}
		return &graph
	}

	// Build road edges between consecutive nodes in each way (bidirectional)
	for _, way := range ways {
		for i := 0; i < len(way)-1; i++ {
			fromID := way[i]
			toID := way[i+1]

			fromNode, ok1 := osmNodes[fromID]
			toNode, ok2 := osmNodes[toID]
			if !ok1 || !ok2 {
				continue
			}

			fromKey := fmt.Sprintf("osm:%d", fromID)
			toKey := fmt.Sprintf("osm:%d", toID)
			dist := Haversine(fromNode.Lat, fromNode.Lon, toNode.Lat, toNode.Lon)

			graph[fromKey] = append(graph[fromKey], model.Edge{To: toKey, Weight: dist})
			graph[toKey] = append(graph[toKey], model.Edge{To: fromKey, Weight: dist})
		}
	}

	// Snap user to nearest road node
	if userSnapID, ok := SnapToGraph(user.Lat, user.Lon, osmNodes); ok {
		userSnapKey := fmt.Sprintf("osm:%d", userSnapID)
		userSnapNode := osmNodes[userSnapID]
		snapDistUser := Haversine(user.Lat, user.Lon, userSnapNode.Lat, userSnapNode.Lon)

		graph[user.Name] = []model.Edge{{To: userSnapKey, Weight: snapDistUser}}
		graph[userSnapKey] = append(graph[userSnapKey], model.Edge{To: user.Name, Weight: snapDistUser})
	} else {
		graph[user.Name] = []model.Edge{}
	}

	// Snap each event to nearest road node
	for _, event := range events {
		if eventSnapID, ok := SnapToGraph(event.Lat, event.Lon, osmNodes); ok {
			eventSnapKey := fmt.Sprintf("osm:%d", eventSnapID)
			eventSnapNode := osmNodes[eventSnapID]
			snapDistEvent := Haversine(event.Lat, event.Lon, eventSnapNode.Lat, eventSnapNode.Lon)

			graph[event.Name] = []model.Edge{{To: eventSnapKey, Weight: snapDistEvent}}
			graph[eventSnapKey] = append(graph[eventSnapKey], model.Edge{To: event.Name, Weight: snapDistEvent})
		} else {
			graph[event.Name] = []model.Edge{}
		}
	}

	return &graph
}
