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

	// 设置日志格式
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s %s | %s\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
			param.Keys["log"],
		)
	}))

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

			// 管理员专用路由
			adminOnly := authorized.Group("/", handler.RootRequired())
			{
				adminOnly.POST("/projects", handler.CreateProject)
				adminOnly.PUT("/projects/:id", handler.UpdateProject)
				adminOnly.DELETE("/projects/:id", handler.DeleteProject)
				adminOnly.GET("/users", handler.GetUsers)
				adminOnly.POST("/users/:id/reset-password", handler.ResetUserPassword)
			}

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

		// 添加日志
		log.Printf("用户信息 - 用户名: %s, 角色: %s", username, role)

		data := gin.H{
			"title": "工作管理系统",
			"User": gin.H{
				"Username": username,
				"Role":     role,
			},
			"isAdmin": role == "admin",
		}

		// 打印传递给模板的数据
		log.Printf("传递给模板的数据: %+v", data)

		c.HTML(http.StatusOK, "index.html", data)
	})

	r.GET("/projects", handler.AuthMiddleware(), func(c *gin.Context) {
		username, _ := c.Cookie("username")
		role, _ := c.Cookie("role")

		c.HTML(http.StatusOK, "projects.html", gin.H{
			"title": "项目管理",
			"User": gin.H{
				"Username": username,
				"Role":     role,
			},
			"isAdmin": role == "admin",
		})
	})

	r.GET("/analytics", handler.AuthMiddleware(), func(c *gin.Context) {
		username, _ := c.Cookie("username")
		role, _ := c.Cookie("role")

		c.HTML(http.StatusOK, "analytics.html", gin.H{
			"title": "统计分析",
			"User": gin.H{
				"Username": username,
				"Role":     role,
			},
			"isAdmin": role == "admin",
		})
	})

	r.GET("/reports/new", handler.AuthMiddleware(), func(c *gin.Context) {
		username, _ := c.Cookie("username")
		role, _ := c.Cookie("role")

		c.HTML(http.StatusOK, "write_report.html", gin.H{
			"title": "写日报",
			"User": gin.H{
				"Username": username,
				"Role":     role,
			},
			"isAdmin": role == "admin",
		})
	})

	r.GET("/reports", handler.AuthMiddleware(), func(c *gin.Context) {
		username, _ := c.Cookie("username")
		role, _ := c.Cookie("role")

		c.HTML(http.StatusOK, "reports.html", gin.H{
			"title": "我的日报",
			"User": gin.H{
				"Username": username,
				"Role":     role,
			},
			"isAdmin": role == "admin",
		})
	})

	// 修改用户管理页面路由，与项目管理保持一致
	r.GET("/users", handler.AuthMiddleware(), func(c *gin.Context) {
		username, _ := c.Cookie("username")
		role, _ := c.Cookie("role")

		// 检查是否为管理员
		if role != "admin" {
			c.Redirect(http.StatusFound, "/")
			return
		}

		c.HTML(http.StatusOK, "users", gin.H{
			"title": "用户管理",
			"User": gin.H{
				"Username": username,
				"Role":     role,
			},
			"isAdmin": role == "admin",
		})
	})
}
