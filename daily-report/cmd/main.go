package main

import (
	"fmt"
	"log"

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
		public.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{
				"title": "工作管理系统",
			})
		})

		public.GET("/register", func(c *gin.Context) {
			c.HTML(200, "register.html", gin.H{
				"title": "工作管理系统",
			})
		})

		public.GET("/login", func(c *gin.Context) {
			c.HTML(200, "login.html", gin.H{
				"title": "工作管理系统",
			})
		})
	}

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		api.POST("/register", handler.RegisterHandler)
		api.POST("/login", handler.LoginHandler)

		// 需要认证的路由
		authorized := api.Group("/")
		authorized.Use(handler.AuthMiddleware())
		{
			// 用户相关
			authorized.GET("/user/info", handler.GetUserInfo)

			// 项目列表API（所有认证用户可访问）
			authorized.GET("/projects", handler.GetProjects)

			// 项目管理相关（需要管理员权限）
			projectAdmin := authorized.Group("/projects")
			projectAdmin.Use(handler.RootRequired())
			{
				projectAdmin.POST("", handler.CreateProject)
				projectAdmin.PUT("/:id", handler.UpdateProject)
				projectAdmin.DELETE("/:id", handler.DeleteProject)
				projectAdmin.GET("/:id", handler.GetProject)
			}

			// 日报相关
			reports := authorized.Group("/reports")
			{
				reports.POST("", handler.CreateReport)
				reports.GET("", handler.GetReports)
				reports.GET("/status", handler.GetReportSubmissionStatus)
				reports.GET("/stats/monthly", handler.GetMonthlyStats)
				reports.GET("/:id", handler.GetReport)
				reports.DELETE("/:id", handler.DeleteReport)
			}
		}
	}

	// 需要认证的页面路由
	authorized := r.Group("/")
	authorized.Use(handler.AuthMiddleware())
	{
		// 项目管理页面（需要管理员权限）
		projectPages := authorized.Group("/projects")
		projectPages.Use(handler.RootRequired())
		{
			projectPages.GET("", func(c *gin.Context) {
				c.HTML(200, "projects.html", gin.H{
					"title": "项目管理",
				})
			})
		}

		// 日报页面路由
		authorized.GET("/reports/new", func(c *gin.Context) {
			c.HTML(200, "write_report.html", gin.H{
				"title": "写日报",
			})
		})

		authorized.GET("/reports", func(c *gin.Context) {
			c.HTML(200, "reports.html", gin.H{
				"title": "我的日报",
			})
		})
	}
}
