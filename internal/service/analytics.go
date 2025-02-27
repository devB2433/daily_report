package service

import (
	"daily-report/internal/database"
	"daily-report/internal/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AnalyticsService 统计分析服务
type AnalyticsService struct {
	DB *gorm.DB
}

// NewAnalyticsService 创建新的统计分析服务
func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		DB: database.GetDB(),
	}
}

// GetAnalyticsSummary 获取完整的分析摘要
func (s *AnalyticsService) GetAnalyticsSummary(startDate, endDate, preset string) (*model.AnalyticsSummary, error) {
	summary := &model.AnalyticsSummary{
		TimeRange: model.TimeRange{
			StartDate: startDate,
			EndDate:   endDate,
			Preset:    preset,
		},
	}

	// 1. 获取基本统计数据
	var totalUsers, totalProjects int64
	s.DB.Model(&model.User{}).Count(&totalUsers)
	s.DB.Model(&model.Project{}).Count(&totalProjects)

	summary.TotalUsers = totalUsers
	summary.TotalProjects = totalProjects

	// 2. 获取报告统计
	var reportCount int64
	s.DB.Model(&model.Report{}).Where("date BETWEEN ? AND ?", startDate, endDate).Count(&reportCount)
	summary.TotalReports = reportCount

	// 3. 计算总工时
	var totalHours float64
	s.DB.Model(&model.Task{}).
		Joins("JOIN reports ON tasks.report_id = reports.id").
		Where("reports.date BETWEEN ? AND ?", startDate, endDate).
		Select("SUM(tasks.hours)").
		Row().Scan(&totalHours)

	summary.TotalHours = totalHours

	// 4. 获取项目工时分布
	projectHours, err := s.GetProjectHoursDistribution(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("获取项目工时分布失败: %w", err)
	}
	summary.ProjectHours = projectHours

	// 5. 获取每日统计
	dailyStats, err := s.GetDailyStats(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("获取每日统计失败: %w", err)
	}
	summary.DailyStats = dailyStats

	// 6. 获取用户工时统计
	userHours, err := s.GetUserHoursStats(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("获取用户工时统计失败: %w", err)
	}
	summary.UserHours = userHours

	// 7. 获取用户提交率
	userSubmissions, err := s.GetUserSubmissionRates(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("获取用户提交率失败: %w", err)
	}
	summary.UserSubmissions = userSubmissions

	return summary, nil
}

// GetProjectHoursDistribution 获取项目工时分布
func (s *AnalyticsService) GetProjectHoursDistribution(startDate, endDate string) ([]model.ProjectHoursStat, error) {
	var results []model.ProjectHoursStat

	query := `
		SELECT 
			p.id AS project_id, 
			p.name AS project_name, 
			SUM(t.hours) AS hours
		FROM 
			tasks t
		JOIN 
			projects p ON t.project_id = p.id
		JOIN 
			reports r ON t.report_id = r.id
		WHERE 
			r.date BETWEEN ? AND ?
		GROUP BY 
			p.id, p.name
		ORDER BY 
			hours DESC
	`

	err := s.DB.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

// GetDailyStats 获取每日统计数据
func (s *AnalyticsService) GetDailyStats(startDate, endDate string) ([]model.ProjectDailyHoursStat, error) {
	// 获取日期范围内的所有日期
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var results []model.ProjectDailyHoursStat

	// 对每一天进行统计
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		dailyStat := model.ProjectDailyHoursStat{
			Date: dateStr,
		}

		// 获取当天的项目工时
		var projectHours []model.ProjectHoursStat
		projectHoursQuery := `
			SELECT 
				p.id AS project_id, 
				p.name AS project_name, 
				SUM(t.hours) AS hours
			FROM 
				tasks t
			JOIN 
				projects p ON t.project_id = p.id
			JOIN 
				reports r ON t.report_id = r.id
			WHERE 
				r.date = ?
			GROUP BY 
				p.id, p.name
		`
		s.DB.Raw(projectHoursQuery, dateStr).Scan(&projectHours)

		dailyStat.ProjectHours = projectHours
		results = append(results, dailyStat)
	}

	return results, nil
}

// GetUserHoursStats 获取用户工时统计
func (s *AnalyticsService) GetUserHoursStats(startDate, endDate string) ([]model.UserHoursStat, error) {
	var results []model.UserHoursStat

	query := `
		SELECT 
			u.username, 
			u.chinese_name, 
			SUM(t.hours) AS hours
		FROM 
			tasks t
		JOIN 
			reports r ON t.report_id = r.id
		JOIN 
			users u ON r.user_id = u.id
		WHERE 
			r.date BETWEEN ? AND ?
		GROUP BY 
			u.username, u.chinese_name
		ORDER BY 
			hours DESC
	`

	err := s.DB.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

// GetUserSubmissionRates 获取用户提交率
func (s *AnalyticsService) GetUserSubmissionRates(startDate, endDate string) ([]model.UserSubmissionStat, error) {
	var results []model.UserSubmissionStat

	// 获取所有用户
	var users []model.User
	s.DB.Find(&users)

	// 计算工作日数量
	workdays := s.calculateWorkdays(startDate, endDate)

	for _, user := range users {
		// 跳过已删除的用户
		if user.DeletedAt.Valid {
			continue
		}

		// 获取用户在该时间段内提交的报告数量
		var submittedCount int64
		s.DB.Model(&model.Report{}).
			Where("user_id = ? AND date BETWEEN ? AND ?", user.ID, startDate, endDate).
			Count(&submittedCount)

		// 计算提交率
		submissionRate := float64(submittedCount) / float64(workdays)
		if submissionRate > 1 {
			submissionRate = 1 // 最高100%
		}

		stat := model.UserSubmissionStat{
			Username:       user.Username,
			ChineseName:    user.ChineseName,
			TotalWorkdays:  workdays,
			SubmittedDays:  int(submittedCount),
			MissingDays:    workdays - int(submittedCount),
			SubmissionRate: submissionRate,
		}

		results = append(results, stat)
	}

	return results, nil
}

// calculateWorkdays 计算工作日数量（排除周末）
func (s *AnalyticsService) calculateWorkdays(startDate, endDate string) int {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	workdays := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		// 排除周六(6)和周日(0)
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			workdays++
		}
	}

	return workdays
}
