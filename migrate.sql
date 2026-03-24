-- 数据迁移脚本
-- 手动映射字段以适应新数据库结构

-- 先删除新数据库中的数据
DELETE FROM wallet_transactions;
DELETE FROM verify_codes;
DELETE FROM user_wallets;
DELETE FROM user_traffic_histories;
DELETE FROM user_roles;
DELETE FROM user_permissions;
DELETE FROM traffic_monitor_tasks;
DELETE FROM tasks;
DELETE FROM system_images;
DELETE FROM system_configs;
DELETE FROM sub_user_relations;
DELETE FROM site_configs;
DELETE FROM roles;
DELETE FROM resource_reservations;
DELETE FROM redemption_codes;
DELETE FROM redemption_code_usages;
DELETE FROM providers;
DELETE FROM provider_traffic_histories;
DELETE FROM products;
DELETE FROM product_purchases;
DELETE FROM ports;
DELETE FROM pmacct_traffic_records;
DELETE FROM pmacct_monitors;
DELETE FROM performance_metrics;
DELETE FROM pending_deletions;
DELETE FROM password_resets;
DELETE FROM orders;
DELETE FROM o_auth2_providers;
DELETE FROM kyc_records;
DELETE FROM jwt_secrets;
DELETE FROM invite_codes;
DELETE FROM invite_code_usages;
DELETE FROM instances;
DELETE FROM instance_traffic_histories;
DELETE FROM domains;
DELETE FROM domain_configs;
DELETE FROM configuration_tasks;
DELETE FROM commissions;
DELETE FROM captchas;
DELETE FROM audit_logs;
DELETE FROM announcements;
DELETE FROM agents;
DELETE FROM users;

