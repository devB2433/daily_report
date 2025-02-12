package handler

import (
	"log"
	"net/http"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
)

// RootRequired 检查是否为管理员用户的中间件
func RootRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未登录",
			})
			c.Abort()
			return
		}

		if role.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "权限不足，需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetProjects 获取项目列表
func GetProjects(c *gin.Context) {
	db := database.GetDB()
	var projects []model.Project

	// 查询所有进行中的项目，按名称排序
	if err := db.Where("status = ?", "active").Order("name ASC").Find(&projects).Error; err != nil {
		log.Printf("Failed to get projects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取项目列表失败",
		})
		return
	}

	// 如果没有找到任何项目
	if len(projects) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []model.Project{},
			"message": "暂无进行中的项目",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
	})
}

// CreateProject 创建新项目
func CreateProject(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	// 验证项目状态
	validStatuses := map[string]bool{
		"active":    true,
		"completed": true,
		"suspended": true,
	}
	if !validStatuses[project.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的项目状态，必须是 active、completed 或 suspended",
		})
		return
	}

	db := database.GetDB()
	if err := db.Create(&project).Error; err != nil {
		log.Printf("创建项目失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建项目失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建项目成功",
		"data":    project,
	})
}

// UpdateProject 更新项目信息
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project model.Project

	db := database.GetDB()
	if err := db.First(&project, id).Error; err != nil {
		log.Printf("Failed to find project: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "项目不存在",
		})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	if err := db.Save(&project).Error; err != nil {
		log.Printf("Failed to update project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新项目失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新项目成功",
		"data":    project,
	})
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	if err := db.Delete(&model.Project{}, id).Error; err != nil {
		log.Printf("Failed to delete project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除项目失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除项目成功",
	})
}

// GetProject 获取项目详情
func GetProject(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var project model.Project
	if err := db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "项目不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
	})
}
