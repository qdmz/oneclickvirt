-- 最终数据迁移脚本
-- 将旧数据库中的用户、节点、任务、实例和端口数据迁移到新数据库

-- 先删除新数据库中可能存在的旧数据，避免冲突
DELETE FROM ports;
DELETE FROM instances;
DELETE FROM tasks;
DELETE FROM providers;
DELETE FROM users;

-- 导入用户数据（修正字段顺序和类型）
INSERT INTO `users` (`id`, `uuid`, `username`, `email`, `password`, `user_type`, `level`, `status`, `nickname`, `phone`, `email_verified`, `real_name_verified`, `level_expire_at`, `max_instances`, `max_cpu`, `max_memory`, `max_disk`, `created_at`, `updated_at`, `deleted_at`, `telegram`, `qq`, `avatar`, `used_quota`, `total_quota`, `total_traffic`, `traffic_reset_at`, `traffic_limited`, `max_bandwidth`, `invite_code`, `last_login_at`, `o_auth2_provider_id`, `o_auth2_uid`, `o_auth2_username`, `o_auth2_email`, `o_auth2_avatar`, `o_auth2_extra`)
VALUES 
(1,'6939f1fd-a925-4d7a-8958-8ee19a9aee7a','admin','admin@ypvps.com','$2a$10$nyAOEDpaP9GH01njj4ymDuPJn4tZA0ibGQT3TQARAQyWw5PDYq/O6','admin',5,1,'admin','',0,0,'2026-06-04 16:41:56.325',1,2,1024,2048,'2026-03-06 15:37:07.000','2026-03-12 21:08:17.339',NULL,'','','',3590,0,512000,NULL,0,500,'','2026-03-12 21:08:17.338',0,'','','','','',0,0),
(2,'61a9e805-157f-41e6-adae-d6e0bc036e2d','user','qdmz@vip.qq.com','$2a$10$mH153cHmiB.q50yDQLSlKulSG2FDojHHWUenf8hBI8wWivJlj8OXK','user',1,1,'user','',0,0,'2026-04-05 16:50:43.199',1,1,512,10240,'2026-03-06 15:37:07.000','2026-03-12 20:30:03.892',NULL,'','','',0,0,0,NULL,0,100,'','2026-03-06 15:49:57.953',0,'','','','','',0,0),
(3,'44ff4a62-46c3-4c6a-a8e9-fe8b559bac22','byzhenyu','','$2a$10$K67b4BN4.JMkoa0jsV/MEOpnJbfwNkdA3v247E7H6WB7gjZ28I1ta','user',3,0,'','',0,0,'2026-04-06 19:23:14.508',1,1,128,512,'2026-03-06 19:23:14.508','2026-03-07 09:30:57.376',NULL,'','','',632,0,0,'2026-04-01 00:00:00.000',0,2,'','2026-03-06 19:23:14.582',0,'','','','','',0,0),
(4,'b3d86fdb-5cf4-46c0-a6b6-7c6c81656bab','399107679@qq.com','','$2a$10$hFtXgT3XpkbIxHjD2sCs1ew9t69njfWQv8NHF5blFdF7QsY4BJVoy','user',2,1,'','',0,0,'2026-05-06 20:11:59.037',1,2,1024,2048,'2026-03-06 20:11:59.038','2026-03-10 22:13:25.970',NULL,'','','',106,0,0,'2026-04-01 00:00:00.000',0,2,'','2026-03-07 23:46:27.800',0,'','','','','',0,0),
(5,'cab8d3bb-cb29-4f57-b405-c39b1fa42cc5','ypvps','admin@qdmz.biters.edu.kg','$2a$10$vcLBkvhxq64RaDci9DdUXOLk1wan0U8A8qs3.Ds0Rib9o8idQJ8w.','user',1,1,'ypvps','13885852635',0,0,'2026-04-07 10:36:40.085',1,1,128,512,'2026-03-07 10:36:40.085','2026-03-09 11:53:59.900',NULL,'','','',0,0,0,'2026-04-01 00:00:00.000',0,100,'','2026-03-07 11:56:18.716',0,'','','','','',0,0);

-- 导入节点数据
INSERT INTO `providers` (`id`, `uuid`, `name`, `type`, `endpoint`, `port_ip`, `ssh_port`, `username`, `password`, `status`, `region`, `country`, `country_code`, `city`, `version`, `container_enabled`, `virtual_machine_enabled`, `allow_claim`, `ipv4_port_mapping_method`, `ipv6_port_mapping_method`, `used_quota`, `total_quota`, `architecture`, `expires_at`, `is_frozen`, `storage_pool`, `network_type`, `default_inbound_bandwidth`, `default_outbound_bandwidth`, `max_inbound_bandwidth`, `max_outbound_bandwidth`, `enable_traffic_control`, `max_traffic`, `traffic_limited`, `traffic_reset_at`, `created_at`, `updated_at`)
VALUES 
(1,'bd0fa1e1-82c3-43ce-8e8d-25b7cfd93b66','heyun','proxmox','154.12.84.134','',22,'root','thanks12A#','active','hk','中国香港','HK','','9.1.6',1,1,1,'iptables','native',0,0,'amd64','2029-03-28 08:00:00.000',0,'local','nat_ipv4_ipv6',300,300,1000,1000,0,1048576,0,'2026-04-01 00:00:00.000','2026-03-06 16:21:13.956','2026-03-13 11:39:17.629'),
(2,'fa00a303-0753-49b3-961a-111f8d717bc6','heyunus','proxmox','38.165.47.49','',22,'root','thanks12A#','active','usa','美国','US','洛杉矶','9.1.6',1,1,1,'iptables','native',0,0,'amd64','2029-03-22 08:00:00.000',0,'local','nat_ipv4_ipv6',300,300,1000,1000,0,1048576,0,'2026-04-01 00:00:00.000','2026-03-06 19:14:03.477','2026-03-13 11:39:16.364');

-- 完成迁移
SELECT '数据迁移完成' AS message;