# 工作日报系统

这是一个简单的工作日报管理系统，用于记录和管理日常工作报告。

## 功能特性

- 用户认证和授权
- 日报的创建、编辑、查看和管理
- 数据统计和分析
- 友好的Web界面

## 技术栈

- Go 1.20+
- Gin Web Framework
- GORM
- SQLite3
- Bootstrap 5

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/yourusername/daily-report.git
```

2. 安装依赖
```bash
go mod download
```

3. 运行项目
```bash
go run cmd/main.go
```

4. 访问系统
打开浏览器访问 http://localhost:8080

## 项目结构

```
daily-report/
├── cmd/                    # 主程序入口
├── internal/              # 内部包
│   ├── handler/          # HTTP处理器
│   ├── model/            # 数据模型
│   └── service/          # 业务逻辑
└── web/                  # Web资源
    ├── static/           # 静态文件
    └── templates/        # HTML模板
```

## 配置说明

系统配置文件位于 `config/config.yaml`，包含数据库、服务器等配置项。

## 许可证

MIT License 