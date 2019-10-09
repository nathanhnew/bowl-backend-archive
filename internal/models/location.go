package models

type location struct {
	City     string   `bson: "city"`
	State    string   `bson: "state"`
	Geometry geometry `bson: "geometry"`
}

type geometry struct {
	Type        string    `bson: "type"`
	coordinates []float64 `bson: "coordinates"`
}
