package handler

import (
	"net/http"
	"regexp"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID          uint   `json:"id"`
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
	}
}

// UpdateUserInfo 更新用户信息（仅管理员可用）
func UpdateUserInfo(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		ChineseName string `json:"chinese_name" binding:"required"`
		Department  string `json:"department" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	// 验证中文姓名（2-10个中文字符）
	if !regexp.MustCompile(`^[\p{Han}]{2,10}$`).MatchString(req.ChineseName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "姓名必须是2-10个中文字符",
		})
		return
	}

	// 验证部门
	if req.Department != "交付" && req.Department != "产品研发测试" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "部门必须是'交付'或'产品研发测试'",
		})
		return
	}

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
	user.ChineseName = req.ChineseName
	user.Department = req.Department

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
		"data":    user,
	})
}
