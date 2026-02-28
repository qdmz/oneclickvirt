-- OneClickVirt 数据库初始化脚本
-- 使用方法: mysql -uroot oneclickvirt < init.sql

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 1. 创建默认角色
-- ============================================
INSERT INTO `roles` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`) VALUES
('admin', 'admin', '系统管理员角色', 1, NOW(), NOW()),
('user', 'user', '普通用户角色', 1, NOW(), NOW());

-- ============================================
-- 2. 创建默认管理员账户
-- 默认密码: admin123456
-- ============================================
INSERT INTO `users` (`uuid`, `username`, `password`, `nickname`, `email`, `phone`, `status`, `level`, `level_expire_at`, `user_type`, `created_at`, `updated_at`) VALUES
('admin-uuid-001', 'admin', '$2a$10$AKvQPFPqSVBQWv6J0hCCeujsogZRK2dZReuye1ZZXJFEzWgL61IZK', '管理员', 'admin@example.com', '13800138000', 1, 5, '2099-12-31 23:59:59', 'admin', NOW(), NOW());

-- 创建用户角色关联
INSERT INTO `user_roles` (`user_id`, `role_id`, `created_at`, `updated_at`) VALUES
(1, 1, NOW(), NOW());

-- ============================================
-- 3. 创建默认公告
-- ============================================
INSERT INTO `announcements` (`title`, `content`, `content_html`, `type`, `status`, `priority`, `is_sticky`, `created_at`, `updated_at`) VALUES
('欢迎使用虚拟化管理平台', '欢迎使用虚拟化管理平台，支持Docker、LXD、Incus、Proxmox VE等多种虚拟化技术。本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。', '<p>欢迎使用虚拟化管理平台，支持<strong>Docker</strong>、<strong>LXD</strong>、<strong>Incus</strong>、<strong>Proxmox VE</strong>等多种虚拟化技术。</p><p>本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。</p>', 'homepage', 1, 10, 1, NOW(), NOW()),
('系统维护通知', '为了提供更好的服务质量，我们会定期进行系统维护。维护期间可能会影响部分功能的使用，请您谅解。', '<p>为了提供更好的服务质量，我们会定期进行系统维护。</p>', 'topbar', 1, 5, 0, NOW(), NOW()),
('新手使用指南', '如果您是第一次使用本平台，建议先阅读使用文档。您可以在右上角的帮助菜单中找到详细的操作指南。', '<p>如果您是第一次使用本平台，建议先阅读使用文档。</p>', 'homepage', 1, 8, 0, NOW(), NOW());

-- ============================================
-- 4. 创建默认产品套餐
-- ============================================
INSERT INTO `products` (`name`, `description`, `price`, `billing_cycle`, `cpu_limit`, `memory_limit`, `disk_limit`, `bandwidth_limit`, `instance_limit`, `features`, `status`, `sort_order`, `created_at`, `updated_at`) VALUES
('入门套餐', '适合个人用户的基础套餐，包含基本的虚拟化功能', 0.00, 'monthly', 1, 512, 10240, 100, 1, '{"cpu": "1核", "memory": "512MB", "disk": "10GB", "bandwidth": "100Mbps", "instances": "1个实例"}', 1, 1, NOW(), NOW()),
('标准套餐', '适合小型团队的标准套餐，包含更多资源', 9.90, 'monthly', 2, 1024, 20480, 200, 3, '{"cpu": "2核", "memory": "1GB", "disk": "20GB", "bandwidth": "200Mbps", "instances": "3个实例"}', 1, 2, NOW(), NOW()),
('专业套餐', '适合中型团队的专业套餐，包含完整功能', 29.90, 'monthly', 4, 2048, 40960, 500, 5, '{"cpu": "4核", "memory": "2GB", "disk": "40GB", "bandwidth": "500Mbps", "instances": "5个实例"}', 1, 3, NOW(), NOW()),
('企业套餐', '适合大型团队的企业套餐，包含无限资源', 99.90, 'monthly', 8, 4096, 102400, 1000, 10, '{"cpu": "8核", "memory": "4GB", "disk": "100GB", "bandwidth": "1000Mbps", "instances": "10个实例"}', 1, 4, NOW(), NOW());

-- ============================================
-- 5. 创建默认系统配置
-- ============================================
INSERT INTO `system_configs` (`key`, `value`, `description`, `created_at`, `updated_at`) VALUES
('site_name', 'OneClickVirt', '网站名称', NOW(), NOW()),
('site_description', '虚拟化管理平台', '网站描述', NOW(), NOW()),
('site_keywords', '虚拟化,Docker,LXD,Incus,Proxmox', '网站关键词', NOW(), NOW()),
('enable_registration', 'true', '是否开启注册', NOW(), NOW()),
('enable_email_verify', 'false', '是否开启邮箱验证', NOW(), NOW()),
('default_user_level', '1', '默认用户等级', NOW(), NOW()),
('max_instances_per_user', '10', '每个用户最大实例数', NOW(), NOW()),
('default_instance_expiry_days', '30', '默认实例过期天数', NOW(), NOW());

-- ============================================
-- 6. 创建默认站点配置
-- ============================================
INSERT INTO `site_configs` (`key`, `value`, `type`, `group`, `description`, `created_at`, `updated_at`) VALUES
('site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW()),
('site_icon_url', '/favicon.ico', 'string', 'basic', '网站图标URL', NOW(), NOW()),
('site_logo_url', '/logo.png', 'string', 'basic', '网站Logo URL', NOW(), NOW()),
('footer_text', '© 2024 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW()),
('icp_number', '', 'string', 'basic', 'ICP备案号', NOW(), NOW()),
('police_number', '', 'string', 'basic', '公安备案号', NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 初始化完成
-- 默认管理员账户:
--   用户名: admin
--   密码: admin123456
-- ============================================
