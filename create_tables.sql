-- 创建所有必要的表结构

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 1. 创建 roles 表
-- ============================================
CREATE TABLE IF NOT EXISTS `roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `code` varchar(64) NOT NULL,
  `description` text,
  `status` int DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 2. 创建 user_roles 表
-- ============================================
CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `role_id` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `user_roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `user_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 3. 创建 announcements 表
-- ============================================
CREATE TABLE IF NOT EXISTS `announcements` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `title` varchar(255) NOT NULL,
  `content` text,
  `content_html` text,
  `type` varchar(32) DEFAULT NULL,
  `status` int DEFAULT '1',
  `priority` int DEFAULT '0',
  `is_sticky` int DEFAULT '0',
  `start_at` datetime(3) DEFAULT NULL,
  `end_at` datetime(3) DEFAULT NULL,
  `sort_order` int DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `type` (`type`),
  KEY `status` (`status`),
  KEY `is_sticky` (`is_sticky`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 4. 创建 products 表
-- ============================================
CREATE TABLE IF NOT EXISTS `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text,
  `price` decimal(10,2) DEFAULT '0.00',
  `duration` int DEFAULT '30',
  `level` int DEFAULT '1',
  `max_instances` int DEFAULT '1',
  `max_cpu` int DEFAULT '1',
  `max_memory` int DEFAULT '512',
  `max_disk` int DEFAULT '10240',
  `max_bandwidth` int DEFAULT '100',
  `max_traffic` int DEFAULT '0',
  `traffic_limited` int DEFAULT '0',
  `status` int DEFAULT '1',
  `type` varchar(32) DEFAULT 'standard',
  `is_featured` int DEFAULT '0',
  `is_recommended` int DEFAULT '0',
  `sort_order` int DEFAULT '0',
  `icon` varchar(255) DEFAULT NULL,
  `cpu_limit` int DEFAULT '1',
  `memory_limit` int DEFAULT '512',
  `disk_limit` int DEFAULT '10240',
  `bandwidth_limit` int DEFAULT '100',
  `traffic_limit` int DEFAULT '0',
  `instance_limit` int DEFAULT '1',
  `enable_auto_renewal` int DEFAULT '0',
  `auto_renewal_discount` decimal(5,2) DEFAULT '0.00',
  `billing_cycle` varchar(32) DEFAULT 'monthly',
  `setup_fee` decimal(10,2) DEFAULT '0.00',
  `recurring_fee` decimal(10,2) DEFAULT '0.00',
  `trial_duration` int DEFAULT '0',
  `trial_enabled` int DEFAULT '0',
  `refund_policy` text,
  `terms_of_service` text,
  `recommended_for` text,
  `target_audience` text,
  `tags` text,
  `metadata` text,
  `stock` int DEFAULT '-1',
  `sold_count` int DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `status` (`status`),
  KEY `type` (`type`),
  KEY `is_featured` (`is_featured`),
  KEY `is_recommended` (`is_recommended`),
  KEY `sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 5. 创建 system_configs 表
-- ============================================
CREATE TABLE IF NOT EXISTS `system_configs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `key` varchar(128) NOT NULL,
  `value` text,
  `description` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 6. 创建 site_configs 表
-- ============================================
CREATE TABLE IF NOT EXISTS `site_configs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `key` varchar(128) NOT NULL,
  `value` text,
  `type` varchar(32) DEFAULT 'string',
  `group` varchar(32) DEFAULT 'basic',
  `description` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 7. 创建 domain_configs 表
-- ============================================
CREATE TABLE IF NOT EXISTS `domain_configs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `max_domains_per_user` int DEFAULT '3',
  `max_domains_per_agent_user` int DEFAULT '5',
  `default_ttl` int DEFAULT '300',
  `auto_ssl` int DEFAULT '0',
  `allowed_suffixes` text,
  `dns_type` varchar(32) DEFAULT 'dnsmasq',
  `dns_config_path` varchar(255) DEFAULT '/etc/dnsmasq.d/oneclickvirt-hosts.conf',
  `nginx_config_path` varchar(255) DEFAULT '/etc/nginx/conf.d/oneclickvirt-domains',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 8. 创建 domains 表
