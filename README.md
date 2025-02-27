# 工作日报系统

一个简单的工作日报管理系统，支持项目管理、工时统计和数据分析。使用Cursor编写。

## 功能特性

- 用户管理：支持用户注册、登录和权限控制
- 项目管理：创建和管理项目，跟踪项目状态
- 日报管理：便捷的日报提交和查看功能
- 数据分析：工时统计和项目进展可视化
- 响应式设计：支持多种设备访问

## 系统要求

- Docker 20.10.0 或更高版本
- Docker Compose 2.0.0 或更高版本
- 至少 2GB 可用内存
- 至少 1GB 可用磁盘空间

## 快速开始

### 开发环境

1. 克隆代码仓库：
   ```bash
   git clone <repository-url>
   cd daily-report
   ```

2. 使用开发配置启动系统：
   ```bash
   docker-compose -f docker-compose.dev.yml up --build
   ```

3. 导入测试数据（可选）：
   ```bash
   # 连接到 MySQL 容器
   docker exec -it daily_report-mysql-1 mysql -udaily_report -pdaily_report_password daily_report

   # 在 MySQL 命令行中执行
   source /docker-entrypoint-initdb.d/mock_users.sql
   source /docker-entrypoint-initdb.d/mock_projects.sql
   source /docker-entrypoint-initdb.d/mock_reports.sql
   ```

### 生产环境

1. 构建并启动系统：
   ```bash
   # 强制重新构建所有镜像
   docker-compose -f docker-compose.yml build --no-cache
   
   # 启动系统
   docker-compose -f docker-compose.yml up -d
   ```

## 默认账户

系统首次启动时会自动创建管理员账户：
- 邮箱：admin@blingsec.cn
- 密码：Kj#9mP$2nL5vB@8x

建议首次登录后立即修改密码。

## 开发与生产环境区别

### 开发环境 (docker-compose.dev.yml)
- 启用调试模式
- 代码热重载
- 详细日志输出
- 暴露数据库端口便于调试
- 挂载源代码目录实现实时更新

### 生产环境 (docker-compose.yml)
- 发布模式运行
- 优化的性能配置
- 最小化日志输出
- 仅暴露必要端口
- 资源使用限制
- 容器健康检查
- 自动服务恢复

## 测试数据说明

系统提供了三个测试数据脚本：

1. `mock_users.sql`: 创建测试用户账号
   - 包含不同角色的用户
   - 预设的用户名和密码

2. `mock_projects.sql`: 创建示例项目
   - 包含不同状态的项目
   - 模拟真实项目数据

3. `mock_reports.sql`: 生成历史日报数据
   - 包含多个时间段的日报
   - 不同项目的工时分配

## 系统维护

### 数据备份
```bash
# 备份数据库
docker exec daily_report-mysql-1 mysqldump -u daily_report -p daily_report > backup.sql

# 恢复数据库
docker exec -i daily_report-mysql-1 mysql -u daily_report -p daily_report < backup.sql
```

### 日志查看
```bash
# 查看应用日志
docker-compose logs -f app

# 查看数据库日志
docker-compose logs -f mysql
```


### 系统更新
```bash
# 放弃本地修改（之前本地修改过Dockerfile的代理设置，默认是注释的）
git fetch origin main
git reset --hard origin/main

# 拉取最新代码
git pull origin main

#检查环境时间，确保是北京时间正确无误
date


#如果时间不对，先同步
ntpdate ntp.aliyun.com

#最好加到crond里定时获取

#修改Dockerfile，取消这一行注释，让服务器优先使用国内的代理地址，加快编译速度
ENV GOPROXY=https://goproxy.cn,direct

# 重新构建并启动
docker-compose -f docker-compose.yml down
docker-compose -f docker-compose.yml build --no-cache
docker-compose -f docker-compose.yml up -d
```

## 关于如何使用
1、管理员登陆后，可以创建项目。
2、普通成员首先应注册账号，成员登陆后，可以提交日报，查看日报。
3、管理员可以查看所有成员的日报，也可以查看统计。


