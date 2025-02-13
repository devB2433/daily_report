-- Clean existing data
DELETE FROM tasks;
DELETE FROM reports;

-- Add mock reports for admin (user_id = 1) and normal users
-- November 2024
INSERT INTO reports (user_id, status, date, submitted_at, created_at, updated_at) VALUES
-- Admin reports (Every Monday and Wednesday)
(1, 'submitted', '2024-11-04', '2024-11-04 17:30:00', '2024-11-04 17:30:00', '2024-11-04 17:30:00'),
(1, 'submitted', '2024-11-06', '2024-11-06 17:45:00', '2024-11-06 17:45:00', '2024-11-06 17:45:00'),
(1, 'submitted', '2024-11-11', '2024-11-11 17:30:00', '2024-11-11 17:30:00', '2024-11-11 17:30:00'),
(1, 'submitted', '2024-11-13', '2024-11-13 17:45:00', '2024-11-13 17:45:00', '2024-11-13 17:45:00'),
(1, 'submitted', '2024-11-18', '2024-11-18 17:30:00', '2024-11-18 17:30:00', '2024-11-18 17:30:00'),
(1, 'submitted', '2024-11-20', '2024-11-20 17:45:00', '2024-11-20 17:45:00', '2024-11-20 17:45:00'),
(1, 'submitted', '2024-11-25', '2024-11-25 17:30:00', '2024-11-25 17:30:00', '2024-11-25 17:30:00'),
(1, 'submitted', '2024-11-27', '2024-11-27 17:45:00', '2024-11-27 17:45:00', '2024-11-27 17:45:00'),

-- zhangsan reports (Every Tuesday and Thursday)
(2, 'submitted', '2024-11-05', '2024-11-05 17:30:00', '2024-11-05 17:30:00', '2024-11-05 17:30:00'),
(2, 'submitted', '2024-11-07', '2024-11-07 17:45:00', '2024-11-07 17:45:00', '2024-11-07 17:45:00'),
(2, 'submitted', '2024-11-12', '2024-11-12 17:30:00', '2024-11-12 17:30:00', '2024-11-12 17:30:00'),
(2, 'submitted', '2024-11-14', '2024-11-14 17:45:00', '2024-11-14 17:45:00', '2024-11-14 17:45:00'),
(2, 'submitted', '2024-11-19', '2024-11-19 17:30:00', '2024-11-19 17:30:00', '2024-11-19 17:30:00'),
(2, 'submitted', '2024-11-21', '2024-11-21 17:45:00', '2024-11-21 17:45:00', '2024-11-21 17:45:00'),
(2, 'submitted', '2024-11-26', '2024-11-26 17:30:00', '2024-11-26 17:30:00', '2024-11-26 17:30:00'),
(2, 'submitted', '2024-11-28', '2024-11-28 17:45:00', '2024-11-28 17:45:00', '2024-11-28 17:45:00'),

-- lisi reports (Every Monday and Wednesday)
(3, 'submitted', '2024-11-04', '2024-11-04 17:30:00', '2024-11-04 17:30:00', '2024-11-04 17:30:00'),
(3, 'submitted', '2024-11-06', '2024-11-06 17:45:00', '2024-11-06 17:45:00', '2024-11-06 17:45:00'),
(3, 'submitted', '2024-11-11', '2024-11-11 17:30:00', '2024-11-11 17:30:00', '2024-11-11 17:30:00'),
(3, 'submitted', '2024-11-13', '2024-11-13 17:45:00', '2024-11-13 17:45:00', '2024-11-13 17:45:00'),
(3, 'submitted', '2024-11-18', '2024-11-18 17:30:00', '2024-11-18 17:30:00', '2024-11-18 17:30:00'),
(3, 'submitted', '2024-11-20', '2024-11-20 17:45:00', '2024-11-20 17:45:00', '2024-11-20 17:45:00'),
(3, 'submitted', '2024-11-25', '2024-11-25 17:30:00', '2024-11-25 17:30:00', '2024-11-25 17:30:00'),
(3, 'submitted', '2024-11-27', '2024-11-27 17:45:00', '2024-11-27 17:45:00', '2024-11-27 17:45:00');

