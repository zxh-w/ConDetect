package model

type MeasureRecord struct {
	BaseModel
	Name     string  `json:"name"`
	Baseline float64 `json:"baseline"`
	Vuln     float64 `json:"vuln"`
	Macious  float64 `json:"macious"`
	Measure  float64 `json:"measure"`
}
