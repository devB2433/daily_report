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
	reportTime, err := time.Parse("2006-01-02 15:04", req.ReportTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "日报时间格式无效",
			"error":   err.Error(),
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

	// 检查当天是否已经提交过日报
	startOfDay := time.Date(reportTime.Year(), reportTime.Month(), reportTime.Day(), 0, 0, 0, 0, reportTime.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var existingReport model.Report
	if err := db.Where("user_id = ? AND date >= ? AND date < ?", userID, startOfDay, endOfDay).First(&existingReport).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("当天已经提交过日报（ID: %d），你需要先删除旧的日报，才能重新提交", existingReport.ID),
		})
		return
	}

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

// DeleteReport 删除日报
func DeleteReport(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	// 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除日报失败：无法开始事务",
		})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 先检查日报是否存在且属于当前用户
	var report model.Report
	if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&report).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "日报不存在或无权删除",
		})
		return
	}

	// 删除相关的任务
	if err := tx.Where("report_id = ?", id).Delete(&model.Task{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除日报任务失败",
		})
		return
	}

	// 删除日报
	if err := tx.Delete(&report).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除日报失败",
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除日报失败：提交事务失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "日报删除成功",
	})
}

// ReportSubmissionStatus 表示日报提交状态
type ReportSubmissionStatus struct {
	Date      string `json:"date"`
	Submitted bool   `json:"submitted"`
}

// GetReportSubmissionStatus 获取日报提交状态
func GetReportSubmissionStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	year := c.Query("year")
	month := c.Query("month")

	if year == "" || month == "" {
		now := time.Now()
		year = fmt.Sprintf("%d", now.Year())
		month = fmt.Sprintf("%02d", now.Month())
	}

	// 解析年月
	yearNum, _ := time.Parse("2006", year)
	monthNum, _ := time.Parse("01", month)

	// 计算月份的开始和结束时间
	startDate := time.Date(yearNum.Year(), monthNum.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	db := database.GetDB()

	// 获取该月所有的日报记录
	var reports []model.Report
	if err := db.Where("user_id = ? AND date >= ? AND date < ?", userID, startDate, endDate).
		Select("date").
		Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取日报记录失败",
		})
		return
	}

	// 创建提交记录映射
	submittedDates := make(map[string]bool)
	for _, report := range reports {
		dateStr := report.Date.Format("2006-01-02")
		submittedDates[dateStr] = true
	}

	// 生成当月所有工作日的状态
	var result []ReportSubmissionStatus
	currentDate := startDate
	for currentDate.Before(endDate) {
		// 跳过周末
		if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
			dateStr := currentDate.Format("2006-01-02")
			result = append(result, ReportSubmissionStatus{
				Date:      dateStr,
				Submitted: submittedDates[dateStr],
			})
		} else {
			// 对于周末，只在提交了日报的情况下添加记录
			dateStr := currentDate.Format("2006-01-02")
			if submittedDates[dateStr] {
				result = append(result, ReportSubmissionStatus{
					Date:      dateStr,
					Submitted: true,
				})
			}
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ProjectHoursStat 项目工时统计
type ProjectHoursStat struct {
	ProjectID   uint    `json:"project_id"`
	ProjectName string  `json:"project_name"`
	TotalHours  float64 `json:"total_hours"`
}

// DailyHoursStat 每日工时统计
type DailyHoursStat struct {
	Date         string             `json:"date"`
	TotalHours   float64            `json:"total_hours"`
	ProjectHours []ProjectHoursStat `json:"project_hours"`
}

// GetMonthlyStats 获取月度统计数据
func GetMonthlyStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	year := c.Query("year")
	month := c.Query("month")

	if year == "" || month == "" {
		now := time.Now()
		year = fmt.Sprintf("%d", now.Year())
		month = fmt.Sprintf("%02d", now.Month())
	}

	// 解析年月
	yearNum, _ := time.Parse("2006", year)
	monthNum, _ := time.Parse("01", month)

	// 计算月份的开始和结束时间
	startDate := time.Date(yearNum.Year(), monthNum.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	db := database.GetDB()

	// 1. 获取当月项目工时统计
	var projectStats []ProjectHoursStat
	err := db.Table("tasks").
		Select("tasks.project_id, projects.name as project_name, SUM(tasks.hours) as total_hours").
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("reports.user_id = ? AND reports.date >= ? AND reports.date < ?", userID, startDate, endDate).
		Group("tasks.project_id").
		Order("total_hours DESC").
		Scan(&projectStats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取项目统计失败",
		})
		return
	}

	// 2. 获取每日工时统计
	var dailyStats []DailyHoursStat
	currentDate := startDate
	for currentDate.Before(endDate) {
		dayStart := currentDate
		dayEnd := currentDate.Add(24 * time.Hour)

		var dayStats DailyHoursStat
		dayStats.Date = currentDate.Format("2006-01-02")

		// 获取当天的项目工时
		var projectHours []ProjectHoursStat
		err := db.Table("tasks").
			Select("tasks.project_id, projects.name as project_name, SUM(tasks.hours) as total_hours").
			Joins("JOIN reports ON tasks.report_id = reports.id").
			Joins("JOIN projects ON tasks.project_id = projects.id").
			Where("reports.user_id = ? AND reports.date >= ? AND reports.date < ?", userID, dayStart, dayEnd).
			Group("tasks.project_id").
			Scan(&projectHours).Error

		if err == nil {
			dayStats.ProjectHours = projectHours
			// 计算当天总工时
			for _, ph := range projectHours {
				dayStats.TotalHours += ph.TotalHours
			}
		}

		dailyStats = append(dailyStats, dayStats)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"project_stats": projectStats,
			"daily_stats":   dailyStats,
		},
	})
}