-- December 2024 (Similar pattern for other users...)
INSERT INTO reports (user_id, status, date, submitted_at, created_at, updated_at) VALUES
-- Admin reports
(1, 'submitted', '2024-12-02', '2024-12-02 17:30:00', '2024-12-02 17:30:00', '2024-12-02 17:30:00'),
(1, 'submitted', '2024-12-04', '2024-12-04 17:45:00', '2024-12-04 17:45:00', '2024-12-04 17:45:00'),
(1, 'submitted', '2024-12-09', '2024-12-09 17:30:00', '2024-12-09 17:30:00', '2024-12-09 17:30:00'),
(1, 'submitted', '2024-12-11', '2024-12-11 17:45:00', '2024-12-11 17:45:00', '2024-12-11 17:45:00'),
(1, 'submitted', '2024-12-16', '2024-12-16 17:30:00', '2024-12-16 17:30:00', '2024-12-16 17:30:00'),
(1, 'submitted', '2024-12-18', '2024-12-18 17:45:00', '2024-12-18 17:45:00', '2024-12-18 17:45:00'),
(1, 'submitted', '2024-12-23', '2024-12-23 17:30:00', '2024-12-23 17:30:00', '2024-12-23 17:30:00'),
(1, 'submitted', '2024-12-25', '2024-12-25 17:45:00', '2024-12-25 17:45:00', '2024-12-25 17:45:00');

-- Add tasks for each report
-- Admin's tasks (focusing on system management and security projects)
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at) 
SELECT 
    r.id,
    pid.project_id,
    ROUND(RAND() * 4 + 4, 1) as hours,
    CASE 
        WHEN pid.project_id = 1 THEN 'Conducted security assessment and vulnerability scanning for web applications'
        ELSE 'Performed mobile app security testing and code review'
    END as content,
    'completed' as status,
    r.created_at,
    r.updated_at
FROM reports r
CROSS JOIN (SELECT CASE WHEN RAND() < 0.5 THEN 1 ELSE 2 END as project_id) pid
WHERE r.user_id = 1;

-- zhangsan's tasks (Web Security Scanner System project)
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at)
SELECT 
    r.id,
    1 as project_id,
    ROUND(RAND() * 2 + 6, 1) as hours,
    CASE WHEN RAND() < 0.5 
        THEN 'Implemented new vulnerability detection modules for the scanner'
        ELSE 'Enhanced report generation functionality and fixed scanner bugs'
    END as content,
    'completed' as status,
    r.created_at,
    r.updated_at
FROM reports r
WHERE r.user_id = 2;

-- lisi's tasks (Mobile App Security Testing Platform project)
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at)
SELECT 
    r.id,
    2 as project_id,
    ROUND(RAND() * 2 + 6, 1) as hours,
    CASE WHEN RAND() < 0.5 
        THEN 'Developed new test cases for iOS application security testing'
        ELSE 'Improved Android dynamic analysis capabilities'
    END as content,
    'completed' as status,
    r.created_at,
    r.updated_at
FROM reports r
WHERE r.user_id = 3;

-- January 2025
INSERT INTO reports (user_id, status, date, submitted_at, created_at, updated_at) VALUES
-- Admin reports
(1, 'submitted', '2025-01-06', '2025-01-06 17:30:00', '2025-01-06 17:30:00', '2025-01-06 17:30:00'),
(1, 'submitted', '2025-01-08', '2025-01-08 17:45:00', '2025-01-08 17:45:00', '2025-01-08 17:45:00'),
(1, 'submitted', '2025-01-13', '2025-01-13 17:30:00', '2025-01-13 17:30:00', '2025-01-13 17:30:00'),
(1, 'submitted', '2025-01-15', '2025-01-15 17:45:00', '2025-01-15 17:45:00', '2025-01-15 17:45:00'),
(1, 'submitted', '2025-01-20', '2025-01-20 17:30:00', '2025-01-20 17:30:00', '2025-01-20 17:30:00'),
(1, 'submitted', '2025-01-22', '2025-01-22 17:45:00', '2025-01-22 17:45:00', '2025-01-22 17:45:00'),
(1, 'submitted', '2025-01-27', '2025-01-27 17:30:00', '2025-01-27 17:30:00', '2025-01-27 17:30:00'),
(1, 'submitted', '2025-01-29', '2025-01-29 17:45:00', '2025-01-29 17:45:00', '2025-01-29 17:45:00'),

