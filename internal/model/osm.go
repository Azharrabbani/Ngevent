package model

type OSMNode struct {
	ID         int64
	Lat        float64
	Lon        float64
	StreetName string
}

type OSMWay struct {
	Nodes   []int64
	Name    string
	Highway string
	Oneway  bool
	Reverse bool
}
