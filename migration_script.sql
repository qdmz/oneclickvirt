-- 数据迁移脚本
-- 将旧数据库中的用户、节点、任务、实例和端口数据迁移到新数据库

-- 先删除新数据库中可能存在的旧数据，避免冲突
DELETE FROM ports;
DELETE FROM instances;
DELETE FROM tasks;
DELETE FROM providers;
DELETE FROM users;

-- 导入用户数据
INSERT INTO `users` (`id`, `uuid`, `username`, `email`, `password`, `user_type`, `level`, `status`, `nickname`, `phone`, `email_verified`, `real_name_verified`, `level_expire_at`, `max_instances`, `max_cpu`, `max_memory`, `max_disk`, `created_at`, `updated_at`, `deleted_at`, `telegram`, `qq`, `avatar`, `used_quota`, `total_quota`, `total_traffic`, `traffic_reset_at`, `traffic_limited`, `max_bandwidth`, `invite_code`, `last_login_at`, `o_auth2_provider_id`, `o_auth2_uid`, `o_auth2_username`, `o_auth2_email`, `o_auth2_avatar`, `o_auth2_extra`)
SELECT 
  id, uuid, username, email, password, user_type, level, status, nickname, phone, 
  email_verified, real_name_verified, level_expire_at, max_instances, max_cpu, max_memory, max_disk, 
  created_at, updated_at, deleted_at, telegram, qq, avatar, used_quota, total_quota, total_traffic, 
  traffic_reset_at, traffic_limited, max_bandwidth, invite_code, last_login_at, 
  o_auth2_provider_id, o_auth2_uid, o_auth2_username, o_auth2_email, o_auth2_avatar, o_auth2_extra
FROM `oneclickvirt_old`.`users`;

-- 导入节点数据
INSERT INTO `providers` (`id`, `uuid`, `created_at`, `updated_at`, `name`, `type`, `endpoint`, `port_ip`, `ssh_port`, `username`, `password`, `ssh_key`, `token`, `config`, `status`, `region`, `country`, `country_code`, `city`, `version`, `container_enabled`, `virtual_machine_enabled`, `supported_types`, `allow_claim`, `ipv4_port_mapping_method`, `ipv6_port_mapping_method`, `used_quota`, `total_quota`, `architecture`, `expires_at`, `is_frozen`, `storage_pool`, `storage_pool_path`, `cert_path`, `key_path`, `ca_cert_path`, `cert_fingerprint`, `trusted_fingerprint`, `api_status`, `ssh_status`, `last_api_check`, `last_ssh_check`, `auth_config`, `config_version`, `auto_configured`, `last_config_update`, `config_backup_path`, `cert_content`, `key_content`, `token_content`, `node_cpu_cores`, `node_memory_total`, `node_disk_total`, `allow_concurrent_tasks`, `max_concurrent_tasks`, `ssh_connect_timeout`, `ssh_execute_timeout`, `task_poll_interval`, `enable_task_polling`, `execution_rule`, `max_container_instances`, `max_vm_instances`, `container_limit_cpu`, `container_limit_memory`, `container_limit_disk`, `vm_limit_cpu`, `vm_limit_memory`, `vm_limit_disk`, `default_port_count`, `port_range_start`, `port_range_end`, `next_available_port`, `network_type`, `default_inbound_bandwidth`, `default_outbound_bandwidth`, `max_inbound_bandwidth`, `max_outbound_bandwidth`, `enable_traffic_control`, `max_traffic`, `traffic_limited`, `traffic_reset_at`, `traffic_count_mode`, `traffic_multiplier`, `traffic_stats_mode`, `traffic_collect_interval`, `traffic_collect_batch_size`, `traffic_limit_check_interval`, `traffic_limit_check_batch_size`, `traffic_auto_reset_interval`, `traffic_auto_reset_batch_size`, `used_cpu_cores`, `used_memory`, `used_disk`, `container_count`, `vm_count`, `resource_synced`, `resource_synced_at`, `count_cache_expiry`, `available_cpu_cores`, `available_memory`, `used_instances`, `level_limits`, `host_name`, `container_privileged`, `container_allow_nesting`, `container_enable_lxcfs`, `container_cpu_allowance`, `container_memory_swap`, `container_max_processes`, `container_disk_io_limit`)
SELECT 
  id, uuid, created_at, updated_at, name, type, endpoint, port_ip, ssh_port, username, password, ssh_key, token, config, status, region, country, country_code, city, version, 
  container_enabled, virtual_machine_enabled, supported_types, allow_claim, ipv4_port_mapping_method, ipv6_port_mapping_method, used_quota, total_quota, architecture, expires_at, 
  is_frozen, storage_pool, storage_pool_path, cert_path, key_path, ca_cert_path, cert_fingerprint, trusted_fingerprint, api_status, ssh_status, last_api_check, last_ssh_check, 
  auth_config, config_version, auto_configured, last_config_update, config_backup_path, cert_content, key_content, token_content, node_cpu_cores, node_memory_total, node_disk_total, 
  allow_concurrent_tasks, max_concurrent_tasks, ssh_connect_timeout, ssh_execute_timeout, task_poll_interval, enable_task_polling, execution_rule, max_container_instances, max_vm_instances, 
  container_limit_cpu, container_limit_memory, container_limit_disk, vm_limit_cpu, vm_limit_memory, vm_limit_disk, default_port_count, port_range_start, port_range_end, next_available_port, 
  network_type, default_inbound_bandwidth, default_outbound_bandwidth, max_inbound_bandwidth, max_outbound_bandwidth, enable_traffic_control, max_traffic, traffic_limited, traffic_reset_at, 
  traffic_count_mode, traffic_multiplier, traffic_stats_mode, traffic_collect_interval, traffic_collect_batch_size, traffic_limit_check_interval, traffic_limit_check_batch_size, 
  traffic_auto_reset_interval, traffic_auto_reset_batch_size, used_cpu_cores, used_memory, used_disk, container_count, vm_count, resource_synced, resource_synced_at, count_cache_expiry, 
  available_cpu_cores, available_memory, used_instances, level_limits, host_name, container_privileged, container_allow_nesting, container_enable_lxcfs, container_cpu_allowance, 
  container_memory_swap, container_max_processes, container_disk_io_limit