-- zhangsan reports
(2, 'submitted', '2025-01-07', '2025-01-07 17:30:00', '2025-01-07 17:30:00', '2025-01-07 17:30:00'),
(2, 'submitted', '2025-01-09', '2025-01-09 17:45:00', '2025-01-09 17:45:00', '2025-01-09 17:45:00'),
(2, 'submitted', '2025-01-14', '2025-01-14 17:30:00', '2025-01-14 17:30:00', '2025-01-14 17:30:00'),
(2, 'submitted', '2025-01-16', '2025-01-16 17:45:00', '2025-01-16 17:45:00', '2025-01-16 17:45:00'),
(2, 'submitted', '2025-01-21', '2025-01-21 17:30:00', '2025-01-21 17:30:00', '2025-01-21 17:30:00'),
(2, 'submitted', '2025-01-23', '2025-01-23 17:45:00', '2025-01-23 17:45:00', '2025-01-23 17:45:00'),
(2, 'submitted', '2025-01-28', '2025-01-28 17:30:00', '2025-01-28 17:30:00', '2025-01-28 17:30:00'),
(2, 'submitted', '2025-01-30', '2025-01-30 17:45:00', '2025-01-30 17:45:00', '2025-01-30 17:45:00'),

-- lisi reports
(3, 'submitted', '2025-01-06', '2025-01-06 17:30:00', '2025-01-06 17:30:00', '2025-01-06 17:30:00'),
(3, 'submitted', '2025-01-08', '2025-01-08 17:45:00', '2025-01-08 17:45:00', '2025-01-08 17:45:00'),
(3, 'submitted', '2025-01-13', '2025-01-13 17:30:00', '2025-01-13 17:30:00', '2025-01-13 17:30:00'),
(3, 'submitted', '2025-01-15', '2025-01-15 17:45:00', '2025-01-15 17:45:00', '2025-01-15 17:45:00'),
(3, 'submitted', '2025-01-20', '2025-01-20 17:30:00', '2025-01-20 17:30:00', '2025-01-20 17:30:00'),
(3, 'submitted', '2025-01-22', '2025-01-22 17:45:00', '2025-01-22 17:45:00', '2025-01-22 17:45:00'),
(3, 'submitted', '2025-01-27', '2025-01-27 17:30:00', '2025-01-27 17:30:00', '2025-01-27 17:30:00'),
(3, 'submitted', '2025-01-29', '2025-01-29 17:45:00', '2025-01-29 17:45:00', '2025-01-29 17:45:00');

-- February 2025
INSERT INTO reports (user_id, status, date, submitted_at, created_at, updated_at) VALUES
-- Admin reports
(1, 'submitted', '2025-02-03', '2025-02-03 17:30:00', '2025-02-03 17:30:00', '2025-02-03 17:30:00'),
(1, 'submitted', '2025-02-05', '2025-02-05 17:45:00', '2025-02-05 17:45:00', '2025-02-05 17:45:00'),
(1, 'submitted', '2025-02-10', '2025-02-10 17:30:00', '2025-02-10 17:30:00', '2025-02-10 17:30:00'),
(1, 'submitted', '2025-02-12', '2025-02-12 17:45:00', '2025-02-12 17:45:00', '2025-02-12 17:45:00'),

-- zhangsan reports
(2, 'submitted', '2025-02-04', '2025-02-04 17:30:00', '2025-02-04 17:30:00', '2025-02-04 17:30:00'),
(2, 'submitted', '2025-02-06', '2025-02-06 17:45:00', '2025-02-06 17:45:00', '2025-02-06 17:45:00'),
(2, 'submitted', '2025-02-11', '2025-02-11 17:30:00', '2025-02-11 17:30:00', '2025-02-11 17:30:00'),
(2, 'submitted', '2025-02-13', '2025-02-13 17:45:00', '2025-02-13 17:45:00', '2025-02-13 17:45:00'),

-- lisi reports
(3, 'submitted', '2025-02-03', '2025-02-03 17:30:00', '2025-02-03 17:30:00', '2025-02-03 17:30:00'),
(3, 'submitted', '2025-02-05', '2025-02-05 17:45:00', '2025-02-05 17:45:00', '2025-02-05 17:45:00'),
(3, 'submitted', '2025-02-10', '2025-02-10 17:30:00', '2025-02-10 17:30:00', '2025-02-10 17:30:00'),
(3, 'submitted', '2025-02-12', '2025-02-12 17:45:00', '2025-02-12 17:45:00', '2025-02-12 17:45:00');

