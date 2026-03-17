-- OneClickVirt 数据库初始化脚本
-- 使用方法: mysql -uroot [数据库名称] < init.sql

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 0. 创建基础表结构
-- ============================================

-- 创建域名相关表
CREATE TABLE IF NOT EXISTS `domains` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `instance_id` bigint unsigned NOT NULL,
  `domain` varchar(255) NOT NULL,
  `protocol` varchar(10) DEFAULT 'http',
  `internal_ip` varchar(45) NOT NULL,
  `internal_port` int NOT NULL,
  `external_port` int DEFAULT 0,
  `ssl` tinyint(1) DEFAULT 0,
  `status` int DEFAULT 0,
  `agent_id` bigint unsigned DEFAULT NULL,
  `expires_at` datetime DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `domain` (`domain`),
  INDEX `user_id` (`user_id`),
  INDEX `instance_id` (`instance_id`),
  INDEX `agent_id` (`agent_id`),
  INDEX `status` (`status`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `domain_configs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `max_domains_per_user` bigint DEFAULT 3,
  `max_domains_per_agent_user` bigint DEFAULT 5,
  `default_ttl` bigint DEFAULT 300,
  `auto_ssl` bigint DEFAULT 0,
  `allowed_suffixes` text,
  `dns_type` varchar(50) DEFAULT 'dnsmasq',
  `dns_config_path` varchar(255) DEFAULT '/etc/dnsmasq.d/oneclickvirt-hosts.conf',
  `nginx_config_path` varchar(255) DEFAULT '/etc/nginx/conf.d/oneclickvirt-domains',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建订单相关表
CREATE TABLE IF NOT EXISTS `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `order_no` varchar(50) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `status` int DEFAULT 0,
  `payment_method` varchar(20) DEFAULT NULL,
  `payment_transaction_id` varchar(100) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `order_no` (`order_no`),
  INDEX `user_id` (`user_id`),
  INDEX `product_id` (`product_id`),
  INDEX `status` (`status`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建兑换码相关表
CREATE TABLE IF NOT EXISTS `redemption_codes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL,
  `type` varchar(20) NOT NULL,
  `value` decimal(10,2) NOT NULL,
  `uses` int DEFAULT 0,
  `max_uses` int DEFAULT 1,
  `expires_at` datetime DEFAULT NULL,
  `status` int DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  INDEX `status` (`status`),
  INDEX `expires_at` (`expires_at`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建代理商相关表
CREATE TABLE IF NOT EXISTS `agents` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `code` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `status` int DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `user_id` (`user_id`),
  INDEX `status` (`status`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `agent_sub_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `agent_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `agent_id_user_id` (`agent_id`,`user_id`),
  INDEX `agent_id` (`agent_id`),
  INDEX `user_id` (`user_id`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `commissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `agent_id` bigint unsigned NOT NULL,
  `order_id` bigint unsigned NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `status` int DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `agent_id` (`agent_id`),
  INDEX `order_id` (`order_id`),
  INDEX `status` (`status`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建钱包相关表
CREATE TABLE IF NOT EXISTS `wallets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `balance` decimal(10,2) DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `wallet_transactions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `type` varchar(20) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `user_id` (`user_id`),
  INDEX `type` (`type`),
  INDEX `created_at` (`created_at`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建实名认证相关表
CREATE TABLE IF NOT EXISTS `kyc_verifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `real_name` varchar(50) NOT NULL,
  `id_card` varchar(20) NOT NULL,
  `status` int DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`),
  INDEX `status` (`status`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建产品购买记录相关表
CREATE TABLE IF NOT EXISTS `product_purchases` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `order_id` bigint unsigned DEFAULT NULL,
  `level` int NOT NULL,
  `start_date` datetime NOT NULL,
  `end_date` datetime DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_purchases_order_id` (`order_id`),
  KEY `idx_product_purchases_user_id` (`user_id`),
  KEY `idx_product_purchases_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
INSERT INTO `products` (`name`, `description`, `price`, `billing_cycle`, `cpu_limit`, `memory_limit`, `disk_limit`, `bandwidth_limit`, `instance_limit`, `features`, `status`, `sort_order`, `is_enabled`, `cpu`, `memory`, `disk`, `bandwidth`, `traffic`, `period`, `allow_repeat`, `created_at`, `updated_at`) VALUES
('入门套餐', '适合个人用户的基础套餐，包含基本的虚拟化功能', 0, 'monthly', 1, 512, 10240, 100, 1, '{}', 1, 1, 1, 1, 512, 10240, 100, 0, 30, 1, NOW(), NOW()),
('标准套餐', '适合小型团队的标准套餐，包含更多资源', 990, 'monthly', 2, 1024, 20480, 200, 3, '{}', 1, 2, 1, 2, 1024, 20480, 200, 0, 30, 1, NOW(), NOW()),
('专业套餐', '适合中型团队的专业套餐，包含完整功能', 2990, 'monthly', 4, 2048, 40960, 500, 5, '{}', 1, 3, 1, 4, 2048, 40960, 500, 0, 30, 1, NOW(), NOW()),
('企业套餐', '适合大型团队的企业套餐，包含无限资源', 9990, 'monthly', 8, 4096, 102400, 1000, 10, '{}', 1, 4, 1, 8, 4096, 102400, 1000, 0, 30, 1, NOW(), NOW());

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
('footer_text', '© 2025 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW()),
('icp_number', '', 'string', 'basic', 'ICP备案号', NOW(), NOW()),
('police_number', '', 'string', 'basic', '公安备案号', NOW(), NOW());

-- ============================================
-- 7. 创建默认域名配置
-- ============================================
INSERT INTO `domain_configs` (`max_domains_per_user`, `max_domains_per_agent_user`, `default_ttl`, `auto_ssl`, `allowed_suffixes`, `dns_type`, `dns_config_path`, `nginx_config_path`, `created_at`, `updated_at`) VALUES
(3, 5, 300, 0, '', 'dnsmasq', '/etc/dnsmasq.d/oneclickvirt-hosts.conf', '/etc/nginx/conf.d/oneclickvirt-domains', NOW(), NOW());

-- ============================================
-- 8. 新增系统配置项 (email verify + real name + agent)
-- ============================================
INSERT INTO `system_configs` (`key`, `value`, `description`, `created_at`, `updated_at`) VALUES
('enable_email_verification', 'false', '是否开启邮箱验证（注册后需验证邮箱）', NOW(), NOW()),
('email_activation_expire_hours', '24', '邮箱激活链接过期时间（小时）', NOW(), NOW()),
('enable_real_name', 'false', '是否开启实名认证', NOW(), NOW()),
('require_real_name', 'false', '是否强制实名认证后才能使用服务', NOW(), NOW()),
('enable_agent', 'true', '是否开启代理商功能', NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 9. 为现有用户创建钱包
-- ============================================
INSERT INTO `wallets` (`user_id`, `balance`) SELECT `id`, 0 FROM `users` WHERE `id` NOT IN (SELECT `user_id` FROM `wallets`);

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 初始化完成
-- 默认管理员账户:
--   用户名: admin
--   密码: admin123456
--
-- 所有必要的表都已创建:
--   domains          - 域名绑定表
--   domain_configs   - 域名系统配置表
--   orders           - 订单表
--   redemption_codes - 兑换码表
--   agents           - 代理商表
--   agent_sub_users  - 代理商子用户关系表
--   commissions      - 佣金记录表
--   wallets          - 钱包表
--   wallet_transactions - 钱包交易记录表
--   kyc_verifications - 实名认证记录表
--   product_purchases - 产品购买记录表
-- ============================================
