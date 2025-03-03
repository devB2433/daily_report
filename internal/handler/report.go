package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// getTimeZone 获取时区
/*func getTimeZone() *time.Location {
	// 尝试从系统获取时区
	if tz, err := time.LoadLocation("Local"); err == nil {
		return tz
	}

	// 如果无法获取系统时区，使用默认时区
	if loc, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		return loc
	}

	// 如果都失败了，返回 UTC
	return time.UTC
}
*/

// CreateReport 创建日报
func CreateReport(c *gin.Context) {
	// log.Printf("开始处理日报创建请求")

	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// log.Printf("请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	// log.Printf("成功解析请求数据: %+v", req)

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		// log.Printf("未找到用户ID")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录",
		})
		return
	}

	// 转换用户ID为uint
	var uid uint
	switch v := userID.(type) {
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			// log.Printf("用户ID转换失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "用户ID无效",
			})
			return
		}
		uid = uint(id)
	case float64:
		uid = uint(v)
	default:
		// log.Printf("用户ID类型转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID类型无效",
		})
		return
	}

	// log.Printf("当前用户ID: %v", uid)

	db := database.GetDB()

	// 解析日报时间
	reportTime, err := time.Parse("2006-01-02 15:04", req.ReportTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的日期格式",
		})
		return
	}

	// 验证日报日期是否与当前日期相同
	now := time.Now() // 直接使用服务器当前时间，不进行时区转换
	reportDate := reportTime.Format("2006-01-02")
	currentDate := now.Format("2006-01-02")
	fmt.Printf("Report Date: %s, Current Date: %s\n", reportDate, currentDate) // 添加调试日志
	if reportDate != currentDate {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "日报提交失败：只能提交当天的日报",
		})
		return
	}

	// 创建日报记录 - 修复时区问题，确保日期不会变化
	report := model.Report{
		UserID: uid,
		// 只保留日期部分，不包含时间，避免时区转换导致日期变化
		Date:        time.Date(reportTime.Year(), reportTime.Month(), reportTime.Day(), 0, 0, 0, 0, time.Local),
		SubmittedAt: time.Now(),
	}

	// 检查是否已存在当天的日报
	var existingReport model.Report
	err = db.Where("user_id = ? AND DATE(date) = DATE(?)", uid, reportTime).First(&existingReport).Error
	if err == nil {
		// log.Printf("当天已存在日报记录: ID=%d", existingReport.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "今天已经提交过日报了",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		// log.Printf("查询已存在日报时发生错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "系统错误",
		})
		return
	}

	// 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		// log.Printf("开始事务失败: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "系统错误",
		})
		return
	}

	// 确保事务在函数返回时提交或回滚
	defer func() {
		if r := recover(); r != nil {
			// log.Printf("发生panic，回滚事务: %v", r)
			tx.Rollback()
		}
	}()

	if err := tx.Create(&report).Error; err != nil {
		tx.Rollback()
		// log.Printf("创建日报记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建日报失败",
		})
		return
	}

	// log.Printf("成功创建日报记录: ID=%d", report.ID)

	// 创建工作项记录
	for _, item := range req.Items {
		// log.Printf("处理工作项: ProjectID=%d, Hours=%.1f", item.ProjectID, item.Hours)

		// 验证项目是否存在
		var project model.Project
		if err := tx.First(&project, item.ProjectID).Error; err != nil {
			tx.Rollback()
			// log.Printf("项目验证失败: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("项目不存在: ID=%d", item.ProjectID),
			})
			return
		}

		task := model.Task{
			ReportID:  report.ID,
			ProjectID: item.ProjectID,
			Content:   item.Content,
			Hours:     item.Hours,
		}

		if err := tx.Create(&task).Error; err != nil {
			tx.Rollback()
			// log.Printf("创建工作项失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "创建工作项失败",
			})
			return
		}

		// log.Printf("成功创建工作项: ID=%d", task.ID)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		// log.Printf("提交事务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存日报失败",
		})
		return
	}

	// log.Printf("事务提交成功")

	response := gin.H{
		"success": true,
		"message": "日报提交成功",
		"data": gin.H{
			"id": report.ID,
		},
	}

	// log.Printf("返回成功响应: %+v", response)
	c.JSON(http.StatusOK, response)
}

