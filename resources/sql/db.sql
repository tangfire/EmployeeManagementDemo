use employee_db;



-- 插入部门数据（5个部门）
INSERT INTO departments (depart) VALUES
                                     ('技术部'),
                                     ('市场部'),
                                     ('财务部'),
                                     ('人力资源部'),
                                     ('销售部');

-- 插入员工数据（每个部门8-12人，薪资范围不同）
INSERT INTO employees (dep_id, username, password, position, gender, email, phone, salary) VALUES
-- 技术部（高薪资范围：10000-30000）
((SELECT dep_id FROM departments WHERE depart = '技术部'), 'tech1', 'hashed_password', '高级工程师', '男', 'tech1@tech.com', '13800000001', 28000),
((SELECT dep_id FROM departments WHERE depart = '技术部'), 'tech2', 'hashed_password', '前端工程师', '女', 'tech2@tech.com', '13800000002', 22000),
((SELECT dep_id FROM departments WHERE depart = '技术部'), 'tech3', 'hashed_password', '后端工程师', '男', 'tech3@tech.com', '13800000003', 26000),
((SELECT dep_id FROM departments WHERE depart = '技术部'), 'tech4', 'hashed_password', '架构师', '男', 'tech4@tech.com', '13800000004', 30000),
((SELECT dep_id FROM departments WHERE depart = '技术部'), 'tech5', 'hashed_password', '测试工程师', '女', 'tech5@tech.com', '13800000005', 18000),

-- 市场部（中等薪资范围：8000-20000）
((SELECT dep_id FROM departments WHERE depart = '市场部'), 'mkt1', 'hashed_password', '市场经理', '女', 'mkt1@mkt.com', '13800000006', 20000),
((SELECT dep_id FROM departments WHERE depart = '市场部'), 'mkt2', 'hashed_password', '市场专员', '男', 'mkt2@mkt.com', '13800000007', 12000),
((SELECT dep_id FROM departments WHERE depart = '市场部'), 'mkt3', 'hashed_password', '品牌策划', '女', 'mkt3@mkt.com', '13800000008', 15000),

-- 财务部（稳定薪资范围：10000-25000）
((SELECT dep_id FROM departments WHERE depart = '财务部'), 'finance1', 'hashed_password', '财务总监', '女', 'finance1@finance.com', '13800000009', 25000),
((SELECT dep_id FROM departments WHERE depart = '财务部'), 'finance2', 'hashed_password', '会计', '男', 'finance2@finance.com', '13800000010', 15000),
((SELECT dep_id FROM departments WHERE depart = '财务部'), 'finance3', 'hashed_password', '出纳', '女', 'finance3@finance.com', '13800000011', 12000),

-- 人力资源部（中等薪资范围：8000-18000）
((SELECT dep_id FROM departments WHERE depart = '人力资源部'), 'hr1', 'hashed_password', 'HRBP', '女', 'hr1@hr.com', '13800000012', 18000),
((SELECT dep_id FROM departments WHERE depart = '人力资源部'), 'hr2', 'hashed_password', '招聘经理', '男', 'hr2@hr.com', '13800000013', 16000),

-- 销售部（底薪+提成模式，此处模拟底薪：5000-15000）
((SELECT dep_id FROM departments WHERE depart = '销售部'), 'sales1', 'hashed_password', '销售总监', '男', 'sales1@sales.com', '13800000014', 15000),
((SELECT dep_id FROM departments WHERE depart = '销售部'), 'sales2', 'hashed_password', '大客户经理', '女', 'sales2@sales.com', '13800000015', 12000),
((SELECT dep_id FROM departments WHERE depart = '销售部'), 'sales3', 'hashed_password', '销售代表', '男', 'sales3@sales.com', '13800000016', 8000);
