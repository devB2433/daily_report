package main

import (
	"fmt"
	"log"
	"net/http"

	"daily-report/internal/config"
	"daily-report/internal/database"
	"daily-report/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.GetConfig()

	// 设置gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	database.GetDB()

	// 创建Gin引擎实例
	r := gin.Default()

	// 设置模板目录
	r.LoadHTMLGlob("web/templates/*")

	// 设置静态文件目录
	r.Static("/static", "web/static")

	// 设置路由
	setupRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
	log.Printf("服务器启动在 http://localhost:%d\n", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// 公开路由
	public := r.Group("/")
	{
		// 登录页面
		public.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title": "工作管理系统",
			})
		})

		// 注册页面
		public.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", gin.H{
				"title": "工作管理系统",
			})
		})
	}

	// API路由组
	api := r.Group("/api")
	{
		// 公开的API路由
		api.POST("/register", handler.RegisterHandler)
		api.POST("/login", handler.LoginHandler)
		api.POST("/logout", handler.LogoutHandler)

		// 需要认证的API路由
		authorized := api.Group("/", handler.AuthMiddleware())
		{
			authorized.GET("/user/info", handler.GetUserInfo)
			authorized.GET("/projects", handler.GetProjects)
			authorized.POST("/projects", handler.CreateProject)
			authorized.PUT("/projects/:id", handler.UpdateProject)
			authorized.DELETE("/projects/:id", handler.DeleteProject)
			authorized.GET("/projects/:id", handler.GetProject)
			authorized.POST("/reports", handler.CreateReport)
			authorized.GET("/reports", handler.GetReports)
			authorized.GET("/reports/status", handler.GetReportSubmissionStatus)
			authorized.GET("/reports/stats/monthly", handler.GetMonthlyStats)
			authorized.GET("/reports/:id", handler.GetReport)
			authorized.DELETE("/reports/:id", handler.DeleteReport)
			authorized.GET("/analytics/summary", handler.GetAnalyticsSummary)
		}
	}

	// 页面路由
	r.GET("/", handler.AuthMiddleware(), func(c *gin.Context) {
		username, _ := c.Cookie("username")
		role, _ := c.Cookie("role")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "工作管理系统",
			"username": username,
			"role":     role,
			"isAdmin":  role == "admin",
		})
	})

	r.GET("/projects", handler.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "projects.html", gin.H{
			"title": "项目管理",
		})
	})

	r.GET("/analytics", handler.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "analytics.html", gin.H{
			"title": "统计分析",
		})
	})

	r.GET("/reports/new", handler.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "write_report.html", gin.H{
			"title": "写日报",
		})
	})

	r.GET("/reports", handler.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "reports.html", gin.H{
			"title": "我的日报",
		})
	})
}