// GetReports 获取日报列表
func GetReports(c *gin.Context) {
	userID, _ := c.Get("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	db := database.GetDB()

	query := db.Preload("Tasks").Preload("Tasks.Project").Preload("User").
		Where("user_id = ?", userID)

	// 如果提供了日期范围，添加日期过滤条件
	if startDate != "" && endDate != "" {
		query = query.Where("DATE(date) >= DATE(?) AND DATE(date) <= DATE(?)", startDate, endDate)
	}

	var reports []model.Report
	if err := query.Order("date DESC").Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取日报列表失败",
		})
		return
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
	role, _ := c.Get("role")
	fmt.Printf("Getting report: id=%s, userID=%v, role=%s\n", id, userID, role)

	db := database.GetDB()

	var report model.Report
	query := db.Preload("Tasks").Preload("Tasks.Project").Preload("User")

	// 如果不是管理员，只能查看自己的日报
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
		fmt.Printf("Adding user_id condition for non-admin user\n")
	}

	if err := query.First(&report, id).Error; err != nil {
		fmt.Printf("Error finding report: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "日报不存在",
		})
		return
	}

	fmt.Printf("Found report: %+v\n", report)
	fmt.Printf("Tasks: %+v\n", report.Tasks)

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
	if err := tx.Unscoped().Where("id = ? AND user_id = ?", id, userID).First(&report).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "日报不存在或无权删除",
		})
		return
	}

	// 物理删除相关的任务
	if err := tx.Unscoped().Where("report_id = ?", id).Delete(&model.Task{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除日报任务失败",
		})
		return
	}

	// 物理删除日报
	if err := tx.Unscoped().Delete(&report).Error; err != nil {
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
		Where("reports.user_id = ? AND reports.date >= ? AND reports.date < ? AND reports.deleted_at IS NULL AND tasks.deleted_at IS NULL", userID, startDate, endDate).
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
			Where("reports.user_id = ? AND reports.date >= ? AND reports.date < ? AND reports.deleted_at IS NULL AND tasks.deleted_at IS NULL", userID, dayStart, dayEnd).
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

// 计算项目工时占比
func calculateProjectHoursPercentage(projectHours []model.ProjectHoursStat) []map[string]interface{} {
	var result []map[string]interface{}
	var totalHours float64

	// 计算总工时
	for _, ph := range projectHours {
		totalHours += ph.Hours
	}

	// 计算每个项目的工时占比
	for _, ph := range projectHours {
		percentage := (ph.Hours / totalHours) * 100
		result = append(result, map[string]interface{}{
			"name":       ph.ProjectName,
			"hours":      ph.Hours,
			"percentage": percentage,
		})
	}

	return result
}

// GetAllReports 管理员获取所有用户的日报
func GetAllReports(c *gin.Context) {
	// 验证是否为管理员
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "权限不足",
		})
		return
	}

	// 获取时间范围参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供开始和结束日期",
		})
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 50 // 默认每页50条记录

	if pageStr := c.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSizeNum, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeNum > 0 {
			pageSize = pageSizeNum
		}
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	db := database.GetDB()
	var reports []model.Report
	var total int64

	// 构建基础查询
	baseQuery := db.Model(&model.Report{}).
		Where("date >= ? AND date < DATE_ADD(?, INTERVAL 1 DAY)", startDate, endDate)

	// 获取总记录数
	baseQuery.Count(&total)

	// 查询指定时间范围内的所有日报，带分页
	query := baseQuery.
		Preload("Tasks").
		Preload("Tasks.Project").
		Preload("User").
		Order("date DESC, user_id ASC").
		Limit(pageSize).
		Offset(offset)

	if err := query.Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取日报列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"data":     reports,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
