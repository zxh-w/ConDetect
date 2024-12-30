package dto

import "time"

type CheckBaselineOprate struct {
	Option bool `json:"option"`
}

type SearchBaselineWithPage struct {
	PageInfo
	Level string `json:"level"`
}

type BaselineInfo struct {
	Start   time.Time     `json:"start"`
	End     time.Time     `json:"end"`
	Score   uint          `json:"score"`
	Checks  uint          `json:"checks"`
	Results []CheckResult `json:"results"`
}

type CheckResult struct {
	Id      string `json:"id"`
	Class   string `json:"class"`
	Desc    string `json:"desc"`
	Result  string `json:"result"`
	Details string `json:"details"`
}

type DockerBench struct {
	Start  uint64          `json:"start"`
	End    uint64          `json:"end"`
	Score  uint            `json:"score"`
	Checks uint            `json:"checks"`
	Tests  []BenchTestItem `json:"tests"`
}

type BenchTestItem struct {
	Id      string            `json:"id"`
	Desc    string            `json:"desc"`
	Results []BenchTestResult `json:"results"`
}

type BenchTestResult struct {
	Id      string `json:"id"`
	Desc    string `json:"desc"`
	Result  string `json:"result"`
	Details string `json:"details"`
}
