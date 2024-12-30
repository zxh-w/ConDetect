package dto

import (
	"ConDetect/backend/utils/entropyscan"
	"time"
)

type SandflyCreate struct {
	Name             string `json:"name"`
	Path             string `json:"path"`
	InfectedStrategy string `json:"infectedStrategy"`
	InfectedDir      string `json:"infectedDir"`
	Description      string `json:"description"`
}
type SandflyUpdate struct {
	ID uint `json:"id"`

	Name             string `json:"name"`
	Path             string `json:"path"`
	InfectedStrategy string `json:"infectedStrategy"`
	InfectedDir      string `json:"infectedDir"`
	Description      string `json:"description"`
}
type SandflyDelete struct {
	RemoveRecord   bool   `json:"removeRecord"`
	RemoveInfected bool   `json:"removeInfected"`
	Ids            []uint `json:"ids" validate:"required"`
}

type SandflyInfo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	Name             string `json:"name"`
	Path             string `json:"path"`
	InfectedStrategy string `json:"infectedStrategy"`
	InfectedDir      string `json:"infectedDir"`
	LastHandleDate   string `json:"lastHandleDate"`
	Description      string `json:"description"`
}
type SearchSandflyWithPage struct {
	PageInfo
	Info    string `json:"info"`
	OrderBy string `json:"orderBy" validate:"required,oneof=name status created_at"`
	Order   string `json:"order" validate:"required,oneof=null ascending descending"`
}
type SandflyRecordSearch struct {
	PageInfo

	SandflyID uint      `json:"sandflyID"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
type SandflyRecord struct {
	Name         string                 `json:"name"`
	Path         string                 `json:"path"`
	Status       string                 `json:"status"`
	TotalMacious uint                   `json:"totalMacious"`
	CreatedAt    time.Time              `json:"createdAt"`
	Macious      []entropyscan.FileData `json:"macious"`
}
