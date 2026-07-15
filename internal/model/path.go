package model

type Location struct {
	Name string
	Lat  float64
	Lon  float64
}

type Graph map[string][]Edge

type Edge struct {
	To     string
	Weight float64
}
