package handler

import (
	"fmt"
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

	// 查询所有项目，按状态和名称排序
	if err := db.Order("FIELD(status, 'active', 'suspended', 'completed'), name ASC").Find(&projects).Error; err != nil {
		// log.Printf("Failed to get projects: %v", err)
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
			"message": "暂无项目",
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
	// log.Printf("开始处理创建项目请求")

	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		// log.Printf("解析项目数据失败: %v, 请求体: %+v", err, c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	// log.Printf("接收到的项目数据: %+v", project)

	// 验证必填字段
	if project.Name == "" || project.Code == "" {
		// log.Printf("必填字段验证失败 - name: %s, code: %s", project.Name, project.Code)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目名称和代号不能为空",
		})
		return
	}

	// 验证项目状态
	validStatus := map[string]bool{
		"active":    true,
		"completed": true,
		"suspended": true,
	}
	if !validStatus[project.Status] {
		// log.Printf("项目状态无效: %s", project.Status)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目状态无效",
		})
		return
	}

	// log.Printf("开始数据库操作")
	db := database.GetDB()

	// 检查项目名称是否已存在
	var existingProject model.Project
	if err := db.Where("name = ?", project.Name).First(&existingProject).Error; err == nil {
		// log.Printf("项目名称已存在: %s", project.Name)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目名称已存在",
		})
		return
	}

	// 检查项目代号是否已存在
	if err := db.Where("code = ?", project.Code).First(&existingProject).Error; err == nil {
		// log.Printf("项目代号已存在: %s", project.Code)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目代号已存在",
		})
		return
	}

	// 创建项目
	if err := db.Create(&project).Error; err != nil {
		// log.Printf("创建项目失败: %+v, 错误: %v", project, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建项目失败",
		})
		return
	}

	// log.Printf("项目创建成功: ID=%d, Name=%s", project.ID, project.Name)
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

	// 检查是否存在与该项目相关联的日报
	var count int64
	db.Model(&model.Task{}).Where("project_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无法删除项目，因为存在与该项目相关联的日报",
		})
		return
	}

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
	log.Printf("获取项目详情请求，项目ID: %s", id)

	db := database.GetDB()

	var project model.Project
	if err := db.First(&project, id).Error; err != nil {
		log.Printf("查找项目失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("项目不存在: %v", err),
		})
		return
	}

	log.Printf("成功获取项目详情: ID=%d, Name=%s", project.ID, project.Name)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
	})
}