-- Tasks for January 2025 reports
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at) 
SELECT 
  r.id,
  CASE 
    WHEN r.user_id = 1 THEN 1  -- Admin works on Web Security Scanner
    WHEN r.user_id = 2 THEN 2  -- zhangsan works on Mobile App Security Testing
    WHEN r.user_id = 3 THEN 3  -- lisi works on Security Operations Center
  END as project_id,
  8.0 as hours,
  CASE 
    WHEN r.user_id = 1 THEN 'Conducted security assessment of web application modules. Identified and documented potential vulnerabilities. Updated scanning rules.'
    WHEN r.user_id = 2 THEN 'Enhanced mobile app testing procedures. Implemented new test cases for API security. Reviewed and updated documentation.'
    WHEN r.user_id = 3 THEN 'Monitored security alerts and incidents. Performed threat analysis. Updated incident response procedures.'
  END as content,
  'completed' as status,
  r.created_at,
  r.updated_at
FROM reports r
WHERE r.submitted_at >= '2025-01-01' AND r.submitted_at < '2025-02-01';

-- Tasks for February 2025 reports
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at) 
SELECT 
  r.id,
  CASE 
    WHEN r.user_id = 1 THEN 1  -- Admin works on Web Security Scanner
    WHEN r.user_id = 2 THEN 2  -- zhangsan works on Mobile App Security Testing
    WHEN r.user_id = 3 THEN 3  -- lisi works on Security Operations Center
  END as project_id,
  8.0 as hours,
  CASE 
    WHEN r.user_id = 1 THEN 'Implemented new security scanning features. Optimized scanning performance. Updated vulnerability database.'
    WHEN r.user_id = 2 THEN 'Developed automated security test scripts. Conducted penetration testing on mobile apps. Updated security guidelines.'
    WHEN r.user_id = 3 THEN 'Managed security operations workflow. Investigated security incidents. Updated security metrics dashboard.'
  END as content,
  'completed' as status,
  r.created_at,
  r.updated_at
FROM reports r
WHERE r.submitted_at >= '2025-02-01' AND r.submitted_at < '2025-03-01';

-- Add reports for other users (qianqi, sunba, zhoujiu, wushi)
-- November 2024
INSERT INTO reports (user_id, status, date, submitted_at, created_at, updated_at) VALUES
-- qianqi reports (Every Tuesday and Thursday)
(6, 'submitted', '2024-11-05', '2024-11-05 17:30:00', '2024-11-05 17:30:00', '2024-11-05 17:30:00'),
(6, 'submitted', '2024-11-07', '2024-11-07 17:45:00', '2024-11-07 17:45:00', '2024-11-07 17:45:00'),
(6, 'submitted', '2024-11-12', '2024-11-12 17:30:00', '2024-11-12 17:30:00', '2024-11-12 17:30:00'),
(6, 'submitted', '2024-11-14', '2024-11-14 17:45:00', '2024-11-14 17:45:00', '2024-11-14 17:45:00'),
(6, 'submitted', '2024-11-19', '2024-11-19 17:30:00', '2024-11-19 17:30:00', '2024-11-19 17:30:00'),
(6, 'submitted', '2024-11-21', '2024-11-21 17:45:00', '2024-11-21 17:45:00', '2024-11-21 17:45:00'),
(6, 'submitted', '2024-11-26', '2024-11-26 17:30:00', '2024-11-26 17:30:00', '2024-11-26 17:30:00'),
(6, 'submitted', '2024-11-28', '2024-11-28 17:45:00', '2024-11-28 17:45:00', '2024-11-28 17:45:00'),

-- sunba reports (Every Monday and Wednesday)
(7, 'submitted', '2024-11-04', '2024-11-04 17:30:00', '2024-11-04 17:30:00', '2024-11-04 17:30:00'),
(7, 'submitted', '2024-11-06', '2024-11-06 17:45:00', '2024-11-06 17:45:00', '2024-11-06 17:45:00'),
(7, 'submitted', '2024-11-11', '2024-11-11 17:30:00', '2024-11-11 17:30:00', '2024-11-11 17:30:00'),
(7, 'submitted', '2024-11-13', '2024-11-13 17:45:00', '2024-11-13 17:45:00', '2024-11-13 17:45:00'),
(7, 'submitted', '2024-11-18', '2024-11-18 17:30:00', '2024-11-18 17:30:00', '2024-11-18 17:30:00'),
(7, 'submitted', '2024-11-20', '2024-11-20 17:45:00', '2024-11-20 17:45:00', '2024-11-20 17:45:00'),
(7, 'submitted', '2024-11-25', '2024-11-25 17:30:00', '2024-11-25 17:30:00', '2024-11-25 17:30:00'),
(7, 'submitted', '2024-11-27', '2024-11-27 17:45:00', '2024-11-27 17:45:00', '2024-11-27 17:45:00'),

