package model

import (
	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;uniqueIndex" json:"name"`       // 项目名称
	Code        string `gorm:"size:50;not null;uniqueIndex" json:"code"`        // 项目代号
	Description string `gorm:"size:500" json:"description"`                     // 项目描述
	Status      string `gorm:"size:20;not null;default:'active'" json:"status"` // 项目状态：active, completed, suspended
	StartDate   string `gorm:"size:10" json:"start_date"`                       // 开始日期
	EndDate     string `gorm:"size:10" json:"end_date"`                         // 结束日期
	Manager     string `gorm:"size:50" json:"manager"`                          // 项目经理
	Client      string `gorm:"size:100" json:"client"`                          // 客户名称
}

// BeforeCreate GORM的钩子，在创建项目前执行
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.Status == "" {
		p.Status = "active"
	}
	return nil
}
