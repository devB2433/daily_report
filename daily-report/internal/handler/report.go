package handler

import (
	"fmt"
	"net/http"
	"time"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
)

// ReportItem 表示一个工作项
type ReportItem struct {
	ProjectID uint    `json:"project_id" binding:"required"`
	Hours     float64 `json:"hours" binding:"required,min=0.5,max=24"`
	Content   string  `json:"content" binding:"required,min=1"`
}

// CreateReportRequest 创建日报的请求结构
type CreateReportRequest struct {
	ReportTime string       `json:"report_time" binding:"required"`
	Items      []ReportItem `json:"items" binding:"required,min=1,dive"`
}

// CreateReport 创建日报
func CreateReport(c *gin.Context) {
	// 1. 解析请求数据
	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
			"error":   err.Error(),
		})
		return
	}

	// 2. 解析日报时间
	reportTime, err := time.Parse("2006-01-02-15-04", req.ReportTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "日报时间格式无效",
		})
		return
	}

	// 3. 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	db := database.GetDB()

	// 4. 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建日报失败：无法开始事务",
		})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 5. 创建日报记录
	report := model.Report{
		UserID:      userID.(uint),
		Date:        reportTime,
		Status:      "submitted",
		SubmittedAt: time.Now(),
	}

	if err := tx.Create(&report).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("创建日报失败：%v", err),
		})
		return
	}

	// 6. 创建工作项
	for _, item := range req.Items {
		// 验证项目是否存在且处于活动状态
		var project model.Project
		if err := tx.Where("id = ? AND status = ?", item.ProjectID, "active").First(&project).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("项目不存在或已关闭：%d", item.ProjectID),
			})
			return
		}

		task := model.Task{
			ReportID:  report.ID,
			ProjectID: item.ProjectID,
			Hours:     item.Hours,
			Content:   item.Content,
			Status:    "completed",
		}

		if err := tx.Create(&task).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("创建工作项失败：%v", err),
			})
			return
		}
	}

	// 7. 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("提交日报失败：%v", err),
		})
		return
	}

	// 8. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "日报提交成功",
		"data": gin.H{
			"id": report.ID,
		},
	})
}

// GetReports 获取日报列表
func GetReports(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var reports []model.Report
	if err := db.Preload("Tasks").
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取日报列表失败",
		})
		return
	}

	// 手动加载每个任务的项目信息
	for i := range reports {
		for j := range reports[i].Tasks {
			var project model.Project
			if err := db.First(&project, reports[i].Tasks[j].ProjectID).Error; err == nil {
				reports[i].Tasks[j].Project = &project
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    reports,
	})
}

// GetReport 获取日报详情
func GetReport(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var report model.Report
	if err := db.Preload("Tasks").Preload("Tasks.Project").
		Where("id = ? AND user_id = ?", id, userID).
		First(&report).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "日报不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
	})
}