-- zhoujiu reports (Every Tuesday and Thursday)
(8, 'submitted', '2024-11-05', '2024-11-05 17:30:00', '2024-11-05 17:30:00', '2024-11-05 17:30:00'),
(8, 'submitted', '2024-11-07', '2024-11-07 17:45:00', '2024-11-07 17:45:00', '2024-11-07 17:45:00'),
(8, 'submitted', '2024-11-12', '2024-11-12 17:30:00', '2024-11-12 17:30:00', '2024-11-12 17:30:00'),
(8, 'submitted', '2024-11-14', '2024-11-14 17:45:00', '2024-11-14 17:45:00', '2024-11-14 17:45:00'),
(8, 'submitted', '2024-11-19', '2024-11-19 17:30:00', '2024-11-19 17:30:00', '2024-11-19 17:30:00'),
(8, 'submitted', '2024-11-21', '2024-11-21 17:45:00', '2024-11-21 17:45:00', '2024-11-21 17:45:00'),
(8, 'submitted', '2024-11-26', '2024-11-26 17:30:00', '2024-11-26 17:30:00', '2024-11-26 17:30:00'),
(8, 'submitted', '2024-11-28', '2024-11-28 17:45:00', '2024-11-28 17:45:00', '2024-11-28 17:45:00'),

-- wushi reports (Every Monday and Wednesday)
(9, 'submitted', '2024-11-04', '2024-11-04 17:30:00', '2024-11-04 17:30:00', '2024-11-04 17:30:00'),
(9, 'submitted', '2024-11-06', '2024-11-06 17:45:00', '2024-11-06 17:45:00', '2024-11-06 17:45:00'),
(9, 'submitted', '2024-11-11', '2024-11-11 17:30:00', '2024-11-11 17:30:00', '2024-11-11 17:30:00'),
(9, 'submitted', '2024-11-13', '2024-11-13 17:45:00', '2024-11-13 17:45:00', '2024-11-13 17:45:00'),
(9, 'submitted', '2024-11-18', '2024-11-18 17:30:00', '2024-11-18 17:30:00', '2024-11-18 17:30:00'),
(9, 'submitted', '2024-11-20', '2024-11-20 17:45:00', '2024-11-20 17:45:00', '2024-11-20 17:45:00'),
(9, 'submitted', '2024-11-25', '2024-11-25 17:30:00', '2024-11-25 17:30:00', '2024-11-25 17:30:00'),
(9, 'submitted', '2024-11-27', '2024-11-27 17:45:00', '2024-11-27 17:45:00', '2024-11-27 17:45:00');

-- December 2024
INSERT INTO reports (user_id, status, date, submitted_at, created_at, updated_at) VALUES
-- qianqi reports
(6, 'submitted', '2024-12-03', '2024-12-03 17:30:00', '2024-12-03 17:30:00', '2024-12-03 17:30:00'),
(6, 'submitted', '2024-12-05', '2024-12-05 17:45:00', '2024-12-05 17:45:00', '2024-12-05 17:45:00'),
(6, 'submitted', '2024-12-10', '2024-12-10 17:30:00', '2024-12-10 17:30:00', '2024-12-10 17:30:00'),
(6, 'submitted', '2024-12-12', '2024-12-12 17:45:00', '2024-12-12 17:45:00', '2024-12-12 17:45:00'),
(6, 'submitted', '2024-12-17', '2024-12-17 17:30:00', '2024-12-17 17:30:00', '2024-12-17 17:30:00'),
(6, 'submitted', '2024-12-19', '2024-12-19 17:45:00', '2024-12-19 17:45:00', '2024-12-19 17:45:00'),
(6, 'submitted', '2024-12-24', '2024-12-24 17:30:00', '2024-12-24 17:30:00', '2024-12-24 17:30:00'),
(6, 'submitted', '2024-12-26', '2024-12-26 17:45:00', '2024-12-26 17:45:00', '2024-12-26 17:45:00'),