-- ============================================
CREATE TABLE IF NOT EXISTS `domains` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `user_id` int NOT NULL,
  `domain` varchar(255) NOT NULL,
  `instance_id` int DEFAULT NULL,
  `status` int DEFAULT '1',
  `ssl_status` int DEFAULT '0',
  `ssl_expiry` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `domain` (`domain`),
  KEY `user_id` (`user_id`),
  KEY `instance_id` (`instance_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 9. 创建 agents 表
-- ============================================
CREATE TABLE IF NOT EXISTS `agents` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `commission_rate` decimal(5,2) DEFAULT '0.00',
  `total_commission` decimal(10,2) DEFAULT '0.00',
  `withdrawn_commission` decimal(10,2) DEFAULT '0.00',
  `status` int DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 10. 创建 sub_user_relations 表
-- ============================================
CREATE TABLE IF NOT EXISTS `sub_user_relations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `agent_id` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`),
  KEY `agent_id` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 11. 创建 commissions 表
-- ============================================
CREATE TABLE IF NOT EXISTS `commissions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `agent_id` int NOT NULL,
  `user_id` int NOT NULL,
  `order_id` int DEFAULT NULL,
  `amount` decimal(10,2) DEFAULT '0.00',
  `status` int DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `agent_id` (`agent_id`),
  KEY `user_id` (`user_id`),
  KEY `order_id` (`order_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 12. 创建 kyc_records 表
-- ============================================
CREATE TABLE IF NOT EXISTS `kyc_records` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `real_name` varchar(64) DEFAULT NULL,
  `id_card` varchar(32) DEFAULT NULL,
  `id_card_front` varchar(255) DEFAULT NULL,
  `id_card_back` varchar(255) DEFAULT NULL,
  `status` int DEFAULT '0',
  `admin_id` int DEFAULT NULL,
  `admin_notes` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 13. 创建 orders 表
-- ============================================
CREATE TABLE IF NOT EXISTS `orders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `user_id` int NOT NULL,
  `product_id` int NOT NULL,
  `order_no` varchar(64) NOT NULL,
  `amount` decimal(10,2) DEFAULT '0.00',
  `status` varchar(32) DEFAULT 'pending',
  `payment_method` varchar(32) DEFAULT NULL,
  `payment_time` datetime(3) DEFAULT NULL,
  `refund_time` datetime(3) DEFAULT NULL,
  `refund_reason` text,
  `refund_amount` decimal(10,2) DEFAULT '0.00',
  `expires_at` datetime(3) DEFAULT NULL,
  `auto_renew` int DEFAULT '0',
  `renewal_price` decimal(10,2) DEFAULT '0.00',
  `renewal_discount` decimal(5,2) DEFAULT '0.00',
  `original_order_id` int DEFAULT NULL,
  `affiliate_id` int DEFAULT NULL,
  `commission_amount` decimal(10,2) DEFAULT '0.00',
  `coupon_code` varchar(64) DEFAULT NULL,
  `discount_amount` decimal(10,2) DEFAULT '0.00',
  `actual_amount` decimal(10,2) DEFAULT '0.00',
  `billing_cycle` varchar(32) DEFAULT 'monthly',
  `setup_fee` decimal(10,2) DEFAULT '0.00',
  `recurring_fee` decimal(10,2) DEFAULT '0.00',
  `trial_used` int DEFAULT '0',
  `trial_end_at` datetime(3) DEFAULT NULL,
  `metadata` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `order_no` (`order_no`),
  KEY `user_id` (`user_id`),
  KEY `product_id` (`product_id`),
  KEY `status` (`status`),
  KEY `payment_method` (`payment_method`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 14. 创建 user_wallets 表
-- ============================================
CREATE TABLE IF NOT EXISTS `user_wallets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `balance` decimal(10,2) DEFAULT '0.00',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 15. 创建 wallet_transactions 表
-- ============================================
CREATE TABLE IF NOT EXISTS `wallet_transactions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `type` varchar(32) DEFAULT NULL,
  `amount` decimal(10,2) DEFAULT '0.00',
  `balance` decimal(10,2) DEFAULT '0.00',
  `description` text,
  `order_id` int DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `type` (`type`),
  KEY `order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 16. 创建 providers 表
-- ============================================
CREATE TABLE IF NOT EXISTS `providers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(64) NOT NULL,
  `type` varchar(32) NOT NULL,
  `endpoint` varchar(255) DEFAULT NULL,
  `port_ip` varchar(255) DEFAULT NULL,
  `ssh_port` bigint DEFAULT '22',
  `username` varchar(128) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `ssh_key` text,
  `token` varchar(255) DEFAULT NULL,
  `config` text,
  `status` varchar(16) DEFAULT 'active',
  `region` varchar(64) DEFAULT NULL,
  `country` varchar(64) DEFAULT NULL,
  `country_code` varchar(8) DEFAULT NULL,
  `city` varchar(64) DEFAULT NULL,
  `version` varchar(32) DEFAULT '',
  `container_enabled` tinyint(1) DEFAULT '1',
  `virtual_machine_enabled` tinyint(1) DEFAULT '0',
  `supported_types` varchar(128) DEFAULT NULL,
  `allow_claim` tinyint(1) DEFAULT '1',
  `ipv4_port_mapping_method` varchar(16) DEFAULT 'device_proxy',
  `ipv6_port_mapping_method` varchar(16) DEFAULT 'device_proxy',
  `used_quota` bigint DEFAULT '0',
  `total_quota` bigint DEFAULT '0',
  `architecture` varchar(16) DEFAULT 'amd64',
  `expires_at` datetime(3) DEFAULT NULL,
  `is_frozen` tinyint(1) DEFAULT '0',
  `storage_pool` varchar(64) DEFAULT 'local',
  `storage_pool_path` varchar(255) DEFAULT '',
  `cert_path` varchar(512) DEFAULT NULL,
  `key_path` varchar(512) DEFAULT NULL,
  `ca_cert_path` varchar(512) DEFAULT NULL,
  `cert_fingerprint` varchar(128) DEFAULT NULL,
  `trusted_fingerprint` varchar(128) DEFAULT NULL,
  `api_status` varchar(16) DEFAULT 'unknown',
  `ssh_status` varchar(16) DEFAULT 'unknown',
  `last_api_check` datetime(3) DEFAULT NULL,
  `last_ssh_check` datetime(3) DEFAULT NULL,
  `auth_config` text,
  `config_version` bigint DEFAULT '0',
  `auto_configured` tinyint(1) DEFAULT '0',
  `last_config_update` datetime(3) DEFAULT NULL,
  `config_backup_path` varchar(512) DEFAULT NULL,
  `cert_content` text,
  `key_content` text,
  `token_content` text,
  `node_cpu_cores` bigint DEFAULT '0',
  `node_memory_total` bigint DEFAULT '0',
  `node_disk_total` bigint DEFAULT '0',
  `allow_concurrent_tasks` tinyint(1) DEFAULT '0',
  `max_concurrent_tasks` bigint DEFAULT '1',
  `ssh_connect_timeout` bigint DEFAULT '30',
  `ssh_execute_timeout` bigint DEFAULT '300',
  `task_poll_interval` bigint DEFAULT '60',
  `enable_task_polling` tinyint(1) DEFAULT '1',
  `execution_rule` varchar(16) DEFAULT 'auto',
  `max_container_instances` bigint DEFAULT '0',
  `max_vm_instances` bigint DEFAULT '0',
  `container_limit_cpu` tinyint(1) DEFAULT '0',
  `container_limit_memory` tinyint(1) DEFAULT '0',
  `container_limit_disk` tinyint(1) DEFAULT '1',
  `vm_limit_cpu` tinyint(1) DEFAULT '1',
  `vm_limit_memory` tinyint(1) DEFAULT '1',
  `vm_limit_disk` tinyint(1) DEFAULT '1',
  `default_port_count` bigint DEFAULT '10',
  `port_range_start` bigint DEFAULT '10000',
  `port_range_end` bigint DEFAULT '65535',
  `next_available_port` bigint DEFAULT '10000',
  `network_type` varchar(32) NOT NULL DEFAULT 'nat_ipv4',
  `default_inbound_bandwidth` bigint DEFAULT '300',
  `default_outbound_bandwidth` bigint DEFAULT '300',
  `max_inbound_bandwidth` bigint DEFAULT '1000',
  `max_outbound_bandwidth` bigint DEFAULT '1000',
  `enable_traffic_control` tinyint(1) DEFAULT '0',
  `max_traffic` bigint DEFAULT '1048576',
  `traffic_limited` tinyint(1) DEFAULT '0',
  `traffic_reset_at` datetime(3) DEFAULT NULL,
  `traffic_count_mode` varchar(16) DEFAULT 'both',
  `traffic_multiplier` double DEFAULT '1',
  `traffic_stats_mode` varchar(16) DEFAULT 'light',
  `traffic_collect_interval` bigint DEFAULT '300',
  `traffic_collect_batch_size` bigint DEFAULT '10',
  `traffic_limit_check_interval` bigint DEFAULT '600',
  `traffic_limit_check_batch_size` bigint DEFAULT '10',
  `traffic_auto_reset_interval` bigint DEFAULT '1800',
  `traffic_auto_reset_batch_size` bigint DEFAULT '10',
  `used_cpu_cores` bigint DEFAULT '0',
  `used_memory` bigint DEFAULT '0',
  `used_disk` bigint DEFAULT '0',
  `container_count` bigint DEFAULT '0',
  `vm_count` bigint DEFAULT '0',
  `resource_synced` tinyint(1) DEFAULT '0',
  `resource_synced_at` datetime(3) DEFAULT NULL,
  `count_cache_expiry` datetime(3) DEFAULT NULL,
  `available_cpu_cores` bigint DEFAULT '0',
  `available_memory` bigint DEFAULT '0',
  `used_instances` bigint DEFAULT '0',
  `level_limits` text,
  `host_name` varchar(128) DEFAULT NULL,
  `container_privileged` tinyint(1) DEFAULT '0',
  `container_allow_nesting` tinyint(1) DEFAULT '0',
  `container_enable_lxcfs` tinyint(1) DEFAULT '1',
  `container_cpu_allowance` varchar(16) DEFAULT '100%',
  `container_memory_swap` tinyint(1) DEFAULT '1',
  `container_max_processes` bigint DEFAULT '0',
  `container_disk_io_limit` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `name` (`name`),
  KEY `type` (`type`),
  KEY `status` (`status`),
  KEY `region` (`region`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 17. 创建 instances 表
-- ============================================
CREATE TABLE IF NOT EXISTS `instances` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `provider` varchar(32) NOT NULL,
  `provider_id` bigint unsigned NOT NULL,
  `status` varchar(32) DEFAULT NULL,
  `image` varchar(128) DEFAULT NULL,
  `instance_type` varchar(16) DEFAULT 'container',
  `cpu` bigint DEFAULT '1',
  `memory` bigint DEFAULT '512',
  `disk` bigint DEFAULT '10240',
  `bandwidth` bigint DEFAULT '10',
  `network` varchar(64) DEFAULT NULL,
  `private_ip` varchar(64) DEFAULT NULL,
  `public_ip` varchar(64) DEFAULT NULL,
  `ipv6_address` varchar(128) DEFAULT NULL,
  `public_ipv6` varchar(128) DEFAULT NULL,
  `ssh_port` bigint DEFAULT '22',
  `port_range_start` bigint DEFAULT NULL,
  `port_range_end` bigint DEFAULT NULL,
  `username` varchar(64) DEFAULT NULL,
  `password` varchar(128) DEFAULT NULL,
  `os_type` varchar(64) DEFAULT NULL,
  `region` varchar(64) DEFAULT NULL,
  `max_traffic` bigint DEFAULT '0',
  `traffic_limited` tinyint(1) DEFAULT '0',
  `traffic_limit_reason` varchar(16) DEFAULT '',
  `pmacct_interface_v4` varchar(32) DEFAULT NULL,
  `pmacct_interface_v6` varchar(32) DEFAULT NULL,
  `expired_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `name` (`name`),
  KEY `provider` (`provider`),
  KEY `provider_id` (`provider_id`),
  KEY `status` (`status`),
  KEY `instance_type` (`instance_type`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 18. 创建 ports 表
-- ============================================
CREATE TABLE IF NOT EXISTS `ports` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `instance_id` bigint unsigned NOT NULL,
  `provider_id` bigint unsigned NOT NULL,
  `protocol` varchar(16) DEFAULT 'tcp',
  `private_port` bigint DEFAULT NULL,
  `public_port` bigint DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `status` varchar(16) DEFAULT 'active',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `instance_id` (`instance_id`),
  KEY `provider_id` (`provider_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 19. 创建 user_permissions 表
-- ============================================
CREATE TABLE IF NOT EXISTS `user_permissions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` int NOT NULL,
  `user_types` varchar(255) DEFAULT '',
  `level` int DEFAULT 1,
  `is_active` tinyint(1) DEFAULT 1,
  `remark` varchar(255) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 20. 创建 system_images 表
-- ============================================
CREATE TABLE IF NOT EXISTS `system_images` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `name` varchar(128) NOT NULL,
  `description` text,
  `type` varchar(32) DEFAULT NULL,
  `status` varchar(32) DEFAULT 'active',
  `provider` varchar(32) DEFAULT NULL,
  `os_type` varchar(64) DEFAULT NULL,
  `version` varchar(64) DEFAULT NULL,
  `size` bigint DEFAULT '0',
  `checksum` varchar(256) DEFAULT NULL,
  `url` varchar(512) DEFAULT NULL,
  `is_default` int DEFAULT '0',
  `sort_order` int DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `name` (`name`),
  KEY `type` (`type`),
  KEY `status` (`status`),
  KEY `provider` (`provider`),
  KEY `is_default` (`is_default`),
  KEY `sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 21. 创建 audit_logs 表
-- ============================================
CREATE TABLE IF NOT EXISTS `audit_logs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int DEFAULT NULL,
  `username` varchar(64) DEFAULT NULL,
  `action` varchar(128) DEFAULT NULL,
  `resource_type` varchar(64) DEFAULT NULL,
  `resource_id` int DEFAULT NULL,
  `ip_address` varchar(64) DEFAULT NULL,
  `user_agent` text,
  `details` text,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `action` (`action`),
  KEY `resource_type` (`resource_type`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 22. 创建 captchas 表
-- ============================================
CREATE TABLE IF NOT EXISTS `captchas` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(64) NOT NULL,
  `answer` varchar(64) NOT NULL,
  `expires_at` datetime(3) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 23. 创建 password_resets 表
-- ============================================
CREATE TABLE IF NOT EXISTS `password_resets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` datetime(3) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `email` (`email`),
  KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 24. 创建 verify_codes 表
-- ============================================
CREATE TABLE IF NOT EXISTS `verify_codes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL,
  `code` varchar(16) NOT NULL,
  `type` varchar(32) DEFAULT 'email',
  `expires_at` datetime(3) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `email` (`email`),
  KEY `type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 25. 创建 tasks 表
-- ============================================
CREATE TABLE IF NOT EXISTS `tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `provider_id` bigint unsigned DEFAULT NULL,
  `instance_id` bigint unsigned DEFAULT NULL,
  `type` varchar(32) NOT NULL,
  `status` varchar(32) DEFAULT 'pending',
  `priority` int DEFAULT '0',
  `progress` int DEFAULT '0',
  `message` text,
  `error` text,
  `data` text,
  `metadata` text,
  `started_at` datetime(3) DEFAULT NULL,
  `completed_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `user_id` (`user_id`),
  KEY `provider_id` (`provider_id`),
  KEY `instance_id` (`instance_id`),
  KEY `type` (`type`),
  KEY `status` (`status`),
  KEY `priority` (`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 26. 创建 pending_deletions 表
-- ============================================
CREATE TABLE IF NOT EXISTS `pending_deletions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int DEFAULT NULL,
  `instance_id` int DEFAULT NULL,
  `type` varchar(32) DEFAULT NULL,
  `status` varchar(32) DEFAULT 'pending',
  `scheduled_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `instance_id` (`instance_id`),
  KEY `type` (`type`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 27. 创建 configuration_tasks 表
-- ============================================
CREATE TABLE IF NOT EXISTS `configuration_tasks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `provider_id` int NOT NULL,
  `type` varchar(32) DEFAULT NULL,
  `status` varchar(32) DEFAULT 'pending',
  `result` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `provider_id` (`provider_id`),
  KEY `type` (`type`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 28. 创建 performance_metrics 表
-- ============================================
CREATE TABLE IF NOT EXISTS `performance_metrics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `provider_id` int DEFAULT NULL,
  `cpu_usage` decimal(5,2) DEFAULT '0.00',
  `memory_usage` decimal(5,2) DEFAULT '0.00',
  `disk_usage` decimal(5,2) DEFAULT '0.00',
  `network_rx` bigint DEFAULT '0',
  `network_tx` bigint DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `provider_id` (`provider_id`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 29. 创建 pmacct_monitors 表
-- ============================================
CREATE TABLE IF NOT EXISTS `pmacct_monitors` (
  `id` int NOT NULL AUTO_INCREMENT,
  `provider_id` int NOT NULL,
  `interface` varchar(64) DEFAULT NULL,
  `status` varchar(32) DEFAULT 'active',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `provider_id` (`provider_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 30. 创建 pmacct_traffic_records 表
-- ============================================
CREATE TABLE IF NOT EXISTS `pmacct_traffic_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `provider_id` int NOT NULL,
  `instance_id` int NOT NULL,
  `ip_address` varchar(64) DEFAULT NULL,
  `bytes_rx` bigint DEFAULT '0',
  `bytes_tx` bigint DEFAULT '0',
  `packets_rx` bigint DEFAULT '0',
  `packets_tx` bigint DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `provider_id` (`provider_id`),
  KEY `instance_id` (`instance_id`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 31. 创建 instance_traffic_histories 表
-- ============================================
CREATE TABLE IF NOT EXISTS `instance_traffic_histories` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `instance_id` int NOT NULL,
  `user_id` int NOT NULL,
  `bytes_rx` bigint DEFAULT '0',
  `bytes_tx` bigint DEFAULT '0',
  `total_bytes` bigint DEFAULT '0',
  `period` varchar(32) DEFAULT 'daily',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `instance_id` (`instance_id`),
  KEY `user_id` (`user_id`),
  KEY `period` (`period`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 32. 创建 provider_traffic_histories 表
-- ============================================
CREATE TABLE IF NOT EXISTS `provider_traffic_histories` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `provider_id` int NOT NULL,
  `bytes_rx` bigint DEFAULT '0',
  `bytes_tx` bigint DEFAULT '0',
  `total_bytes` bigint DEFAULT '0',
  `period` varchar(32) DEFAULT 'daily',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `provider_id` (`provider_id`),
  KEY `period` (`period`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 33. 创建 user_traffic_histories 表
-- ============================================
CREATE TABLE IF NOT EXISTS `user_traffic_histories` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `bytes_rx` bigint DEFAULT '0',
  `bytes_tx` bigint DEFAULT '0',
  `total_bytes` bigint DEFAULT '0',
  `period` varchar(32) DEFAULT 'daily',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `period` (`period`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 34. 创建 traffic_monitor_tasks 表
-- ============================================
CREATE TABLE IF NOT EXISTS `traffic_monitor_tasks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `provider_id` int NOT NULL,
  `status` varchar(32) DEFAULT 'pending',
  `result` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `provider_id` (`provider_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 35. 创建 o_auth2_providers 表
-- ============================================
CREATE TABLE IF NOT EXISTS `o_auth2_providers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `client_id` varchar(255) NOT NULL,
  `client_secret` varchar(255) NOT NULL,
  `redirect_uri` varchar(512) NOT NULL,
  `scopes` varchar(255) DEFAULT NULL,
  `auth_url` varchar(512) NOT NULL,
  `token_url` varchar(512) NOT NULL,
  `user_info_url` varchar(512) NOT NULL,
  `enabled` int DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 36. 创建 resource_reservations 表
-- ============================================
CREATE TABLE IF NOT EXISTS `resource_reservations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `provider_id` int NOT NULL,
  `resource_type` varchar(32) DEFAULT NULL,
  `resource_id` int DEFAULT NULL,
  `quantity` int DEFAULT '1',
  `status` varchar(32) DEFAULT 'active',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `provider_id` (`provider_id`),
  KEY `resource_type` (`resource_type`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 37. 创建 product_purchases 表
-- ============================================
CREATE TABLE IF NOT EXISTS `product_purchases` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `product_id` int NOT NULL,
  `order_id` int DEFAULT NULL,
  `status` varchar(32) DEFAULT 'active',
  `start_at` datetime(3) DEFAULT NULL,
  `end_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `product_id` (`product_id`),
  KEY `order_id` (`order_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 38. 创建 redemption_codes 表
-- ============================================
CREATE TABLE IF NOT EXISTS `redemption_codes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(64) NOT NULL,
  `type` varchar(32) DEFAULT 'credit',
  `value` decimal(10,2) DEFAULT '0.00',
  `product_id` int DEFAULT NULL,
  `uses` int DEFAULT '0',
  `max_uses` int DEFAULT '1',
  `expires_at` datetime(3) DEFAULT NULL,
  `status` varchar(32) DEFAULT 'active',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  KEY `type` (`type`),
  KEY `product_id` (`product_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 39. 创建 redemption_code_usages 表
-- ============================================
CREATE TABLE IF NOT EXISTS `redemption_code_usages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code_id` int NOT NULL,
  `user_id` int NOT NULL,
  `order_id` int DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `code_id` (`code_id`),
  KEY `user_id` (`user_id`),
  KEY `order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================
-- 40. 创建 invite_code_usages 表
-- ============================================
CREATE TABLE IF NOT EXISTS `invite_code_usages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code_id` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `code_id` (`code_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;