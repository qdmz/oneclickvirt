-- OneClickVirt 自动修复脚本 (autosetup.sql)
-- 此脚本用于修复已部署系统的数据库问题
-- 使用方法: mysql -uroot -p密码 oneclickvirt < autosetup.sql

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 1. 修复 redemption_codes 表缺失字段
-- ============================================
ALTER TABLE redemption_codes ADD COLUMN IF NOT EXISTS `used_count` int DEFAULT 0 COMMENT '已使用次数' AFTER `max_uses`;

-- ============================================
-- 2. 修复 product_purchases 表外键问题
-- ============================================
ALTER TABLE product_purchases DROP INDEX IF EXISTS `idx_product_purchases_order_id`;
ALTER TABLE product_purchases MODIFY COLUMN `order_id` bigint unsigned DEFAULT NULL;

-- ============================================
-- 3. 确保 users 表有 balance, total_spent, total_orders 字段
-- ============================================
ALTER TABLE users ADD COLUMN IF NOT EXISTS `balance` decimal(10,2) DEFAULT 0.00 AFTER `user_type`;
ALTER TABLE users ADD COLUMN IF NOT EXISTS `total_spent` decimal(10,2) DEFAULT 0.00 AFTER `balance`;
ALTER TABLE users ADD COLUMN IF NOT EXISTS `total_orders` int DEFAULT 0 AFTER `total_spent`;
ALTER TABLE users ADD COLUMN IF NOT EXISTS `last_login_at` datetime DEFAULT NULL AFTER `total_orders`;
ALTER TABLE users ADD COLUMN IF NOT EXISTS `last_login_ip` varchar(64) DEFAULT NULL AFTER `last_login_at`;

-- ============================================
-- 4. 确保 wallets 表存在
-- ============================================
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

-- ============================================
-- 5. 确保 wallet_transactions 表存在
-- ============================================
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

-- ============================================
-- 6. 确保 tickets 表存在
-- ============================================
CREATE TABLE IF NOT EXISTS `tickets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text,
  `type` varchar(20) NOT NULL DEFAULT 'question',
  `priority` varchar(20) NOT NULL DEFAULT 'medium',
  `status` varchar(20) NOT NULL DEFAULT 'open',
  `assigned_to` bigint unsigned DEFAULT NULL,
  `instance_id` bigint unsigned DEFAULT NULL,
  `tags` varchar(255) DEFAULT NULL,
  `resolution_notes` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `closed_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `user_id` (`user_id`),
  INDEX `status` (`status`),
  INDEX `assigned_to` (`assigned_to`),
  INDEX `instance_id` (`instance_id`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 7. 确保 ticket_replies 表存在
-- ============================================
CREATE TABLE IF NOT EXISTS `ticket_replies` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `content` text NOT NULL,
  `is_admin` tinyint(1) DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `ticket_id` (`ticket_id`),
  INDEX `user_id` (`user_id`),
  INDEX `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 8. 确保 commissions 表存在
-- ============================================
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

-- ============================================
-- 9. 确保 agent_sub_users 表存在
-- ============================================
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

-- ============================================
-- 10. 为现有用户创建钱包
-- ============================================
INSERT INTO `wallets` (`user_id`, `balance`)
SELECT `id`, 0 FROM `users` WHERE `id` NOT IN (SELECT `user_id` FROM `wallets`)
ON DUPLICATE KEY UPDATE `balance` = COALESCE(`wallets`.`balance`, 0);

-- ============================================
-- 11. 确保 system_configs 有必要的配置
-- ============================================
INSERT INTO `system_configs` (`key`, `value`, `description`, `created_at`, `updated_at`) VALUES
('enable_agent', 'true', '是否开启代理商功能', NOW(), NOW()),
('enable_email_verification', 'false', '是否开启邮箱验证（注册后需验证邮箱）', NOW(), NOW()),
('enable_real_name', 'false', '是否开启实名认证', NOW(), NOW()),
('require_real_name', 'false', '是否强制实名认证后才能使用服务', NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

-- ============================================
-- 12. 确保 site_configs 有必要的配置
-- ============================================
INSERT INTO `site_configs` (`key`, `value`, `type`, `group`, `description`, `created_at`, `updated_at`) VALUES
('site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW()),
('footer_text', '© 2025 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

-- ============================================
-- 13. 确保 domain_configs 表存在
-- ============================================
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

INSERT INTO `domain_configs` (`max_domains_per_user`, `max_domains_per_agent_user`, `default_ttl`, `auto_ssl`, `allowed_suffixes`, `dns_type`, `dns_config_path`, `nginx_config_path`, `created_at`, `updated_at`)
SELECT 3, 5, 300, 0, '', 'dnsmasq', '/etc/dnsmasq.d/oneclickvirt-hosts.conf', '/etc/nginx/conf.d/oneclickvirt-domains', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM `domain_configs` LIMIT 1);

-- ============================================
-- 14. 确保 kyc_verifications 表存在
-- ============================================
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

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 自动修复完成
-- ============================================
