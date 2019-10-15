package models

type Location struct {
	City     string   `bson:"city"`
	State    string   `bson:"state"`
	Geometry Geometry `bson:"geometry"`
}

type Geometry struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}
