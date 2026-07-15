package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

type SnapResult struct {
	VirtualKey string
	Lat, Lon   float64
	NodeAKey   string
	NodeBKey   string
	DistToA    float64
	DistToB    float64
	IsOneway   bool
	NodeAID    int64
	NodeBID    int64
}

func FetchRoadGraph(minLat, minLon, maxLat, maxLon float64) (map[int64]*model.OSMNode, []model.OSMWay, error) {
	query := fmt.Sprintf(`
    [out:json][timeout:30];
    (

      way["highway"~"^(motorway|trunk|primary|secondary|tertiary|residential|unclassified|primary_link|secondary_link|motorway_link|trunk_link|living_street|service|road)$"]
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

	client := &http.Client{Timeout: 35 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("overpass request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("overpass HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("overpass decode failed: %w", err)
	}

	nodes := make(map[int64]*model.OSMNode)
	var ways []model.OSMWay

	for _, el := range result.Elements {
		switch el.Type {
		case "node":
			nodes[el.ID] = &model.OSMNode{
				ID:  el.ID,
				Lat: el.Lat,
				Lon: el.Lon,
			}
		case "way":
			if len(el.Nodes) < 2 {
				continue
			}
			onewayTag := el.Tags["oneway"]
			isOneway := onewayTag == "yes" || onewayTag == "1" || onewayTag == "true"
			isReverse := onewayTag == "-1" || onewayTag == "reverse"

			if el.Tags["highway"] == "motorway" && onewayTag == "" {
				isOneway = true
			}
			if el.Tags["junction"] == "roundabout" {
				isOneway = true
			}

			ways = append(ways, model.OSMWay{
				Nodes:   el.Nodes,
				Name:    el.Tags["name"],
				Highway: el.Tags["highway"],
				Oneway:  isOneway,
				Reverse: isReverse,
			})
		}
	}

	// Assign street name ke nodes
	for _, way := range ways {
		if way.Name == "" {
			continue
		}
		for _, nodeID := range way.Nodes {
			if n, ok := nodes[nodeID]; ok && n.StreetName == "" {
				n.StreetName = way.Name
			}
		}
	}

	return nodes, ways, nil
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

	pad := paddingKm / 111.0
	minLat -= pad
	maxLat += pad
	minLon -= pad
	maxLon += pad
	return
}

func highwayWeight(highway string) float64 {
	switch highway {

	case "motorway":
		return 10.0

	case "motorway_link":
		return 8.0

	case "trunk":
		return 1.0

	case "trunk_link":
		return 1.05

	case "primary":
		return 1.1

	case "primary_link":
		return 1.15

	case "secondary":
		return 1.3

	case "secondary_link":
		return 1.35

	case "tertiary":
		return 1.6

	case "residential":
		return 2.5

	case "road":
		return 3.0

	case "service":
		return 4.0

	case "living_street":
		return 5.0

	default:
		return 4.0
	}
}

// BuildGraphFromOSM builds a directed graph from OSM data.
// Returns: graph pointer + snapCoords (virtual snap nodes)

func BuildGraphFromOSM(
	user model.Location,
	events []model.Location,
	osmNodes map[int64]*model.OSMNode,
	ways []model.OSMWay,
) (*model.Graph, map[string][2]float64) {
	graph := make(model.Graph)
	snapCoords := make(map[string][2]float64)
	nodesWithEdges := make(map[string]bool)

	if len(osmNodes) == 0 || len(ways) == 0 {
		graph[user.Name] = []model.Edge{}
		for _, event := range events {
			dist := Haversine(user.Lat, user.Lon, event.Lat, event.Lon)
			graph[user.Name] = append(graph[user.Name], model.Edge{To: event.Name, Weight: dist})
			graph[event.Name] = []model.Edge{}
		}
		return &graph, snapCoords
	}

	// Build road edges with oneway rules
	for _, way := range ways {
		nodeIDs := way.Nodes
		if way.Reverse {
			reversed := make([]int64, len(nodeIDs))
			for i, id := range nodeIDs {
				reversed[len(nodeIDs)-1-i] = id
			}
			nodeIDs = reversed
		}

		for i := 0; i < len(nodeIDs)-1; i++ {
			fromID := nodeIDs[i]
			toID := nodeIDs[i+1]
			fromNode, ok1 := osmNodes[fromID]
			toNode, ok2 := osmNodes[toID]
			if !ok1 || !ok2 {
				continue
			}

			fromKey := fmt.Sprintf("osm:%d", fromID)
			toKey := fmt.Sprintf("osm:%d", toID)

			physDist := Haversine(fromNode.Lat, fromNode.Lon, toNode.Lat, toNode.Lon)
			// weightedDist := physDist * highwayWeight(way.Highway) 

			if _, exists := graph[fromKey]; !exists {
				graph[fromKey] = []model.Edge{}
			}
			if _, exists := graph[toKey]; !exists {
				graph[toKey] = []model.Edge{}
			}

			graph[fromKey] = append(graph[fromKey],
				model.Edge{
					To:     toKey,
					Weight: physDist})
			nodesWithEdges[fromKey] = true
			nodesWithEdges[toKey] = true

			if !way.Oneway && !way.Reverse {
				graph[toKey] = append(graph[toKey], model.Edge{To: fromKey, Weight: physDist})
			}
		}
	}

	// Snap user to the nearest road segment
	graph[user.Name] = []model.Edge{}
	if snap, ok := snapToSegment(user.Lat, user.Lon, osmNodes, ways); ok {
		snapCoords[snap.VirtualKey] = [2]float64{snap.Lat, snap.Lon}

		if _, exists := graph[snap.VirtualKey]; !exists {
			graph[snap.VirtualKey] = []model.Edge{}
		}

		// Virtual node can go to both ends (user is not yet on the road, free to choose entry direction)
		graph[snap.VirtualKey] = append(graph[snap.VirtualKey],
			model.Edge{To: snap.NodeAKey, Weight: snap.DistToA},
			model.Edge{To: snap.NodeBKey, Weight: snap.DistToB},
		)

		userSnapDist := Haversine(user.Lat, user.Lon, snap.Lat, snap.Lon)
		graph[user.Name] = append(graph[user.Name], model.Edge{To: snap.VirtualKey, Weight: userSnapDist})
		graph[snap.VirtualKey] = append(graph[snap.VirtualKey], model.Edge{To: user.Name, Weight: userSnapDist})
	}

	// Snap each event to the nearest road segment
	for _, event := range events {
		graph[event.Name] = []model.Edge{}
		if snap, ok := snapToSegment(event.Lat, event.Lon, osmNodes, ways); ok {
			snapCoords[snap.VirtualKey] = [2]float64{snap.Lat, snap.Lon}

			if _, exists := graph[snap.VirtualKey]; !exists {
				graph[snap.VirtualKey] = []model.Edge{}
			}

			// Both ends of the segment can enter the virtual node (event can be reached from two directions)
			graph[snap.NodeAKey] = append(graph[snap.NodeAKey],
				model.Edge{To: snap.VirtualKey, Weight: snap.DistToA},
			)
			if !snap.IsOneway {
				graph[snap.NodeBKey] = append(graph[snap.NodeBKey],
					model.Edge{To: snap.VirtualKey, Weight: snap.DistToB},
				)
			}

			eventSnapDist := Haversine(event.Lat, event.Lon, snap.Lat, snap.Lon)
			graph[snap.VirtualKey] = append(graph[snap.VirtualKey],
				model.Edge{To: event.Name, Weight: eventSnapDist},
			)
			graph[event.Name] = append(graph[event.Name],
				model.Edge{To: snap.VirtualKey, Weight: eventSnapDist},
			)
		}
	}

	for _, way := range ways {
		if way.Name == "Jalan Gerbang Pemuda" {
			log.Printf(
				"Way %s, highway=%s, nodes=%d",
				way.Name,
				way.Highway,
				len(way.Nodes),
			)
		}
	}

	return &graph, snapCoords
}

// snapToSegment finds the nearest road segment to use as a snap point.
func snapToSegment(lat, lon float64, osmNodes map[int64]*model.OSMNode, ways []model.OSMWay) (SnapResult, bool) {
	bestDist := math.Inf(1)
	var best SnapResult
	found := false

	for _, way := range ways {
		nodeIDs := way.Nodes
		if way.Reverse {
			reversed := make([]int64, len(nodeIDs))
			for i, id := range nodeIDs {
				reversed[len(nodeIDs)-1-i] = id
			}
			nodeIDs = reversed
		}

		for i := 0; i < len(nodeIDs)-1; i++ {
			aID := nodeIDs[i]
			bID := nodeIDs[i+1]
			aNode, ok1 := osmNodes[aID]
			bNode, ok2 := osmNodes[bID]
			if !ok1 || !ok2 {
				continue
			}

			projLat, projLon := projectPointToSegment(lat, lon, aNode.Lat, aNode.Lon, bNode.Lat, bNode.Lon)
			realDist := Haversine(lat, lon, projLat, projLon)

			
			if realDist < bestDist {
				bestDist = realDist
				best = SnapResult{
					VirtualKey: fmt.Sprintf("snap:%d:%d", aID, bID),
					Lat:        projLat,
					Lon:        projLon,
					NodeAKey:   fmt.Sprintf("osm:%d", aID),
					NodeBKey:   fmt.Sprintf("osm:%d", bID),
					NodeAID:    aID,
					NodeBID:    bID,
					DistToA:    Haversine(projLat, projLon, aNode.Lat, aNode.Lon),
					DistToB:    Haversine(projLat, projLon, bNode.Lat, bNode.Lon),
					IsOneway:   way.Oneway || way.Reverse,
				}
				found = true
			}
		}
	}

	return best, found
}

func projectPointToSegment(pLat, pLon, aLat, aLon, bLat, bLon float64) (float64, float64) {
	abLat := bLat - aLat
	abLon := bLon - aLon
	apLat := pLat - aLat
	apLon := pLon - aLon

	dot := apLat*abLat + apLon*abLon
	lenSq := abLat*abLat + abLon*abLon
	if lenSq == 0 {
		return aLat, aLon
	}

	t := dot / lenSq
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	return aLat + t*abLat, aLon + t*abLon
}
