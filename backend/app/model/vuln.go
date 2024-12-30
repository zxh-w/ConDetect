package model

type Trivy struct {
	BaseModel

	Name        string `gorm:"type:varchar(64);not null" json:"name"`
	Type        int    `gorm:"type:varchar(64)" json:"type"`
	Target      string `gorm:"type:varchar(64);not null" json:"target"`
	Parallel    int    `gorm:"type:varchar(64)" json:"parallel"`
	Description string `gorm:"type:varchar(64)" json:"description"`
}
