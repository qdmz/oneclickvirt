-- 完整的数据库初始化脚本
-- 包含所有表结构创建和默认数据导入

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 1. 创建 users 表（先创建，因为其他表有外键依赖）
CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `username` varchar(64) NOT NULL,
  `email` varchar(128) NOT NULL,
  `password` varchar(255) NOT NULL,
  `user_type` varchar(32) DEFAULT 'user',
  `level` int DEFAULT '1',
  `status` int DEFAULT '1',
  `nickname` varchar(64) DEFAULT NULL,
  `phone` varchar(32) DEFAULT NULL,
  `email_verified` tinyint(1) DEFAULT '0',
  `real_name_verified` tinyint(1) DEFAULT '0',
  `level_expire_at` datetime DEFAULT NULL,
  `max_instances` int DEFAULT '1',
  `max_cpu` int DEFAULT '1',
  `max_memory` int DEFAULT '512',
  `max_disk` int DEFAULT '10240',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. 创建 roles 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. 创建 user_roles 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. 创建 announcements 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. 创建 products 表
CREATE TABLE IF NOT EXISTS `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text,
  `price` int DEFAULT '0',
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
  `is_enabled` int DEFAULT '1',
  `cpu` int DEFAULT '0',
  `memory` int DEFAULT '0',
  `disk` int DEFAULT '0',
  `bandwidth` int DEFAULT '0',
  `traffic` int DEFAULT '0',
  `period` int DEFAULT '0',
  `allow_repeat` int DEFAULT '0',
  `features` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `status` (`status`),
  KEY `type` (`type`),
  KEY `is_featured` (`is_featured`),
  KEY `is_recommended` (`is_recommended`),
  KEY `sort_order` (`sort_order`),
  KEY `is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 6. 创建 system_configs 表
CREATE TABLE IF NOT EXISTS `system_configs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `key` varchar(128) NOT NULL,
  `value` text,
  `description` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 7. 创建 site_configs 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 8. 创建 domain_configs 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 9. 创建 domains 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 10. 创建 agents 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 11. 创建 sub_user_relations 表
CREATE TABLE IF NOT EXISTS `sub_user_relations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `agent_id` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`),
  KEY `agent_id` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 12. 创建 commissions 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 13. 创建 kyc_records 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 14. 创建 orders 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 15. 创建 user_wallets 表
CREATE TABLE IF NOT EXISTS `user_wallets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `balance` decimal(10,2) DEFAULT '0.00',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 16. 创建 wallet_transactions 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 17. 创建 providers 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 18. 创建 instances 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 19. 创建 ports 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 20. 创建 user_permissions 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 21. 创建 system_images 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 22. 创建 audit_logs 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 23. 创建 captchas 表
CREATE TABLE IF NOT EXISTS `captchas` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(64) NOT NULL,
  `answer` varchar(64) NOT NULL,
  `expires_at` datetime(3) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 24. 创建 password_resets 表