-- 导入用户数据（手动映射字段）
INSERT INTO `users` (
  `id`, `uuid`, `username`, `email`, `password`, `user_type`, `level`, `status`, 
  `nickname`, `phone`, `email_verified`, `real_name_verified`, `level_expire_at`, 
  `max_instances`, `max_cpu`, `max_memory`, `max_disk`, `created_at`, `updated_at`, 
  `deleted_at`, `telegram`, `qq`, `avatar`, `used_quota`, `total_quota`, `total_traffic`, 
  `traffic_reset_at`, `traffic_limited`, `max_bandwidth`, `invite_code`, `last_login_at`, 
  `o_auth2_provider_id`, `o_auth2_uid`, `o_auth2_username`, `o_auth2_email`, `o_auth2_avatar`, `o_auth2_extra`
) VALUES
-- 用户1
(1,'6939f1fd-a925-4d7a-8958-8ee19a9aee7a','admin','admin@ypvps.com','$2a$10$nyAOEDpaP9GH01njj4ymDuPJn4tZA0ibGQT3TQARAQyWw5PDYq/O6','admin',5,1,'admin','',0,0,'2026-06-04 16:41:56.325',1,2,1024,2048,'2026-03-06 15:37:07.000','2026-03-12 21:08:17.339',NULL,'','','',3590,0,512000,NULL,0,500,'','2026-03-12 21:08:17.338',0,'','','','',''),
-- 用户2
(2,'61a9e805-157f-41e6-adae-d6e0bc036e2d','user','qdmz@vip.qq.com','$2a$10$mH153cHmiB.q50yDQLSlKulSG2FDojHHWUenf8hBI8wWivJlj8OXK','user',1,1,'user','',0,0,'2026-04-05 16:50:43.199',1,1,512,10240,'2026-03-06 15:37:07.000','2026-03-12 20:30:03.892',NULL,'','','',0,0,0,NULL,0,100,'','2026-03-06 15:49:57.953',0,'','','','',''),
-- 用户3
(3,'44ff4a62-46c3-4c6a-a8e9-fe8b559bac22','byzhenyu','','$2a$10$K67b4BN4.JMkoa0jsV/MEOpnJbfwNkdA3v247E7H6WB7gjZ28I1ta','user',3,0,'','',0,0,'2026-04-06 19:23:14.508',1,1,128,512,'2026-03-06 19:23:14.508','2026-03-07 09:30:57.376',NULL,'','','',632,0,0,'2026-04-01 00:00:00.000',0,2,'','2026-03-06 19:23:14.582',0,'','','','',''),
-- 用户4
(4,'b3d86fdb-5cf4-46c0-a6b6-7c6c81656bab','399107679@qq.com','','$2a$10$hFtXgT3XpkbIxHjD2sCs1ew9t69njfWQv8NHF5blFdF7QsY4BJVoy','user',2,1,'','',0,0,'2026-05-06 20:11:59.037',1,2,1024,2048,'2026-03-06 20:11:59.038','2026-03-10 22:13:25.970',NULL,'','','',106,0,0,'2026-04-01 00:00:00.000',0,2,'','2026-03-07 23:46:27.800',0,'','','','',''),
-- 用户5
(5,'cab8d3bb-cb29-4f57-b405-c39b1fa42cc5','ypvps','admin@qdmz.biters.edu.kg','$2a$10$vcLBkvhxq64RaDci9DdUXOLk1wan0U8A8qs3.Ds0Rib9o8idQJ8w.','user',1,1,'ypvps','13885852635',0,0,'2026-04-07 10:36:40.085',1,1,128,512,'2026-03-07 10:36:40.085','2026-03-09 11:53:59.900',NULL,'','','',0,0,0,'2026-04-01 00:00:00.000',0,100,'','2026-03-07 11:56:18.716',0,'','','','',''),
-- 用户6
(6,'5e4cdd79-313e-4415-a0db-14fd95c99b93','root','','$2a$10$Oa2bMBEAwB.lAeEyRb4vO.w6/JKsSU28k3e9lLXDBD0Cs43OpZV.e','user',4,1,'','',0,0,'2026-05-07 19:28:02.051',1,4,4096,4096,'2026-03-07 19:28:02.052','2026-03-07 19:34:16.999',NULL,'','','',851,0,0,'2026-04-01 00:00:00.000',0,100,'','2026-03-07 19:28:02.127',0,'','','','',''),
-- 用户7
(7,'0765eb45-c44d-4ba4-b4fd-c5265e59a758','wlqb1981','','$2a$10$IBV6vzDomVs6MBnh1HeaUunc0NgIc.zpWngJL3xU26LlNqOxut.o.','user',1,1,'','',0,0,'2026-04-07 20:05:50.901',1,1,128,512,'2026-03-07 20:05:50.902','2026-03-07 20:14:31.152',NULL,'','','',106,0,0,'2026-04-01 00:00:00.000',0,100,'','2026-03-07 20:07:23.990',0,'','','','',''),
-- 用户8
(8,'d4a785a1-9892-4266-ba95-a3114b8c2340','deuspamm','','$2a$10$lVV7wiEx3ymySIwrlrNWHOFXtYgAguZZt95uK778TUZlrmzP1xFz2','user',1,1,'','',0,0,'2026-03-09 00:04:50.971',1,1,128,512,'2026-03-08 00:04:50.972','2026-03-08 00:04:51.047',NULL,'','','',0,0,0,'2026-04-01 00:00:00.000',0,100,'','2026-03-08 00:04:51.046',0,'','','','',''),
-- 用户9
(9,'ec051e46-55ed-44f8-9fa4-78a1a558b3a8','mytest','qdmz@10000.edu.pl','$2a$12$4NKJx566BsSk3P5/Z2t05.OXxAuOLFFnfn2qxjxhH/mrzadQhAWaK','user',1,1,'mytest','',0,0,'2026-03-12 16:03:08.799',1,1,128,512,'2026-03-11 16:03:08.799','2026-03-11 16:05:04.568',NULL,'','','',0,0,0,'2026-04-01 00:00:00.000',0,100,'','2026-03-11 16:03:08.805',0,'','','','',''),
-- 用户10
(10,'61199b8b-7df0-418a-a506-d53df5db336c','qdmz','','$2a$10$DmgQpkXZd/ZA9EbZooOxQu0op/4fwdaf4d0c6TJ1dKn7UclCI1ZUy','user',1,1,'qdmz','',0,0,'2026-04-11 11:45:29.280',1,1,128,512,'2026-03-12 11:40:57.825','2026-03-12 12:11:38.157',NULL,'','','https://avatars.githubusercontent.com/u/6871334?v=4',106,0,0,NULL,0,100,'','2026-03-12 11:40:57.829',1,'6.871334e+06','qdmz','','https://avatars.githubusercontent.com/u/6871334?v=4','{"avatar_url":"https://avatars.githubusercontent.com/u/6871334?v=4","bio":null,"blog":"https://www.ypvps.com","company":null,"created_at":"2014-03-06T09:35:41Z","email":null,"events_url":"https://api.github.com/users/qdmz/events{/privacy}","followers":0,"followers_url":"https://api.github.com/users/qdmz/followers","following":0,"following_url":"https://api.github.com/users/qdmz/following{/other_user}","gists_url":"https://api.github.com/users/qdmz/gists{/gist_id}","gravatar_id":"","hireable":null,"html_url":"https://github.com/qdmz","id":6871334,"location":null,"login":"qdmz","name":null,"node_id":"MDQ6VXNlcjY4NzEzMzQ=","notification_email":null,"organizations_url":"https://api.github.com/users/qdmz/orgs","public_gists":1,"public_repos":25,"received_events_url":"https://api.github.com/users/qdmz/received_events","repos_url":"https://api.github.com/users/qdmz/repos","site_admin":false,"starred_url":"https://api.github.com/users/qdmz/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/qdmz/subscriptions","twitter_username":"qdmzhost","type":"User","updated_at":"2026-03-02T06:08:34Z","url":"https://api.github.com/users/qdmz","user_view_type":"public"}');

