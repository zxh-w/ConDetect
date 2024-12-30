package dto

type SearchMeasureWithPage struct {
	PageInfo
	Info    string `json:"info"`
	OrderBy string `json:"orderBy" validate:"required,oneof=name status created_at"`
	Order   string `json:"order" validate:"required,oneof=null ascending descending"`
}

type MeasureResult struct {
	Name           string `json:"name"`
	Baseline       string `json:"baseline"`
	Vuln           string `json:"vuln"`
	Macious        string `json:"macious"`
	Measure        string `json:"measure"`
	LastHandleDate string `json:"lastHandleDate"`
}
