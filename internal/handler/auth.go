package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// JWT密钥
const jwtSecret = "your-secret-key"

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
		return []byte(jwtSecret), nil
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
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
		Email:    req.Email,
		Username: req.Username,
		Role:     "user", // 默认角色为普通用户
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
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

	db := database.GetDB()

	// 查找用户
	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "密码错误",
		})
		return
	}

	// 更新最后登录时间
	user.UpdateLastLogin()
	db.Save(&user)

	// 生成JWT令牌
	expirationTime := time.Now().Add(24 * time.Hour)
	if req.RememberMe {
		expirationTime = time.Now().Add(7 * 24 * time.Hour)
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成令牌失败",
		})
		return
	}

	// 设置Cookie
	c.SetCookie("token", tokenString, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"token":    tokenString,
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
		// 从Cookie中获取令牌
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			// 清除可能存在的无效cookie
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未登录",
			})
			c.Abort()
			return
		}

		// 解析并验证token
		claims, err := ParseToken(tokenString)
		if err != nil {
			// token无效，清除cookie
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "登录已过期，请重新登录",
			})
			c.Abort()
			return
		}

		// 检查token是否过期
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.SetCookie("token", "", -1, "/", "", false, true)
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "登录已过期，请重新登录",
				})
				c.Abort()
				return
			}
		}

		// 设置用户信息到上下文
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("email", claims["email"].(string))
		c.Set("username", claims["username"].(string))
		c.Set("role", claims["role"].(string))
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
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}