FROM `oneclickvirt_old`.`providers`;

-- 导入任务数据
INSERT INTO `tasks` (`id`, `uuid`, `user_id`, `provider_id`, `instance_id`, `type`, `status`, `priority`, `progress`, `message`, `error`, `data`, `metadata`, `started_at`, `completed_at`, `created_at`, `updated_at`, `deleted_at`, `task_type`, `error_message`, `cancel_reason`, `status_message`, `task_data`, `estimated_duration`, `timeout_duration`, `preallocated_cpu`, `preallocated_memory`, `preallocated_disk`, `preallocated_bandwidth`, `can_force_stop`, `is_force_stoppable`)
SELECT 
  id, uuid, user_id, provider_id, instance_id, task_type, status, 0, progress, status_message, error_message, task_data, '', started_at, completed_at, 
  created_at, updated_at, deleted_at, task_type, error_message, cancel_reason, status_message, task_data, estimated_duration, timeout_duration, 
  preallocated_cpu, preallocated_memory, preallocated_disk, preallocated_bandwidth, can_force_stop, is_force_stoppable
FROM `oneclickvirt_old`.`tasks`;

-- 导入实例数据
INSERT INTO `instances` (`id`, `uuid`, `created_at`, `updated_at`, `deleted_at`, `name`, `provider`, `provider_id`, `status`, `image`, `instance_type`, `cpu`, `memory`, `disk`, `bandwidth`, `network`, `private_ip`, `public_ip`, `ipv6_address`, `public_ipv6`, `ssh_port`, `port_range_start`, `port_range_end`, `username`, `password`, `os_type`, `region`, `max_traffic`, `traffic_limited`, `traffic_limit_reason`, `pmacct_interface_v4`, `pmacct_interface_v6`, `expired_at`, `user_id`)
SELECT 
  id, uuid, created_at, updated_at, deleted_at, name, provider, provider_id, status, image, instance_type, cpu, memory, disk, bandwidth, network, 
  private_ip, public_ip, ipv6_address, public_ipv6, ssh_port, port_range_start, port_range_end, username, password, os_type, region, 
  max_traffic, traffic_limited, traffic_limit_reason, pmacct_interface_v4, pmacct_interface_v6, expired_at, user_id
FROM `oneclickvirt_old`.`instances`;

-- 导入端口数据
INSERT INTO `ports` (`id`, `instance_id`, `provider_id`, `protocol`, `private_port`, `public_port`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`, `host_port`, `host_port_end`, `guest_port`, `guest_port_end`, `port_count`, `is_ssh`, `is_automatic`, `port_type`, `ipv6_enabled`, `ipv6_address`, `mapping_method`)
SELECT 
  id, instance_id, provider_id, protocol, guest_port, host_port, description, status, created_at, updated_at, deleted_at, 
  host_port, host_port_end, guest_port, guest_port_end, port_count, is_ssh, is_automatic, port_type, ipv6_enabled, ipv6_address, mapping_method
FROM `oneclickvirt_old`.`ports`;

-- 完成迁移
SELECT '数据迁移完成' AS message;