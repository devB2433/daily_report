package model

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username     string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password     string     `gorm:"size:255;not null" json:"-"`
	PasswordHash string     `gorm:"column:password_hash;size:255;not null" json:"-"`
	Email        string     `gorm:"size:100" json:"email"`
	Role         string     `gorm:"size:20;default:'user'" json:"role"` // admin or user
	LastLoginAt  *time.Time `gorm:"default:null" json:"last_login_at"`  // 使用指针类型，允许为null
}

// SetPassword 设置用户密码
func (u *User) SetPassword(password string) error {
	log.Printf("正在为用户 %s 设置密码", u.Username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("生成密码哈希失败: %v", err)
		return err
	}
	log.Printf("生成的密码哈希: %s", string(hashedPassword))
	u.Password = password
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	log.Printf("正在验证密码 - 用户: %s", u.Username)
	log.Printf("存储的密码哈希: %s", u.PasswordHash)
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("密码验证失败 - 错误: %v", err)
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
