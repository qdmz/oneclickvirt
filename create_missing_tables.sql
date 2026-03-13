-- 创建缺失的表

-- 创建providers表
CREATE TABLE IF NOT EXISTS providers (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  type varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  host varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  port int DEFAULT 22,
  username varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  password varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'inactive',
  region varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'default',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_providers_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建instances表
CREATE TABLE IF NOT EXISTS instances (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  uuid varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  name varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  provider_id int NOT NULL,
  provider varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  instance_type varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'creating',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_instances_uuid (uuid),
  KEY idx_instances_provider_id (provider_id),
  KEY idx_instances_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建site_configs表
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

-- 创建user_permissions表
CREATE TABLE IF NOT EXISTS user_permissions (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  role_id bigint unsigned NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_user_permissions_user_role (user_id,role_id),
  KEY idx_user_permissions_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建product_purchases表
CREATE TABLE IF NOT EXISTS product_purchases (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  product_id bigint unsigned NOT NULL,
  order_id varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_product_purchases_order_id (order_id),
  KEY idx_product_purchases_user_id (user_id),
  KEY idx_product_purchases_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建oauth2_providers表
CREATE TABLE IF NOT EXISTS oauth2_providers (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  display_name varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  provider_type varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  client_id varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  client_secret varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  redirect_url varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  auth_url varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  token_url varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  user_info_url varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  user_id_field varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'id',
  username_field varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'username',
  email_field varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'email',
  avatar_field varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'avatar',
  nickname_field varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'nickname',
  trust_level_field varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '',
  max_registrations int DEFAULT -1,
  default_level int DEFAULT 1,
  sort int DEFAULT 0,
  enabled tinyint(1) DEFAULT 1,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_oauth2_providers_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建resource_reservations表
CREATE TABLE IF NOT EXISTS resource_reservations (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  provider_id int NOT NULL,
  instance_id bigint unsigned NOT NULL,
  cpu_cores int DEFAULT 1,
  memory int DEFAULT 512,
  disk int DEFAULT 10240,
  bandwidth int DEFAULT 100,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_resource_reservations_user_id (user_id),
  KEY idx_resource_reservations_provider_id (provider_id),
  KEY idx_resource_reservations_instance_id (instance_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建ports表
CREATE TABLE IF NOT EXISTS ports (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  provider_id int NOT NULL,
  instance_id bigint unsigned NOT NULL,
  protocol varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'tcp',
  external_port int NOT NULL,
  internal_port int NOT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_ports_provider_external (provider_id,external_port),
  KEY idx_ports_instance_id (instance_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建tasks表
CREATE TABLE IF NOT EXISTS tasks (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  provider_id int DEFAULT NULL,
  instance_id bigint unsigned DEFAULT NULL,
  type varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  progress int DEFAULT 0,
  message text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_tasks_user_id (user_id),
  KEY idx_tasks_provider_id (provider_id),
  KEY idx_tasks_instance_id (instance_id),
  KEY idx_tasks_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建verify_codes表
CREATE TABLE IF NOT EXISTS verify_codes (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  code varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  type varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  target varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  used tinyint(1) DEFAULT 0,
  expires_at datetime(3) NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_verify_codes_target_type (target,type),
  KEY idx_verify_codes_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建password_resets表
CREATE TABLE IF NOT EXISTS password_resets (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_uuid varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  token varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  expires_at datetime(3) NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_password_resets_token (token),
  KEY idx_password_resets_user_uuid (user_uuid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建system_configs表
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

-- 创建announcements表
CREATE TABLE IF NOT EXISTS announcements (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  content longtext COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  content_html longtext COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  type varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'homepage',
  status bigint DEFAULT 1,
  priority bigint DEFAULT 0,
  is_sticky tinyint(1) DEFAULT 0,
  start_time datetime(3) DEFAULT NULL,
  end_time datetime(3) DEFAULT NULL,
  created_by bigint unsigned DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_announcements_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建system_images表
CREATE TABLE IF NOT EXISTS system_images (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  type varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  os varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  version varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  arch varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'amd64',
  size bigint DEFAULT 0,
  status varchar(32) DEFAULT 'available',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_system_images_name_type (name,type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建captchas表
CREATE TABLE IF NOT EXISTS captchas (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  captcha_id varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  code varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  expires_at datetime(3) NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_captchas_captcha_id (captcha_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建invite_codes表
CREATE TABLE IF NOT EXISTS invite_codes (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  code varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  creator_id bigint unsigned NOT NULL,
  creator_name varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  max_uses bigint DEFAULT 1,
  used_count bigint DEFAULT 0,
  expires_at datetime(3) DEFAULT NULL,
  status bigint DEFAULT 1,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_invite_codes_code (code),
  KEY idx_invite_codes_creator_id (creator_id),
  KEY idx_invite_codes_expires_at (expires_at),
  KEY idx_invite_codes_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建invite_code_usages表
CREATE TABLE IF NOT EXISTS invite_code_usages (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  invite_code_id bigint unsigned NOT NULL,
  user_id int NOT NULL,
  ip varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  user_agent text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  used_at datetime(3) DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_invite_code_usages_invite_code_id (invite_code_id),
  KEY idx_invite_code_usages_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建audit_logs表
CREATE TABLE IF NOT EXISTS audit_logs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int DEFAULT NULL,
  username varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  action varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  resource_type varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  resource_id bigint unsigned DEFAULT NULL,
  ip varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  user_agent text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  details text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_audit_logs_user_id (user_id),
  KEY idx_audit_logs_action (action),
  KEY idx_audit_logs_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建pending_deletions表
CREATE TABLE IF NOT EXISTS pending_deletions (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  provider_id int NOT NULL,
  instance_id bigint unsigned NOT NULL,
  instance_name varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  deletion_time datetime(3) NOT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_pending_deletions_provider_id (provider_id),
  KEY idx_pending_deletions_status (status),
  KEY idx_pending_deletions_deletion_time (deletion_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建configuration_tasks表
CREATE TABLE IF NOT EXISTS configuration_tasks (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  type varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  provider_id int DEFAULT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  progress int DEFAULT 0,
  message text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_configuration_tasks_type (type),
  KEY idx_configuration_tasks_provider_id (provider_id),
  KEY idx_configuration_tasks_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建jwt_secrets表
CREATE TABLE IF NOT EXISTS jwt_secrets (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  secret_key varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_jwt_secrets_secret_key (secret_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建domains表
CREATE TABLE IF NOT EXISTS domains (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  domain varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  provider_id int DEFAULT NULL,
  instance_id bigint unsigned DEFAULT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_domains_domain (domain),
  KEY idx_domains_user_id (user_id),
  KEY idx_domains_provider_id (provider_id),
  KEY idx_domains_instance_id (instance_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建domain_configs表
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

-- 创建agents表
CREATE TABLE IF NOT EXISTS agents (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  name varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_agents_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建sub_user_relations表
CREATE TABLE IF NOT EXISTS sub_user_relations (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  agent_id bigint unsigned NOT NULL,
  user_id int NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_sub_user_relations_user_id (user_id),
  KEY idx_sub_user_relations_agent_id (agent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建commissions表
CREATE TABLE IF NOT EXISTS commissions (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  agent_id bigint unsigned NOT NULL,
  user_id int NOT NULL,
  order_id varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  amount decimal(10,2) DEFAULT 0,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_commissions_agent_id (agent_id),
  KEY idx_commissions_user_id (user_id),
  KEY idx_commissions_order_id (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建kyc_records表
CREATE TABLE IF NOT EXISTS kyc_records (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  real_name varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  id_number varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  status varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_kyc_records_user_id (user_id),
  KEY idx_kyc_records_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