CREATE TABLE IF NOT EXISTS `password_resets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` datetime(3) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `email` (`email`),
  KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 25. 创建 verify_codes 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 26. 创建 tasks 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 27. 创建 pending_deletions 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 28. 创建 configuration_tasks 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 29. 创建 performance_metrics 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 30. 创建 pmacct_monitors 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 31. 创建 pmacct_traffic_records 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 32. 创建 instance_traffic_histories 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 33. 创建 provider_traffic_histories 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 34. 创建 user_traffic_histories 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 35. 创建 traffic_monitor_tasks 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 36. 创建 o_auth2_providers 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 37. 创建 resource_reservations 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 38. 创建 product_purchases 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 39. 创建 redemption_codes 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 40. 创建 redemption_code_usages 表
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 41. 创建 invite_code_usages 表
CREATE TABLE IF NOT EXISTS `invite_code_usages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code_id` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `code_id` (`code_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 42. 创建 jwt_secrets 表
CREATE TABLE IF NOT EXISTS `jwt_secrets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `secret_key` varchar(512) NOT NULL COMMENT 'JWT签名密钥',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `secret_key` (`secret_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 43. 创建 invite_codes 表
CREATE TABLE IF NOT EXISTS `invite_codes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(32) NOT NULL,
  `creator_id` int NOT NULL,
  `creator_name` varchar(50) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `max_uses` int NOT NULL DEFAULT 1,
  `used_count` int NOT NULL DEFAULT 0,
  `expires_at` datetime(3) DEFAULT NULL,
  `status` int NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  KEY `creator_id` (`creator_id`),
  KEY `expires_at` (`expires_at`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;

-- 导入默认数据

-- 1. 导入角色数据
INSERT IGNORE INTO `roles` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`) VALUES
(1, '管理员', 'admin', '系统管理员角色', 1, NOW(), NOW()),
(2, '普通用户', 'user', '普通用户角色', 1, NOW(), NOW());

-- 2. 导入用户数据（密码：password）
INSERT IGNORE INTO `users` (`id`, `uuid`, `username`, `email`, `password`, `user_type`, `level`, `status`, `nickname`, `created_at`, `updated_at`, `max_instances`, `max_cpu`, `max_memory`, `max_disk`) VALUES
(1, 'user-1', 'admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', 5, 1, '管理员', NOW(), NOW(), 10, 8, 8192, 102400),
(2, 'user-2', 'user', 'user@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', 1, 1, '测试用户', NOW(), NOW(), 1, 1, 512, 10240);

-- 3. 导入用户角色关联
INSERT IGNORE INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`) VALUES
(1, 1, 1, NOW(), NOW()),
(2, 2, 2, NOW(), NOW());

-- 4. 导入产品数据
INSERT IGNORE INTO `products` (`id`, `uuid`, `name`, `description`, `price`, `duration`, `level`, `max_instances`, `max_cpu`, `max_memory`, `max_disk`, `max_bandwidth`, `max_traffic`, `traffic_limited`, `status`, `type`, `is_featured`, `is_recommended`, `sort_order`, `cpu_limit`, `memory_limit`, `disk_limit`, `bandwidth_limit`, `traffic_limit`, `instance_limit`, `billing_cycle`, `is_enabled`, `cpu`, `memory`, `disk`, `bandwidth`, `traffic`, `period`, `allow_repeat`, `features`, `created_at`, `updated_at`) VALUES
(1, 'product-1', '入门套餐', '适合个人用户的基础套餐，包含基本的虚拟化功能', 0, 30, 1, 1, 1, 512, 10240, 100, 0, 0, 1, 'standard', 0, 1, 1, 1, 512, 10240, 100, 0, 1, 'monthly', 1, 1, 512, 10240, 100, 0, 30, 1, '{}', NOW(), NOW()),
(2, 'product-2', '标准套餐', '适合小型团队的标准套餐，包含更多资源', 990, 30, 2, 3, 2, 1024, 20480, 200, 0, 0, 1, 'standard', 0, 0, 2, 2, 1024, 20480, 200, 0, 3, 'monthly', 1, 2, 1024, 20480, 200, 0, 30, 1, '{}', NOW(), NOW()),
(3, 'product-3', '专业套餐', '适合中型团队的专业套餐，包含完整功能', 2990, 30, 3, 5, 4, 2048, 40960, 500, 0, 0, 1, 'standard', 1, 1, 3, 4, 2048, 40960, 500, 0, 5, 'monthly', 1, 4, 2048, 40960, 500, 0, 30, 1, '{}', NOW(), NOW()),
(4, 'product-4', '企业套餐', '适合大型团队的企业套餐，包含无限资源', 9990, 30, 4, 10, 8, 4096, 102400, 1000, 0, 0, 1, 'standard', 1, 0, 4, 8, 4096, 102400, 1000, 0, 10, 'monthly', 1, 8, 4096, 102400, 1000, 0, 30, 1, '{}', NOW(), NOW());

-- 5. 导入系统配置
INSERT IGNORE INTO `system_configs` (`id`, `key`, `value`, `description`, `created_at`, `updated_at`) VALUES
(1, 'site_name', 'OneClickVirt', '网站名称', NOW(), NOW()),
(2, 'site_description', '虚拟化管理平台', '网站描述', NOW(), NOW()),
(3, 'site_keywords', '虚拟化,Docker,LXD,Incus,Proxmox', '网站关键词', NOW(), NOW()),
(4, 'enable_registration', 'true', '是否开启注册', NOW(), NOW()),
(5, 'enable_email_verify', 'false', '是否开启邮箱验证', NOW(), NOW()),
(6, 'default_user_level', '1', '默认用户等级', NOW(), NOW()),
(7, 'max_instances_per_user', '10', '每个用户最大实例数', NOW(), NOW()),
(8, 'default_instance_expiry_days', '30', '默认实例过期天数', NOW(), NOW()),
(9, 'enable_email_verification', 'false', '是否开启邮箱验证（注册后需验证邮箱）', NOW(), NOW()),
(10, 'email_activation_expire_hours', '24', '邮箱激活链接过期时间（小时）', NOW(), NOW()),
(11, 'enable_real_name', 'false', '是否开启实名认证', NOW(), NOW()),
(12, 'require_real_name', 'false', '是否强制实名认证后才能使用服务', NOW(), NOW()),
(13, 'enable_agent', 'true', '是否开启代理商功能', NOW(), NOW());

-- 6. 导入站点配置
INSERT IGNORE INTO `site_configs` (`id`, `key`, `value`, `type`, `group`, `description`, `created_at`, `updated_at`) VALUES
(1, 'site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW()),
(2, 'site_icon_url', '/favicon.ico', 'string', 'basic', '网站图标URL', NOW(), NOW()),
(3, 'site_logo_url', '/logo.png', 'string', 'basic', '网站Logo URL', NOW(), NOW()),
(4, 'footer_text', '© 2025 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW()),
(5, 'icp_number', '', 'string', 'basic', 'ICP备案号', NOW(), NOW()),
(6, 'police_number', '', 'string', 'basic', '公安备案号', NOW(), NOW());

-- 7. 导入域名配置
INSERT IGNORE INTO `domain_configs` (`id`, `max_domains_per_user`, `max_domains_per_agent_user`, `default_ttl`, `auto_ssl`, `allowed_suffixes`, `dns_type`, `dns_config_path`, `nginx_config_path`, `created_at`, `updated_at`) VALUES
(1, 3, 5, 300, 0, '', 'dnsmasq', '/etc/dnsmasq.d/oneclickvirt-hosts.conf', '/etc/nginx/conf.d/oneclickvirt-domains', NOW(), NOW());

-- 8. 导入公告数据
INSERT IGNORE INTO `announcements` (`id`, `uuid`, `title`, `content`, `content_html`, `type`, `status`, `priority`, `is_sticky`, `sort_order`, `created_at`, `updated_at`) VALUES
(1, 'announcement-1', '欢迎使用虚拟化管理平台', '欢迎使用虚拟化管理平台，支持Docker、LXD、Incus、Proxmox VE等多种虚拟化技术。本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。', '<p>欢迎使用虚拟化管理平台，支持<strong>Docker</strong>、<strong>LXD</strong>、<strong>Incus</strong>、<strong>Proxmox VE</strong>等多种虚拟化技术。</p><p>本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。</p>', 'homepage', 1, 10, 1, 1, NOW(), NOW()),
(2, 'announcement-2', '系统维护通知', '为了提供更好的服务质量，我们会定期进行系统维护。维护期间可能会影响部分功能的使用，请您谅解。', '<p>为了提供更好的服务质量，我们会定期进行系统维护。</p>', 'topbar', 1, 5, 0, 2, NOW(), NOW()),
(3, 'announcement-3', '新手使用指南', '如果您是第一次使用本平台，建议先阅读使用文档。您可以在右上角的帮助菜单中找到详细的操作指南。', '<p>如果您是第一次使用本平台，建议先阅读使用文档。</p>', 'homepage', 1, 8, 0, 3, NOW(), NOW());

-- 9. 导入邀请码数据
INSERT IGNORE INTO `invite_codes` (`id`, `code`, `creator_id`, `creator_name`, `description`, `max_uses`, `used_count`, `expires_at`, `status`, `created_at`, `updated_at`) VALUES
(1, 'SC0Q19BW', 1, 'admin', '', 1, 0, NULL, 1, NOW(), NOW());

-- 10. 导入JWT密钥数据
INSERT IGNORE INTO `jwt_secrets` (`id`, `secret_key`, `created_at`, `updated_at`) VALUES
(1, 'b64dca17bf31d0e725285cccf00a6911a43b0e2c8d8d26ed458cdbf16e6a14b5', NOW(), NOW());

-- 11. 导入用户钱包数据
INSERT IGNORE INTO `user_wallets` (`id`, `user_id`, `balance`, `created_at`, `updated_at`) VALUES
(1, 1, 0.00, NOW(), NOW()),
(2, 2, 0.00, NOW(), NOW());

-- 12. 导入用户权限数据
INSERT IGNORE INTO `user_permissions` (`id`, `user_id`, `user_types`, `level`, `is_active`, `created_at`, `updated_at`) VALUES
(1, 1, 'admin', 5, 1, NOW(), NOW()),
(2, 2, 'user', 1, 1, NOW(), NOW());

-- 13. 导入系统镜像数据
INSERT IGNORE INTO `system_images` (`id`, `uuid`, `name`, `description`, `type`, `status`, `provider`, `os_type`, `version`, `is_default`, `sort_order`, `created_at`, `updated_at`) VALUES
(1, 'image-1', 'Ubuntu 22.04 LTS', 'Ubuntu 22.04 LTS 长期支持版本', 'container', 'active', 'docker', 'linux', '22.04', 1, 1, NOW(), NOW()),
(2, 'image-2', 'CentOS 7', 'CentOS 7 稳定版本', 'container', 'active', 'docker', 'linux', '7', 0, 2, NOW(), NOW()),
(3, 'image-3', 'Debian 11', 'Debian 11 稳定版本', 'container', 'active', 'docker', 'linux', '11', 0, 3, NOW(), NOW()),
(4, 'image-4', 'Alpine 3.18', 'Alpine 3.18 轻量级版本', 'container', 'active', 'docker', 'linux', '3.18', 0, 4, NOW(), NOW());

-- 完成初始化
SELECT '数据库初始化完成' AS message;
