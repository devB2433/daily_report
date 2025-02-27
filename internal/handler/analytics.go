package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"daily-report/internal/database"
	"daily-report/internal/model"
	"daily-report/internal/service"

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
	ProjectHours    []model.ProjectHoursStat      `json:"project_hours"`
	UserHours       []model.UserHoursStat         `json:"user_hours"`
	SubmissionRate  []model.SubmissionRateStat    `json:"submission_rate"`
	DailyStats      []model.ProjectDailyHoursStat `json:"daily_stats"`
	UserSubmissions []model.UserSubmissionStat    `json:"user_submissions"`
}

// GetAnalyticsSummary 获取统计分析摘要数据
func GetAnalyticsSummary(c *gin.Context) {
	// 获取时间范围参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	preset := c.Query("preset")

	fmt.Printf("收到统计请求 - 开始日期: %s, 结束日期: %s, 预设: %s\n", startDate, endDate, preset)

	// 如果没有提供时间范围，使用当月
	if startDate == "" || endDate == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		endDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location()).Format("2006-01-02")
		preset = "month"
		fmt.Printf("使用默认时间范围 - 开始日期: %s, 结束日期: %s\n", startDate, endDate)
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

	// 创建分析服务
	analyticsService := service.NewAnalyticsService()

	// 获取分析摘要
	summary, err := analyticsService.GetAnalyticsSummary(startDate, endDate, preset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("获取统计数据失败: %v", err),
		})
		return
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

// ExportReportsCSV 导出报告数据为CSV格式
func ExportReportsCSV(c *gin.Context) {
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

	db := database.GetDB()

	// 查询日报详细数据
	var detailReports []struct {
		Username    string    `json:"username"`
		ChineseName string    `json:"chinese_name"`
		Department  string    `json:"department"`
		Level       string    `json:"level"`
		Date        time.Time `json:"date"`
		ProjectCode string    `json:"project_code"`
		ProjectName string    `json:"project_name"`
		Content     string    `json:"content"`
		Hours       float64   `json:"hours"`
	}

	err = db.Table("tasks").
		Select("users.username, users.chinese_name, users.department, users.level, reports.date, projects.code as project_code, projects.name as project_name, tasks.content, tasks.hours").
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Joins("JOIN users ON reports.user_id = users.id").
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("reports.date >= ? AND reports.date <= ? AND tasks.deleted_at IS NULL AND reports.deleted_at IS NULL",
			start.Format("2006-01-02"), end.Format("2006-01-02")).
		Order("reports.date ASC, users.username ASC").
		Scan(&detailReports).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取日报详细数据失败",
		})
		return
	}

	// 设置响应头
	fileName := fmt.Sprintf("daily_reports_%s_to_%s.csv", startDate, endDate)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// 写入UTF-8 BOM
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// 创建CSV写入器
	writer := csv.NewWriter(c.Writer)

	// 写入日报详细数据表头
	writer.Write([]string{"用户名", "姓名", "部门", "级别", "周数", "日期", "项目编号", "项目名称", "工作内容", "工时"})

	// 写入日报详细数据
	for _, report := range detailReports {
		// 获取周数
		_, week := report.Date.ISOWeek()
		weekStr := fmt.Sprintf("WEEK %d", week)

		writer.Write([]string{
			report.Username,
			report.ChineseName,
			report.Department,
			report.Level,
			weekStr,
			report.Date.Format("2006-01-02"),
			report.ProjectCode,
			report.ProjectName,
			report.Content,
			fmt.Sprintf("%.1f", report.Hours),
		})
	}

	writer.Flush()
}

// 计算两个日期之间的工作日数量（不包括周末）
func calculateWorkdays(start, end time.Time) int {
	// 如果开始日期在结束日期之后，返回0
	if start.After(end) {
		return 0
	}

	// 如果是同一天，检查是否是工作日
	if start.Equal(end) {
		if start.Weekday() != time.Saturday && start.Weekday() != time.Sunday {
			return 1
		}
		return 0
	}

	workdays := 0
	current := start

	for !current.After(end) {
		// 周六是6，周日是0
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			workdays++
		}
		current = current.AddDate(0, 0, 1)
	}

	return workdays
}
