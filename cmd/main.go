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
			// 如果已经登录，重定向到首页
			token, _ := c.Cookie("token")
			if token != "" && handler.ValidateToken(token) {
				c.Redirect(http.StatusFound, "/")
				c.Abort()
				return
			}
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

			// 统计分析相关（需要管理员权限）
			analytics := authorized.Group("/analytics")
			analytics.Use(handler.RootRequired())
			{
				analytics.GET("/summary", handler.GetAnalyticsSummary)
			}
		}
	}

	// 需要认证的页面路由
	authorized := r.Group("/")
	authorized.Use(handler.AuthMiddleware())
	{
		// 首页路由
		authorized.GET("/", func(c *gin.Context) {
			username, _ := c.Get("username")
			role, _ := c.Get("role")
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title":    "工作管理系统",
				"username": username,
				"role":     role,
				"isAdmin":  role == "admin",
			})
		})

		// 项目管理页面（需要管理员权限）
		projectPages := authorized.Group("/projects")
		projectPages.Use(handler.RootRequired())
		{
			projectPages.GET("", func(c *gin.Context) {
				c.HTML(http.StatusOK, "projects.html", gin.H{
					"title": "项目管理",
				})
			})
		}

		// 统计分析页面（需要管理员权限）
		analyticsPages := authorized.Group("/analytics")
		analyticsPages.Use(handler.RootRequired())
		{
			analyticsPages.GET("", func(c *gin.Context) {
				c.HTML(http.StatusOK, "analytics.html", gin.H{
					"title": "统计分析",
				})
			})
		}

		// 日报页面路由
		authorized.GET("/reports/new", func(c *gin.Context) {
			c.HTML(http.StatusOK, "write_report.html", gin.H{
				"title": "写日报",
			})
		})

		authorized.GET("/reports", func(c *gin.Context) {
			c.HTML(http.StatusOK, "reports.html", gin.H{
				"title": "我的日报",
			})
		})
	}
}
