package handler

import (
	"net/http"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
)

// getCurrentUser 从上下文中获取当前用户信息
func getCurrentUser(c *gin.Context) *model.User {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil
	}

	db := database.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil
	}

	return &user
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	// 静态文件
	router.Static("/static", "./web/static")

	// 页面路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", nil)
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register", nil)
	})
	router.GET("/projects", func(c *gin.Context) {
		user := getCurrentUser(c)
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		// 检查是否为管理员
		if user.Role != "admin" {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "projects", gin.H{
			"User": user,
		})
	})
	router.GET("/users", UsersPage)

	// API 路由
	api := router.Group("/api")
	{
		api.POST("/register", RegisterHandler)
		api.POST("/login", LoginHandler)
		api.POST("/logout", LogoutHandler)

		authorized := api.Group("/", AuthMiddleware())
		{
			// 需要管理员权限的API
			adminOnly := authorized.Group("/", RootRequired())
			{
				// 项目管理API
				adminOnly.GET("/projects", GetProjects)
				adminOnly.POST("/projects", CreateProject)
				adminOnly.PUT("/projects/:id", UpdateProject)
				adminOnly.DELETE("/projects/:id", DeleteProject)
				adminOnly.GET("/projects/export", ExportProjects)
				adminOnly.POST("/projects/import", ImportProjects)

				// 用户管理API
				adminOnly.GET("/users", GetUsers)
				adminOnly.POST("/users/:id/reset-password", ResetUserPassword)
				adminOnly.PUT("/users/:id", UpdateUserInfo)
			}
		}
	}
}

// UsersPage 用户管理页面
func UsersPage(c *gin.Context) {
	user := getCurrentUser(c)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否为管理员
	if user.Role != "admin" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "users", gin.H{
		"User": user,
	})
}
