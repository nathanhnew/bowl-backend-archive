package models

type Location struct {
	City     string   `json:"city" bson:"city"`
	State    string   `json:"state" bson:"state"`
	Geometry Geometry `json:"geometry" bson:"geometry"`
}

type Geometry struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}