-- sunba reports
(7, 'submitted', '2024-12-02', '2024-12-02 17:30:00', '2024-12-02 17:30:00', '2024-12-02 17:30:00'),
(7, 'submitted', '2024-12-04', '2024-12-04 17:45:00', '2024-12-04 17:45:00', '2024-12-04 17:45:00'),
(7, 'submitted', '2024-12-09', '2024-12-09 17:30:00', '2024-12-09 17:30:00', '2024-12-09 17:30:00'),
(7, 'submitted', '2024-12-11', '2024-12-11 17:45:00', '2024-12-11 17:45:00', '2024-12-11 17:45:00'),
(7, 'submitted', '2024-12-16', '2024-12-16 17:30:00', '2024-12-16 17:30:00', '2024-12-16 17:30:00'),
(7, 'submitted', '2024-12-18', '2024-12-18 17:45:00', '2024-12-18 17:45:00', '2024-12-18 17:45:00'),
(7, 'submitted', '2024-12-23', '2024-12-23 17:30:00', '2024-12-23 17:30:00', '2024-12-23 17:30:00'),
(7, 'submitted', '2024-12-25', '2024-12-25 17:45:00', '2024-12-25 17:45:00', '2024-12-25 17:45:00'),

-- zhoujiu reports
(8, 'submitted', '2024-12-03', '2024-12-03 17:30:00', '2024-12-03 17:30:00', '2024-12-03 17:30:00'),
(8, 'submitted', '2024-12-05', '2024-12-05 17:45:00', '2024-12-05 17:45:00', '2024-12-05 17:45:00'),
(8, 'submitted', '2024-12-10', '2024-12-10 17:30:00', '2024-12-10 17:30:00', '2024-12-10 17:30:00'),
(8, 'submitted', '2024-12-12', '2024-12-12 17:45:00', '2024-12-12 17:45:00', '2024-12-12 17:45:00'),
(8, 'submitted', '2024-12-17', '2024-12-17 17:30:00', '2024-12-17 17:30:00', '2024-12-17 17:30:00'),
(8, 'submitted', '2024-12-19', '2024-12-19 17:45:00', '2024-12-19 17:45:00', '2024-12-19 17:45:00'),
(8, 'submitted', '2024-12-24', '2024-12-24 17:30:00', '2024-12-24 17:30:00', '2024-12-24 17:30:00'),
(8, 'submitted', '2024-12-26', '2024-12-26 17:45:00', '2024-12-26 17:45:00', '2024-12-26 17:45:00'),

-- wushi reports
(9, 'submitted', '2024-12-02', '2024-12-02 17:30:00', '2024-12-02 17:30:00', '2024-12-02 17:30:00'),
(9, 'submitted', '2024-12-04', '2024-12-04 17:45:00', '2024-12-04 17:45:00', '2024-12-04 17:45:00'),
(9, 'submitted', '2024-12-09', '2024-12-09 17:30:00', '2024-12-09 17:30:00', '2024-12-09 17:30:00'),
(9, 'submitted', '2024-12-11', '2024-12-11 17:45:00', '2024-12-11 17:45:00', '2024-12-11 17:45:00'),
(9, 'submitted', '2024-12-16', '2024-12-16 17:30:00', '2024-12-16 17:30:00', '2024-12-16 17:30:00'),
(9, 'submitted', '2024-12-18', '2024-12-18 17:45:00', '2024-12-18 17:45:00', '2024-12-18 17:45:00'),
(9, 'submitted', '2024-12-23', '2024-12-23 17:30:00', '2024-12-23 17:30:00', '2024-12-23 17:30:00'),
(9, 'submitted', '2024-12-25', '2024-12-25 17:45:00', '2024-12-25 17:45:00', '2024-12-25 17:45:00');

-- January 2025
INSERT INTO reports (user_id, status, date, submitted_at, created_at, updated_at) VALUES
-- qianqi reports
(6, 'submitted', '2025-01-07', '2025-01-07 17:30:00', '2025-01-07 17:30:00', '2025-01-07 17:30:00'),
(6, 'submitted', '2025-01-09', '2025-01-09 17:45:00', '2025-01-09 17:45:00', '2025-01-09 17:45:00'),
(6, 'submitted', '2025-01-14', '2025-01-14 17:30:00', '2025-01-14 17:30:00', '2025-01-14 17:30:00'),
(6, 'submitted', '2025-01-16', '2025-01-16 17:45:00', '2025-01-16 17:45:00', '2025-01-16 17:45:00'),
(6, 'submitted', '2025-01-21', '2025-01-21 17:30:00', '2025-01-21 17:30:00', '2025-01-21 17:30:00'),
(6, 'submitted', '2025-01-23', '2025-01-23 17:45:00', '2025-01-23 17:45:00', '2025-01-23 17:45:00'),
(6, 'submitted', '2025-01-28', '2025-01-28 17:30:00', '2025-01-28 17:30:00', '2025-01-28 17:30:00'),
(6, 'submitted', '2025-01-30', '2025-01-30 17:45:00', '2025-01-30 17:45:00', '2025-01-30 17:45:00'),

