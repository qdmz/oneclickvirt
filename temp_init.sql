-- 创建缺失的表
CREATE TABLE IF NOT EXISTS roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  code varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  description text COLLATE utf8mb4_unicode_ci,
  status bigint DEFAULT 1,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_roles_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  role_id bigint unsigned NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_user_roles_user_role (user_id,role_id),
  KEY idx_user_roles_role_id (role_id),
  CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS products (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  description text COLLATE utf8mb4_unicode_ci COMMENT '产品描述',
  level bigint NOT NULL COMMENT '产品等级',
  price bigint NOT NULL COMMENT '价格(分)',
  period bigint NOT NULL COMMENT '周期(月), 0为永久',
  cpu bigint NOT NULL COMMENT 'CPU核心数',
  memory bigint NOT NULL COMMENT '内存(MB)',
  disk bigint NOT NULL COMMENT '磁盘(MB)',
  bandwidth bigint NOT NULL COMMENT '带宽(Mbps)',
  traffic bigint NOT NULL COMMENT '流量限制(MB)',
  max_instances bigint NOT NULL COMMENT '最大实例数',
  is_enabled bigint DEFAULT 1 COMMENT '是否启用(1:启用, 0:禁用)',
  sort_order bigint DEFAULT 0 COMMENT '排序',
  features text COLLATE utf8mb4_unicode_ci COMMENT '特性(JSON格式)',
  allow_repeat bigint DEFAULT 1 COMMENT '是否允许重复购买(1:允许, 0:不允许)',
  stock bigint DEFAULT -1 COMMENT '库存(-1为无限)',
  sold_count bigint DEFAULT 0 COMMENT '已售数量',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS system_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_system_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS site_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `group` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_site_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS domain_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  max_domains_per_user bigint DEFAULT 3,
  max_domains_per_agent_user bigint DEFAULT 5,
  default_ttl bigint DEFAULT 300,
  auto_ssl bigint DEFAULT 0,
  allowed_suffixes text COLLATE utf8mb4_unicode_ci,
  dns_type varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'dnsmasq',
  dns_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/dnsmasq.d/oneclickvirt-hosts.conf',
  nginx_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/nginx/conf.d/oneclickvirt-domains',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 导入默认数据
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 1. 创建默认角色
INSERT IGNORE INTO `roles` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`) VALUES
('admin', 'admin', '系统管理员角色', 1, NOW(), NOW()),
('user', 'user', '普通用户角色', 1, NOW(), NOW());

-- 2. 创建用户角色关联
INSERT IGNORE INTO `user_roles` (`user_id`, `role_id`, `created_at`, `updated_at`) VALUES
(1, 1, NOW(), NOW());

-- 3. 创建默认产品套餐
INSERT IGNORE INTO `products` (`name`, `description`, `price`, `level`, `period`, `cpu`, `memory`, `disk`, `bandwidth`, `traffic`, `max_instances`, `is_enabled`, `sort_order`, `features`, `allow_repeat`, `stock`, `sold_count`, `created_at`, `updated_at`) VALUES
('入门套餐', '适合个人用户的基础套餐，包含基本的虚拟化功能', 0, 1, 1, 1, 512, 10240, 100, 0, 1, 1, 1, '{"cpu": "1核", "memory": "512MB", "disk": "10GB", "bandwidth": "100Mbps", "instances": "1个实例"}', 1, -1, 0, NOW(), NOW()),
('标准套餐', '适合小型团队的标准套餐，包含更多资源', 990, 2, 1, 2, 1024, 20480, 200, 0, 3, 1, 2, '{"cpu": "2核", "memory": "1GB", "disk": "20GB", "bandwidth": "200Mbps", "instances": "3个实例"}', 1, -1, 0, NOW(), NOW()),
('专业套餐', '适合中型团队的专业套餐，包含完整功能', 2990, 3, 1, 4, 2048, 40960, 500, 0, 5, 1, 3, '{"cpu": "4核", "memory": "2GB", "disk": "40GB", "bandwidth": "500Mbps", "instances": "5个实例"}', 1, -1, 0, NOW(), NOW()),
('企业套餐', '适合大型团队的企业套餐，包含无限资源', 9990, 4, 1, 8, 4096, 102400, 1000, 0, 10, 1, 4, '{"cpu": "8核", "memory": "4GB", "disk": "100GB", "bandwidth": "1000Mbps", "instances": "10个实例"}', 1, -1, 0, NOW(), NOW());

-- 4. 创建默认系统配置
INSERT IGNORE INTO `system_configs` (`key`, `value`, `description`, `created_at`, `updated_at`) VALUES
('site_name', 'OneClickVirt', '网站名称', NOW(), NOW()),
('site_description', '虚拟化管理平台', '网站描述', NOW(), NOW()),
('site_keywords', '虚拟化,Docker,LXD,Incus,Proxmox', '网站关键词', NOW(), NOW()),
('enable_registration', 'true', '是否开启注册', NOW(), NOW()),
('enable_email_verify', 'false', '是否开启邮箱验证', NOW(), NOW()),
('default_user_level', '1', '默认用户等级', NOW(), NOW()),
('max_instances_per_user', '10', '每个用户最大实例数', NOW(), NOW()),
('default_instance_expiry_days', '30', '默认实例过期天数', NOW(), NOW()),
('enable_email_verification', 'false', '是否开启邮箱验证（注册后需验证邮箱）', NOW(), NOW()),
('email_activation_expire_hours', '24', '邮箱激活链接过期时间（小时）', NOW(), NOW()),
('enable_real_name', 'false', '是否开启实名认证', NOW(), NOW()),
('require_real_name', 'false', '是否强制实名认证后才能使用服务', NOW(), NOW()),
('enable_agent', 'true', '是否开启代理商功能', NOW(), NOW());

-- 5. 创建默认站点配置
INSERT IGNORE INTO `site_configs` (`key`, `value`, `type`, `group`, `description`, `created_at`, `updated_at`) VALUES
('site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW()),
('site_icon_url', '/favicon.ico', 'string', 'basic', '网站图标URL', NOW(), NOW()),
('site_logo_url', '/logo.png', 'string', 'basic', '网站Logo URL', NOW(), NOW()),
('footer_text', '© 2025 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW()),
('icp_number', '', 'string', 'basic', 'ICP备案号', NOW(), NOW()),
('police_number', '', 'string', 'basic', '公安备案号', NOW(), NOW());

-- 6. 创建默认域名配置
INSERT IGNORE INTO `domain_configs` (`max_domains_per_user`, `max_domains_per_agent_user`, `default_ttl`, `auto_ssl`, `allowed_suffixes`, `dns_type`, `dns_config_path`, `nginx_config_path`, `created_at`, `updated_at`) VALUES
(3, 5, 300, 0, '', 'dnsmasq', '/etc/dnsmasq.d/oneclickvirt-hosts.conf', '/etc/nginx/conf.d/oneclickvirt-domains', NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;
