-- Add mock project data
INSERT INTO projects (name, code, description, manager, client, status, start_date, end_date, created_at, updated_at) VALUES 
('Web Security Scanner System', 'WSS-2024', 'Web application security scanner based on OWASP Top 10, supporting automated vulnerability detection and report generation', 'Zhang San', 'Tech Corp', 'active', '2024-01-01', '2024-06-30', NOW(), NOW()),

('Mobile App Security Testing Platform', 'MAST-2024', 'Security testing platform for Android and iOS applications, including static analysis and dynamic testing', 'Li Si', 'Finance Group', 'active', '2024-02-01', '2024-08-31', NOW(), NOW()),

('Security Operations Center', 'SOC-2024', 'Enterprise-level security operations center solution, including security device deployment, log analysis, and incident response', 'Wang Wu', 'Government Dept', 'active', '2024-03-01', '2024-12-31', NOW(), NOW()),

('Code Security Audit System', 'CSA-2023', 'Automated source code security audit system supporting multiple programming languages and frameworks', 'Zhao Liu', 'Bank Corp', 'completed', '2023-06-01', '2023-12-31', NOW(), NOW()),

('Vulnerability Management Platform', 'VMP-2023', 'Enterprise vulnerability management platform supporting vulnerability lifecycle management and remediation tracking', 'Qian Qi', 'Internet Corp', 'completed', '2023-07-01', '2024-01-31', NOW(), NOW()),

('Security Awareness Training', 'SAT-2024', 'Enterprise security awareness training solution including online courses and phishing simulation', 'Sun Ba', 'Manufacturing Corp', 'suspended', '2024-01-15', '2024-05-31', NOW(), NOW()),

('Zero Trust Architecture Implementation', 'ZTA-2024', 'Enterprise zero trust security architecture planning and implementation project', 'Zhou Jiu', 'Medical Corp', 'active', '2024-02-15', '2024-09-30', NOW(), NOW()),

('Cloud Security Assessment', 'CSE-2024', 'Security assessment service for cloud platforms and cloud-native applications', 'Wu Shi', 'Education Corp', 'active', '2024-03-15', '2024-07-31', NOW(), NOW()); 