-- sunba reports
(7, 'submitted', '2025-01-06', '2025-01-06 17:30:00', '2025-01-06 17:30:00', '2025-01-06 17:30:00'),
(7, 'submitted', '2025-01-08', '2025-01-08 17:45:00', '2025-01-08 17:45:00', '2025-01-08 17:45:00'),
(7, 'submitted', '2025-01-13', '2025-01-13 17:30:00', '2025-01-13 17:30:00', '2025-01-13 17:30:00'),
(7, 'submitted', '2025-01-15', '2025-01-15 17:45:00', '2025-01-15 17:45:00', '2025-01-15 17:45:00'),
(7, 'submitted', '2025-01-20', '2025-01-20 17:30:00', '2025-01-20 17:30:00', '2025-01-20 17:30:00'),
(7, 'submitted', '2025-01-22', '2025-01-22 17:45:00', '2025-01-22 17:45:00', '2025-01-22 17:45:00'),
(7, 'submitted', '2025-01-27', '2025-01-27 17:30:00', '2025-01-27 17:30:00', '2025-01-27 17:30:00'),
(7, 'submitted', '2025-01-29', '2025-01-29 17:45:00', '2025-01-29 17:45:00', '2025-01-29 17:45:00'),

-- zhoujiu reports
(8, 'submitted', '2025-01-07', '2025-01-07 17:30:00', '2025-01-07 17:30:00', '2025-01-07 17:30:00'),
(8, 'submitted', '2025-01-09', '2025-01-09 17:45:00', '2025-01-09 17:45:00', '2025-01-09 17:45:00'),
(8, 'submitted', '2025-01-14', '2025-01-14 17:30:00', '2025-01-14 17:30:00', '2025-01-14 17:30:00'),
(8, 'submitted', '2025-01-16', '2025-01-16 17:45:00', '2025-01-16 17:45:00', '2025-01-16 17:45:00'),
(8, 'submitted', '2025-01-21', '2025-01-21 17:30:00', '2025-01-21 17:30:00', '2025-01-21 17:30:00'),
(8, 'submitted', '2025-01-23', '2025-01-23 17:45:00', '2025-01-23 17:45:00', '2025-01-23 17:45:00'),
(8, 'submitted', '2025-01-28', '2025-01-28 17:30:00', '2025-01-28 17:30:00', '2025-01-28 17:30:00'),
(8, 'submitted', '2025-01-30', '2025-01-30 17:45:00', '2025-01-30 17:45:00', '2025-01-30 17:45:00'),

-- wushi reports
(9, 'submitted', '2025-01-06', '2025-01-06 17:30:00', '2025-01-06 17:30:00', '2025-01-06 17:30:00'),
(9, 'submitted', '2025-01-08', '2025-01-08 17:45:00', '2025-01-08 17:45:00', '2025-01-08 17:45:00'),
(9, 'submitted', '2025-01-13', '2025-01-13 17:30:00', '2025-01-13 17:30:00', '2025-01-13 17:30:00'),
(9, 'submitted', '2025-01-15', '2025-01-15 17:45:00', '2025-01-15 17:45:00', '2025-01-15 17:45:00'),
(9, 'submitted', '2025-01-20', '2025-01-20 17:30:00', '2025-01-20 17:30:00', '2025-01-20 17:30:00'),
(9, 'submitted', '2025-01-22', '2025-01-22 17:45:00', '2025-01-22 17:45:00', '2025-01-22 17:45:00'),
(9, 'submitted', '2025-01-27', '2025-01-27 17:30:00', '2025-01-27 17:30:00', '2025-01-27 17:30:00'),
(9, 'submitted', '2025-01-29', '2025-01-29 17:45:00', '2025-01-29 17:45:00', '2025-01-29 17:45:00');

