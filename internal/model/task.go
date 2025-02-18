package model

import (
	"gorm.io/gorm"
)

// Task 任务模型
type Task struct {
	gorm.Model
	ReportID  uint     `gorm:"not null" json:"report_id"`
	Report    *Report  `json:"report,omitempty" gorm:"foreignKey:ReportID"`
	ProjectID uint     `gorm:"not null" json:"project_id"`
	Project   *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Hours     float64  `gorm:"not null" json:"hours"`
	Content   string   `gorm:"type:text;not null" json:"content"`
	Status    string   `gorm:"not null;default:'completed'" json:"status"`
}

// BeforeCreate 创建前的钩子
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.Status == "" {
		t.Status = "completed"
	}
	return nil
}
