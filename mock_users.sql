-- Add mock user data (using bcrypt hashed password: User@123)
INSERT INTO users (username, email, password, password_hash, role, created_at, updated_at) VALUES 
('zhangsan', 'zhangsan@blingsec.cn', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', 'user', NOW(), NOW()),

('lisi', 'lisi@blingsec.cn', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', 'user', NOW(), NOW()),

('wangwu', 'wangwu@blingsec.cn', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', 'user', NOW(), NOW()),

('zhaoliu', 'zhaoliu@blingsec.cn', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', 'user', NOW(), NOW()),

('qianqi', 'qianqi@blingsec.cn', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', 'user', NOW(), NOW()),

('sunba', 'sunba@blingsec.cn', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', 'user', NOW(), NOW()),

('zhoujiu', 'zhoujiu@blingsec.cn', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', 'user', NOW(), NOW()),

('wushi', 'wushi@blingsec.cn', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', '$2a$10$dxg8uVUp58TNeqkcFwAd.OFsuYEGlt2pkZOzgiXVTWOy4LDcsmZBa', 'user', NOW(), NOW()); 