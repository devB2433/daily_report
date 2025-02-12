package model

import (
	"time"

	"gorm.io/gorm"
)

// Report 日报模型
type Report struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"user_id"`
	Date        time.Time `gorm:"not null" json:"date"`
	Status      string    `gorm:"not null;default:'submitted'" json:"status"` // submitted, draft
	Tasks       []Task    `gorm:"foreignKey:ReportID" json:"tasks"`
	SubmittedAt time.Time `gorm:"not null" json:"submitted_at"`
}

// BeforeCreate 创建前的钩子
func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.Status == "" {
		r.Status = "submitted"
	}
	if r.SubmittedAt.IsZero() {
		r.SubmittedAt = time.Now()
	}
	return nil
}

// Submit 提交日报
func (r *Report) Submit() {
	r.Status = "submitted"
	r.SubmittedAt = time.Now()
}

// IsSubmitted 检查是否已提交
func (r *Report) IsSubmitted() bool {
	return r.Status == "submitted"
}

// IsDraft 检查是否是草稿
func (r *Report) IsDraft() bool {
	return r.Status == "draft"
}
