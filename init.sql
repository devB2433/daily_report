-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建默认管理员用户
-- 密码为: Admin@123 (已经过bcrypt加密)
INSERT INTO users (username, email, password, role, created_at, updated_at) VALUES 
('admin', 'admin@blingsec.cn', '$2a$10$QOJ1Z6PFr5Qx5pZnJ5vC9.wVJJwlh5FFY5UOvTxNvZ7wxQwPUm3Hy', 'admin', NOW(), NOW());

-- 创建示例项目
INSERT INTO projects (name, code, description, manager, client, status, created_at, updated_at) VALUES
('示例项目1', 'DEMO-001', '这是一个示例项目', '项目经理', '客户A', 'active', NOW(), NOW()),
('示例项目2', 'DEMO-002', '这是另一个示例项目', '项目经理', '客户B', 'active', NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1; 