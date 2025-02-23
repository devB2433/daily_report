package model

// ProjectHoursStat 项目工时统计
type ProjectHoursStat struct {
	ProjectID   uint    `json:"project_id"`
	ProjectName string  `json:"project_name"`
	Hours       float64 `json:"hours"`
}

// UserHoursStat 用户工时统计
type UserHoursStat struct {
	Username    string  `json:"username"`
	ChineseName string  `json:"chinese_name"`
	Hours       float64 `json:"hours"`
}

// SubmissionRateStat 提交率统计
type SubmissionRateStat struct {
	Date string  `json:"date"`
	Rate float64 `json:"rate"`
}

// ProjectDailyHoursStat 项目每日工时统计
type ProjectDailyHoursStat struct {
	Date         string             `json:"date"`
	ProjectHours []ProjectHoursStat `json:"project_hours"`
}

// UserSubmissionStat 用户提交率统计
type UserSubmissionStat struct {
	Username       string  `json:"username"`
	ChineseName    string  `json:"chinese_name"`
	TotalWorkdays  int     `json:"total_workdays"`  // 应提交次数（工作日）
	SubmittedDays  int     `json:"submitted_days"`  // 已提交次数
	MissingDays    int     `json:"missing_days"`    // 未提交次数
	SubmissionRate float64 `json:"submission_rate"` // 提交率
}
