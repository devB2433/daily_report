package handler

import (
	"fmt"
	"net/http"
	"time"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AnalyticsSummary 统计分析摘要
type AnalyticsSummary struct {
	TotalUsers    int64   `json:"total_users"`
	TotalProjects int64   `json:"total_projects"`
	TotalReports  int64   `json:"monthly_reports"`
	TotalHours    float64 `json:"monthly_hours"`
	TimeRange     struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Preset    string `json:"preset"`
	} `json:"time_range"`
	ProjectHours   []model.ProjectHoursStat      `json:"project_hours"`
	UserHours      []model.UserHoursStat         `json:"user_hours"`
	SubmissionRate []model.SubmissionRateStat    `json:"submission_rate"`
	DailyStats     []model.ProjectDailyHoursStat `json:"daily_stats"`
}

// GetAnalyticsSummary 获取统计分析摘要数据
func GetAnalyticsSummary(c *gin.Context) {
	// 获取时间范围参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	preset := c.Query("preset")

	// 如果没有提供时间范围，使用当月
	if startDate == "" || endDate == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		endDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location()).Format("2006-01-02")
		preset = "month"
	}

	// 解析时间
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "开始日期格式无效",
		})
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "结束日期格式无效",
		})
		return
	}

	// 验证时间范围
	if end.Before(start) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "结束日期不能早于开始日期",
		})
		return
	}

	// 限制查询范围不超过1年
	if end.Sub(start) > 365*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "查询时间范围不能超过1年",
		})
		return
	}

	db := database.GetDB()
	var summary AnalyticsSummary

	// 设置时间范围
	summary.TimeRange.StartDate = startDate
	summary.TimeRange.EndDate = endDate
	summary.TimeRange.Preset = preset

	// 1. 获取总用户数
	if err := db.Model(&model.User{}).Count(&summary.TotalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户统计失败",
		})
		return
	}

	// 2. 获取总项目数（只统计活动项目）
	if err := db.Model(&model.Project{}).Where("status = ?", "active").Count(&summary.TotalProjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取项目统计失败",
		})
		return
	}

	// 3. 获取指定时间范围内的日报数
	if err := db.Model(&model.Report{}).Where("date >= ? AND date <= ?", start, end).Count(&summary.TotalReports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取日报统计失败",
		})
		return
	}

	// 4. 获取指定时间范围内的总工时
	var totalHours float64
	if err := db.Model(&model.Task{}).
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Where("reports.date >= ? AND reports.date <= ?", start, end).
		Select("COALESCE(SUM(tasks.hours), 0)").
		Scan(&totalHours).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取工时统计失败",
		})
		return
	}
	summary.TotalHours = totalHours

	// 5. 获取项目工时分布
	var projectHours []model.ProjectHoursStat
	if err := db.Table("tasks").
		Select("projects.id as project_id, projects.name as project_name, COALESCE(SUM(tasks.hours), 0) as hours").
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("reports.date >= ? AND reports.date <= ?", start, end).
		Group("projects.id, projects.name").
		Having("hours > 0").
		Order("hours DESC").
		Scan(&projectHours).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取项目工时统计失败",
		})
		return
	}
	summary.ProjectHours = projectHours

	// 6. 获取用户工时统计
	if err := db.Table("tasks").
		Select("users.username, COALESCE(SUM(tasks.hours), 0) as hours").
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Joins("JOIN users ON reports.user_id = users.id").
		Where("reports.date >= ? AND reports.date <= ?", start, end).
		Group("users.username").
		Having("hours > 0").
		Order("hours DESC").
		Scan(&summary.UserHours).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户工时统计失败",
		})
		return
	}

	// 7. 获取每日项目工时统计
	currentDate := start
	for !currentDate.After(end) {
		dayStart := currentDate
		dayEnd := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 59, 59, 999999999, currentDate.Location())

		var dayStats model.ProjectDailyHoursStat
		dayStats.Date = currentDate.Format("2006-01-02")

		// 获取当天的项目工时
		var dayProjectHours []model.ProjectHoursStat
		if err := db.Table("tasks").
			Select("projects.id as project_id, projects.name as project_name, COALESCE(SUM(tasks.hours), 0) as hours").
			Joins("JOIN reports ON tasks.report_id = reports.id").
			Joins("JOIN projects ON tasks.project_id = projects.id").
			Where("reports.date >= ? AND reports.date <= ?", dayStart, dayEnd).
			Group("projects.id, projects.name").
			Having("hours > 0").
			Order("hours DESC").
			Scan(&dayProjectHours).Error; err == nil {
			dayStats.ProjectHours = dayProjectHours
			summary.DailyStats = append(summary.DailyStats, dayStats)
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// 8. 获取日报提交率趋势
	var submissionStats []struct {
		Date      string
		Submitted int64
		Total     int64
	}

	if err := db.Raw(`
        WITH RECURSIVE dates AS (
            SELECT DATE(?) as date
            UNION ALL
            SELECT DATE_ADD(date, INTERVAL 1 DAY)
            FROM dates
            WHERE date < DATE(?)
        )
        SELECT 
            dates.date,
            COUNT(DISTINCT CASE WHEN reports.id IS NOT NULL THEN reports.user_id END) as submitted,
            COUNT(DISTINCT users.id) as total
        FROM dates
        CROSS JOIN users
        LEFT JOIN reports ON DATE(reports.date) = dates.date AND reports.user_id = users.id
        WHERE users.created_at <= dates.date
        GROUP BY dates.date
        ORDER BY dates.date
    `, startDate, endDate).
		Scan(&submissionStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取提交率统计失败",
		})
		return
	}

	// 计算提交率
	for _, stat := range submissionStats {
		if stat.Total > 0 {
			summary.SubmissionRate = append(summary.SubmissionRate, model.SubmissionRateStat{
				Date: stat.Date,
				Rate: float64(stat.Submitted) / float64(stat.Total),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// 获取项目工时统计
func getProjectHoursStats(db *gorm.DB) ([]model.ProjectHoursStat, error) {
	var stats []model.ProjectHoursStat

	// 使用原生SQL查询来获取项目工时统计
	query := `
        SELECT 
            p.id,
            p.name,
            COALESCE(SUM(t.hours), 0) as hours
        FROM 
            projects p
            LEFT JOIN tasks t ON p.id = t.project_id
            LEFT JOIN reports r ON t.report_id = r.id
        WHERE
            r.date >= DATE_FORMAT(NOW(), '%Y-%m-01')
            AND r.date < DATE_ADD(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL 1 MONTH)
        GROUP BY 
            p.id, p.name
        HAVING 
            hours > 0
        ORDER BY 
            hours DESC
    `

	err := db.Raw(query).Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("查询项目工时统计失败: %v", err)
	}

	return stats, nil
}
