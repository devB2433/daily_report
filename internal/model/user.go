package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username     string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	ChineseName  string     `gorm:"column:chinese_name;size:50;not null" json:"chinese_name"`
	Password     string     `gorm:"size:255;not null" json:"-"`
	PasswordHash string     `gorm:"column:password_hash;size:255;not null" json:"-"`
	Email        string     `gorm:"size:100" json:"email"`
	Role         string     `gorm:"size:20;default:'user'" json:"role"` // admin or user
	Department   string     `gorm:"size:20;not null" json:"department"` // 交付 or 产品研发测试
	Level        string     `gorm:"size:20;default:'初级'" json:"level"`  // 用户级别：初级、中级、高级
	LastLoginAt  *time.Time `gorm:"default:null" json:"last_login_at"`  // 使用指针类型，允许为null
}

// SetPassword 设置用户密码
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = password
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// BeforeCreate GORM的钩子，在创建用户前执行
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = "user"
	}
	return nil
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}