-- February 2025
INSERT INTO reports (user_id, status, date, submitted_at, created_at, updated_at) VALUES
-- qianqi reports
(6, 'submitted', '2025-02-04', '2025-02-04 17:30:00', '2025-02-04 17:30:00', '2025-02-04 17:30:00'),
(6, 'submitted', '2025-02-06', '2025-02-06 17:45:00', '2025-02-06 17:45:00', '2025-02-06 17:45:00'),
(6, 'submitted', '2025-02-11', '2025-02-11 17:30:00', '2025-02-11 17:30:00', '2025-02-11 17:30:00'),
(6, 'submitted', '2025-02-13', '2025-02-13 17:45:00', '2025-02-13 17:45:00', '2025-02-13 17:45:00'),

-- sunba reports
(7, 'submitted', '2025-02-03', '2025-02-03 17:30:00', '2025-02-03 17:30:00', '2025-02-03 17:30:00'),
(7, 'submitted', '2025-02-05', '2025-02-05 17:45:00', '2025-02-05 17:45:00', '2025-02-05 17:45:00'),
(7, 'submitted', '2025-02-10', '2025-02-10 17:30:00', '2025-02-10 17:30:00', '2025-02-10 17:30:00'),
(7, 'submitted', '2025-02-12', '2025-02-12 17:45:00', '2025-02-12 17:45:00', '2025-02-12 17:45:00'),

-- zhoujiu reports
(8, 'submitted', '2025-02-04', '2025-02-04 17:30:00', '2025-02-04 17:30:00', '2025-02-04 17:30:00'),
(8, 'submitted', '2025-02-06', '2025-02-06 17:45:00', '2025-02-06 17:45:00', '2025-02-06 17:45:00'),
(8, 'submitted', '2025-02-11', '2025-02-11 17:30:00', '2025-02-11 17:30:00', '2025-02-11 17:30:00'),
(8, 'submitted', '2025-02-13', '2025-02-13 17:45:00', '2025-02-13 17:45:00', '2025-02-13 17:45:00'),

-- wushi reports
(9, 'submitted', '2025-02-03', '2025-02-03 17:30:00', '2025-02-03 17:30:00', '2025-02-03 17:30:00'),
(9, 'submitted', '2025-02-05', '2025-02-05 17:45:00', '2025-02-05 17:45:00', '2025-02-05 17:45:00'),
(9, 'submitted', '2025-02-10', '2025-02-10 17:30:00', '2025-02-10 17:30:00', '2025-02-10 17:30:00'),
(9, 'submitted', '2025-02-12', '2025-02-12 17:45:00', '2025-02-12 17:45:00', '2025-02-12 17:45:00');

-- Add tasks for other users
-- qianqi's tasks (Security Operations Center project)
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at)
SELECT 
    r.id,
    3 as project_id,
    ROUND(RAND() * 2 + 6, 1) as hours,
    CASE WHEN RAND() < 0.5 
        THEN 'Monitored security alerts and performed incident analysis'
        ELSE 'Updated security monitoring rules and response procedures'
    END as content,
    'completed' as status,
    r.created_at,
    r.updated_at
FROM reports r
WHERE r.user_id = 6;

-- sunba's tasks (Security Awareness Training project)
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at)
SELECT 
    r.id,
    6 as project_id,
    ROUND(RAND() * 2 + 6, 1) as hours,
    CASE WHEN RAND() < 0.5 
        THEN 'Developed security awareness training materials and conducted training sessions'
        ELSE 'Created phishing simulation scenarios and analyzed results'
    END as content,
    'completed' as status,
    r.created_at,
    r.updated_at
FROM reports r
WHERE r.user_id = 7;

-- zhoujiu's tasks (Zero Trust Architecture Implementation project)
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at)
SELECT 
    r.id,
    7 as project_id,
    ROUND(RAND() * 2 + 6, 1) as hours,
    CASE WHEN RAND() < 0.5 
        THEN 'Designed zero trust network architecture and access policies'
        ELSE 'Implemented identity and access management solutions'
    END as content,
    'completed' as status,
    r.created_at,
    r.updated_at
FROM reports r
WHERE r.user_id = 8;

-- wushi's tasks (Cloud Security Assessment project)
INSERT INTO tasks (report_id, project_id, hours, content, status, created_at, updated_at)
SELECT 
    r.id,
    8 as project_id,
    ROUND(RAND() * 2 + 6, 1) as hours,
    CASE WHEN RAND() < 0.5 
        THEN 'Conducted security assessment of cloud infrastructure and services'
        ELSE 'Reviewed cloud security configurations and provided recommendations'
    END as content,
    'completed' as status,
    r.created_at,
    r.updated_at
FROM reports r
WHERE r.user_id = 9; 