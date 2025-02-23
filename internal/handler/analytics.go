package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
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
	if err := db.Model(&model.Report{}).Unscoped().
		Where("date >= ? AND date < DATE_ADD(?, INTERVAL 1 DAY) AND deleted_at IS NULL",
			start.Format("2006-01-02"),
			end.Format("2006-01-02")).
		Count(&summary.TotalReports).Error; err != nil {
		fmt.Printf("获取日报数失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取日报统计失败",
		})
		return
	}
	fmt.Printf("找到 %d 条日报记录\n", summary.TotalReports)

	// 4. 获取指定时间范围内的总工时
	var totalHours float64
	query := db.Model(&model.Task{}).
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Where("reports.date >= ? AND reports.date < DATE_ADD(?, INTERVAL 1 DAY) AND tasks.deleted_at IS NULL AND reports.deleted_at IS NULL",
			start.Format("2006-01-02"),
			end.Format("2006-01-02")).
		Select("COALESCE(SUM(tasks.hours), 0)")

	// 打印SQL查询语句
	sql := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("reports.date >= ? AND reports.date < DATE_ADD(?, INTERVAL 1 DAY)", start.Format("2006-01-02"), end.Format("2006-01-02"))
	})
	fmt.Printf("总工时查询SQL: %s\n", sql)

	if err := query.Scan(&totalHours).Error; err != nil {
		fmt.Printf("获取总工时失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取工时统计失败",
		})
		return
	}
	fmt.Printf("总工时: %.1f 小时\n", totalHours)
	summary.TotalHours = totalHours

	// 5. 获取项目工时分布
	var projectHours []model.ProjectHoursStat
	if err := db.Table("tasks").
		Select("projects.id as project_id, projects.name as project_name, COALESCE(SUM(tasks.hours), 0) as hours").
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("reports.date >= ? AND reports.date < DATE_ADD(?, INTERVAL 1 DAY) AND tasks.deleted_at IS NULL AND reports.deleted_at IS NULL",
			start.Format("2006-01-02"),
			end.Format("2006-01-02")).
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
		Select("users.username, users.chinese_name, COALESCE(SUM(tasks.hours), 0) as hours").
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Joins("JOIN users ON reports.user_id = users.id").
		Where("reports.date >= ? AND reports.date < DATE_ADD(?, INTERVAL 1 DAY)",
			start.Format("2006-01-02"),
			end.Format("2006-01-02")).
		Group("users.username, users.chinese_name").
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

		var dayStats model.ProjectDailyHoursStat
		dayStats.Date = currentDate.Format("2006-01-02")

		// 获取当天的项目工时
		var dayProjectHours []model.ProjectHoursStat
		query := db.Table("tasks").
			Select("projects.id as project_id, projects.name as project_name, COALESCE(SUM(tasks.hours), 0) as hours").
			Joins("JOIN reports ON tasks.report_id = reports.id").
			Joins("JOIN projects ON tasks.project_id = projects.id").
			Where("reports.date >= ? AND reports.date < DATE_ADD(?, INTERVAL 1 DAY) AND tasks.deleted_at IS NULL AND reports.deleted_at IS NULL",
				dayStart.Format("2006-01-02"),
				dayStart.Format("2006-01-02")).
			Group("projects.id, projects.name").
			Having("hours > 0").
			Order("hours DESC")

		// 打印SQL查询语句
		sql := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Where("reports.date >= ? AND reports.date < DATE_ADD(?, INTERVAL 1 DAY)", dayStart.Format("2006-01-02"), dayStart.Format("2006-01-02"))
		})
		fmt.Printf("查询日期 %s 的SQL: %s\n", dayStart.Format("2006-01-02"), sql)

		if err := query.Scan(&dayProjectHours).Error; err == nil {
			fmt.Printf("日期 %s 找到 %d 个项目工时记录\n", dayStart.Format("2006-01-02"), len(dayProjectHours))
			for _, ph := range dayProjectHours {
				fmt.Printf("  - 项目: %s, 工时: %.1f\n", ph.ProjectName, ph.Hours)
			}
			dayStats.ProjectHours = dayProjectHours
			summary.DailyStats = append(summary.DailyStats, dayStats)
		} else {
			fmt.Printf("日期 %s 查询出错: %v\n", dayStart.Format("2006-01-02"), err)
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

	// 9. 获取用户提交率统计
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户列表失败",
		})
		return
	}

	// 计算工作日数量
	workdays := calculateWorkdays(start, end)

	// 获取每个用户在时间范围内的提交记录
	for _, user := range users {
		var submittedCount int64
		if err := db.Model(&model.Report{}).
			Where("user_id = ? AND date >= ? AND date < DATE_ADD(?, INTERVAL 1 DAY) AND deleted_at IS NULL",
				user.ID, start.Format("2006-01-02"), end.Format("2006-01-02")).
			Count(&submittedCount).Error; err != nil {
			continue
		}

		stat := model.UserSubmissionStat{
			Username:      user.Username,
			ChineseName:   user.ChineseName,
			TotalWorkdays: workdays,
			SubmittedDays: int(submittedCount),
			MissingDays:   workdays - int(submittedCount),
		}

		// 处理边缘情况：如果没有工作日，设置提交率为1.0（100%）
		if workdays == 0 {
			stat.SubmissionRate = 1.0
		} else {
			stat.SubmissionRate = float64(submittedCount) / float64(workdays)
		}

		summary.UserSubmissions = append(summary.UserSubmissions, stat)
	}

	// 按提交率降序排序
	sort.Slice(summary.UserSubmissions, func(i, j int) bool {
		return summary.UserSubmissions[i].SubmissionRate > summary.UserSubmissions[j].SubmissionRate
	})

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
		Date        time.Time `json:"date"`
		ProjectName string    `json:"project_name"`
		Content     string    `json:"content"`
		Hours       float64   `json:"hours"`
	}

	err = db.Table("tasks").
		Select("users.username, users.chinese_name, users.department, reports.date, projects.name as project_name, tasks.content, tasks.hours").
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
	writer.Write([]string{"用户名", "姓名", "部门", "日期", "项目名称", "工作内容", "工时"})

	// 写入日报详细数据
	for _, report := range detailReports {
		writer.Write([]string{
			report.Username,
			report.ChineseName,
			report.Department,
			report.Date.Format("2006-01-02"),
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
