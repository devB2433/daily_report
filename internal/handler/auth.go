package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"daily-report/internal/config"
	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// 获取JWT密钥
func getJWTSecret() []byte {
	return []byte(config.GetConfig().JWT.Secret)
}

// ValidateToken 验证token是否有效
func ValidateToken(tokenString string) bool {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return false
	}

	// 检查token是否过期
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return false
		}
	}
	return true
}

// ParseToken 解析JWT token并返回claims
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	ChineseName string `json:"chineseName" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Department  string `json:"department" binding:"required"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

// RegisterHandler 处理用户注册
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效，请检查所有字段是否填写正确",
		})
		return
	}

	// 验证邮箱域名
	if !strings.HasSuffix(req.Email, "@blingsec.cn") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请使用 @blingsec.cn 结尾的邮箱",
		})
		return
	}

	// 验证用户名（只允许英文字母）
	if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户名只能包含英文字母",
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

	// 验证密码强度
	if !isPasswordStrong(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "密码必须包含大小写字母、数字和特殊字符，且长度至少为8位",
		})
		return
	}

	db := database.GetDB()

	// 检查邮箱是否已被注册
	var existingUser model.User
	err := db.Where("email = ?", req.Email).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "系统错误，请稍后重试",
		})
		return
	}
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "该邮箱已被注册",
		})
		return
	}

	// 检查用户名是否已被使用
	err = db.Where("username = ?", req.Username).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "系统错误，请稍后重试",
		})
		return
	}
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "该用户名已被使用",
		})
		return
	}

	// 创建新用户
	user := model.User{
		Email:       req.Email,
		Username:    req.Username,
		ChineseName: req.ChineseName,
		Department:  req.Department,
		Role:        "user", // 默认角色为普通用户
	}

	// 如果用户名是admin或root，设置角色为管理员
	if req.Username == "admin" || req.Username == "root" {
		user.Role = "admin"
	}

	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "系统错误：密码加密失败",
		})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("用户创建失败：%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "注册成功",
	})
}

// LoginHandler 处理用户登录
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// log.Printf("登录请求数据无效: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	// log.Printf("收到登录请求 - 邮箱: %s", req.Email)

	// 验证邮箱域名
	if !strings.HasSuffix(req.Email, "@blingsec.cn") {
		// log.Printf("邮箱域名无效: %s", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请使用 @blingsec.cn 结尾的邮箱",
		})
		return
	}

	db := database.GetDB()

	// 查找用户
	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// log.Printf("用户不存在: %s, 错误: %v", req.Email, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// log.Printf("找到用户 - ID: %d, 用户名: %s, 邮箱: %s, 密码哈希: %s",
	//     user.ID, user.Username, user.Email, user.PasswordHash)

	// 验证密码
	if !user.CheckPassword(req.Password) {
		// log.Printf("密码验证失败 - 用户: %s, 输入的密码: %s", user.Username, req.Password)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "密码错误",
		})
		return
	}

	// log.Printf("密码验证成功 - 用户: %s", user.Username)

	// 更新最后登录时间
	user.UpdateLastLogin()
	if err := db.Save(&user).Error; err != nil {
		// log.Printf("更新最后登录时间失败: %v", err)
	}

	// 设置登录Cookie
	maxAge := 24 * 3600 // 1天
	if req.RememberMe {
		maxAge = 7 * 24 * 3600 // 7天
	}

	// 设置用户信息到cookie
	c.SetCookie("user_id", fmt.Sprintf("%d", user.ID), maxAge, "/", "", false, true)
	c.SetCookie("username", user.Username, maxAge, "/", "", false, true)
	c.SetCookie("role", user.Role, maxAge, "/", "", false, true)

	// log.Printf("登录成功 - 用户: %s, 角色: %s", user.Username, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不需要认证的路径
		publicPaths := []string{"/api/login", "/api/register", "/api/logout"}
		currentPath := c.Request.URL.Path

		// 检查是否是公开路径
		for _, path := range publicPaths {
			if path == currentPath {
				c.Next()
				return
			}
		}

		// 从Cookie中获取用户信息
		userId, err := c.Cookie("user_id")
		if err != nil || userId == "" {
			// 如果是API请求，返回JSON响应
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "未登录",
				})
			} else {
				// 如果是页面请求，重定向到登录页面
				c.Redirect(http.StatusFound, "/login")
			}
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		username, _ := c.Cookie("username")
		role, _ := c.Cookie("role")

		c.Set("user_id", userId)
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

// isPasswordStrong 检查密码强度
func isPasswordStrong(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// GetUserInfo 获取当前登录用户信息
func GetUserInfo(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	db := database.GetDB()
	var user model.User
	if err := db.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user_id":      user.ID,
			"username":     user.Username,
			"chinese_name": user.ChineseName,
			"email":        user.Email,
			"role":         user.Role,
		},
	})
}

// LogoutHandler 处理用户退出登录
func LogoutHandler(c *gin.Context) {
	// 清除所有认证相关的cookie
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.SetCookie("username", "", -1, "/", "", false, true)
	c.SetCookie("role", "", -1, "/", "", false, true)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "退出登录成功",
	})
}

// GetUsers 获取用户列表（仅管理员可用）
func GetUsers(c *gin.Context) {
	db := database.GetDB()
	var users []model.User

	// 查询所有用户，按用户名排序
	if err := db.Select("id, username, chinese_name, email, role, department, level, last_login_at").Order("username ASC").Find(&users).Error; err != nil {
		log.Printf("Failed to get users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// ResetUserPassword 重置用户密码（仅管理员可用）
func ResetUserPassword(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	// 验证密码强度
	if !isPasswordStrong(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "密码必须包含大小写字母、数字和特殊字符，且长度至少为8位",
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

	// 设置新密码
	if err := user.SetPassword(req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "密码重置失败",
		})
		return
	}

	// 保存更改
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "密码重置失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "密码重置成功",
	})
}
