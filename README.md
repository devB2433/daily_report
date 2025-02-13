# 日报管理系统

这是一个基于 Go 语言开发的日报管理系统，使用 Gin 框架和 MySQL 数据库。

## 功能特性

- 用户管理：注册、登录、权限控制
- 项目管理：创建、编辑、删除项目
- 日报管理：提交日报、查看历史日报
- 统计分析：工时统计、项目分析、提交率分析

## 技术栈

- 后端：Go + Gin + GORM
- 前端：Bootstrap + Chart.js
- 数据库：MySQL 8.0
- 容器化：Docker + Docker Compose

## 快速开始

### 使用 Docker Compose 部署

1. 确保已安装 Docker 和 Docker Compose

2. 克隆项目
```bash
git clone [项目地址]
cd daily-report
```

3. 启动服务
```bash
docker-compose up -d
```

4. 访问系统
- 地址：http://localhost:8080
- 默认管理员账号：
  - 邮箱：admin@blingsec.cn
  - 密码：Kj#9mP$2nL5vB@8x

### 手动部署

1. 安装 Go 1.21 或更高版本

2. 安装并配置 MySQL 8.0

3. 克隆项目并安装依赖
```bash
git clone [项目地址]
cd daily-report
go mod download
```

4. 配置数据库连接
- 修改 config/config.yaml 文件中的数据库配置

5. 运行项目
```bash
go run cmd/main.go
```

## 配置说明

### 环境变量

- `DB_HOST`: 数据库主机地址
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `DB_CHARSET`: 数据库字符集
- `DB_LOC`: 数据库时区

## 注意事项

1. 首次使用请修改默认管理员密码
2. 建议在生产环境中修改数据库密码
3. 生产环境部署时建议启用 HTTPS
4. 定期备份数据库数据

## 许可证

MIT License 