-- 导入节点数据
INSERT INTO `providers` (
  `id`, `uuid`, `name`, `type`, `endpoint`, `ssh_port`, `username`, `password`, 
  `status`, `region`, `country`, `created_at`, `updated_at`
) VALUES
-- 节点1
(1,'bd0fa1e1-82c3-43ce-8e8d-25b7cfd93b66','heyun','proxmox','154.12.84.134',22,'root','thanks12A#','active','hk','中国香港','2026-03-06 16:21:13.956','2026-03-13 11:39:17.629'),
-- 节点2
(2,'fa00a303-0753-49b3-961a-111f8d717bc6','heyunus','proxmox','38.165.47.49',22,'root','thanks12A#','active','usa','美国','2026-03-06 19:14:03.477','2026-03-13 11:39:16.364'),
-- 节点3
(3,'3rd-node-uuid','heyhk3','proxmox','156.233.233.6',22,'root','thanks12A#','active','hk','中国香港','2026-03-07 10:13:34.683','2026-03-13 11:39:15.543');

-- 导入实例数据
INSERT INTO `instances` (
  `id`, `uuid`, `name`, `provider`, `provider_id`, `status`, `image`, `instance_type`, 
  `cpu`, `memory`, `disk`, `bandwidth`, `private_ip`, `public_ip`, `ssh_port`, 
  `username`, `password`, `os_type`, `expired_at`, `user_id`, `created_at`, `updated_at`
) VALUES
-- 实例1
(1,'225e25b2-d989-4957-ad4d-1ba6a8297d57','heyun-388c','heyun',1,'running','alpine-3.23-64_cloud','container',1,64,200,1,'172.16.1.16','154.12.84.134',22,'root','j9vopl9ijtsa','alpine','2026-04-05 16:50:43.199',2,'2026-03-06 16:42:26.790','2026-03-10 15:11:02.136'),
-- 实例2
(2,'71be3971-d7ef-465f-93bf-5ab2a4c3921c','heyunus-d9c8-old-1772854755','heyunus',2,'resetting','debian13','vm',3,1536,3072,3,'172.16.1.9','38.165.47.49',22,'root','d0bd6s5dco3r','debian','2026-04-06 20:11:59.037',4,'2026-03-06 19:36:15.357','2026-03-07 11:39:15.190'),
-- 实例3
(3,'04776351-1eda-4a85-b5e2-8adc03bcd95e','heyunus-c85c','heyunus',2,'running','alpine-3.23-64_cloud','container',1,128,512,2,'172.16.1.10','38.165.47.49',22,'root','y3aql2a59c57','alpine','2026-04-05 16:50:43.199',2,'2026-03-06 23:03:00.356','2026-03-11 01:05:22.524'),
-- 实例4
(4,'b8512407-5e16-4fd8-b462-a75b2134601f','heyhk3-c3ae','heyhk3',3,'running','debian13','vm',3,1536,3072,300,'172.16.1.14','156.233.233.6',22,'root','va0qmomd80jg','debian','2026-05-06 20:11:59.037',4,'2026-03-07 10:13:34.683','2026-03-10 22:13:25.969'),
-- 实例5
(5,'6364f5f1-0dec-437c-9d6d-a0e479616b37','heyhk3-f9e0','heyhk3',3,'running','alpine-3.23-64_cloud','container',1,128,512,100,'172.16.1.15','156.233.233.6',22,'root','pfagg0a8xii2','alpine','2026-04-07 10:36:40.085',5,'2026-03-07 10:38:44.683','2026-03-07 10:41:53.352');

-- 完成迁移
SELECT '数据迁移完成' AS message;