package handler

import (
	"net/http"
	"regexp"
	"time"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	ChineseName string    `json:"chinese_name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	ChineseName string `json:"chinese_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

func toUserResponse(user *model.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		ChineseName: user.ChineseName,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求数据",
		})
		return
	}

	// 验证中文姓名
	if !regexp.MustCompile(`^[\p{Han}]{2,10}$`).MatchString(req.ChineseName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "姓名必须是2-10个中文字符",
		})
		return
	}

	// 获取用户ID
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户ID不能为空",
		})
		return
	}

	// 从数据库获取用户
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 更新用户信息
	user.Username = req.Username
	user.ChineseName = req.ChineseName
	user.Email = req.Email
	user.Role = req.Role

	// 保存更新
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户信息更新成功",
		"data":    toUserResponse(&user),
	})
}
