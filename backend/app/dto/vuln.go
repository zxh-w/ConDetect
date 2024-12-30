package dto

import (
	"time"
)

type TrivyCreate struct {
	Name        string `json:"name"`
	Type        uint   `json:"type"`
	Target      string `json:"target"`
	Parallel    uint   `json:"parallel"`
	Description string `json:"description"`
}
type TrivyUpdate struct {
	ID uint `json:"id"`

	Name        string `json:"name"`
	Type        uint   `json:"type"`
	Target      string `json:"target"`
	Parallel    uint   `json:"parallel"`
	Description string `json:"description"`
}
type TrivyDelete struct {
	RemoveRecord bool   `json:"removeRecord"`
	Ids          []uint `json:"ids" validate:"required"`
}

type TrivyInfo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	Name           string `json:"name"`
	Type           uint   `json:"type"`
	Target         string `json:"target"`
	Parallel       uint   `json:"parallel"`
	LastHandleDate string `json:"lastHandleDate"`
	Description    string `json:"description"`
}
type SearchTrivyWithPage struct {
	PageInfo
	Info    string `json:"info"`
	OrderBy string `json:"orderBy" validate:"required,oneof=name status created_at"`
	Order   string `json:"order" validate:"required,oneof=null ascending descending"`
}
type TrivyRecordSearch struct {
	PageInfo

	TrivyID   uint      `json:"trivyID"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
type TrivyRecord struct {
	Name            string          `json:"name"`
	Target          string          `json:"target"`
	Status          string          `json:"status"`
	TotalCVE        uint            `json:"totalCVE"`
	CreatedAt       time.Time       `json:"createdAt"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type Vulnerability struct {
	VulnerabilityID  string `json:"vulnerabilityID"`
	Title            string `json:"title"`
	Status           string `json:"status"`
	Description      string `json:"description"`
	PkgName          string `json:"pkgName"`
	PrimaryURL       string `json:"primaryURL"`
	Severity         string `json:"severity"`
	InstalledVersion string `json:"installedVersion"`
	FixedVersion     string `json:"fixedVersion"`
}
