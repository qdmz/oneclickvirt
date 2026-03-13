-- MySQL dump 10.13  Distrib 8.0.45, for Linux (x86_64)
--
-- Host: localhost    Database: oneclickvirt
-- ------------------------------------------------------
-- Server version	8.0.45

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `nickname` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `telegram` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `qq` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `level` bigint DEFAULT '1',
  `level_expire_at` datetime(3) DEFAULT NULL,
  `user_type` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'user',
  `used_quota` bigint DEFAULT '0',
  `total_quota` bigint DEFAULT '0',
  `total_traffic` bigint DEFAULT '0',
  `traffic_reset_at` datetime(3) DEFAULT NULL,
  `traffic_limited` tinyint(1) DEFAULT '0',
  `max_instances` bigint DEFAULT '1',
  `max_cpu` bigint DEFAULT '1',
  `max_memory` bigint DEFAULT '512',
  `max_disk` bigint DEFAULT '10240',
  `max_bandwidth` bigint DEFAULT '100',
  `invite_code` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `last_login_at` datetime(3) DEFAULT NULL,
  `o_auth2_provider_id` bigint unsigned DEFAULT NULL,
  `o_auth2_uid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `o_auth2_username` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `o_auth2_email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `o_auth2_avatar` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `o_auth2_extra` text COLLATE utf8mb4_unicode_ci,
  `email_verified` tinyint(1) DEFAULT '0' COMMENT 'email verified',
  `real_name_verified` tinyint(1) DEFAULT '0' COMMENT 'real name verified',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_uuid` (`uuid`),
  UNIQUE KEY `idx_users_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_level` (`level`),
  KEY `idx_users_o_auth2_provider_id` (`o_auth2_provider_id`),
  KEY `idx_users_o_auth2_uid` (`o_auth2_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'6939f1fd-a925-4d7a-8958-8ee19a9aee7a','2026-03-06 15:37:07.000','2026-03-12 21:08:17.339',NULL,'admin','$2a$10$nyAOEDpaP9GH01njj4ymDuPJn4tZA0ibGQT3TQARAQyWw5PDYq/O6','admin','admin@ypvps.com','','','','',1,5,'2026-06-04 16:41:56.325','admin',3590,0,512000,NULL,0,1,2,1024,2048,500,'','2026-03-12 21:08:17.338',0,'','','','','',0,0),(2,'61a9e805-157f-41e6-adae-d6e0bc036e2d','2026-03-06 15:37:07.000','2026-03-12 20:30:03.892',NULL,'user','$2a$10$mH153cHmiB.q50yDQLSlKulSG2FDojHHWUenf8hBI8wWivJlj8OXK','user','qdmz@vip.qq.com','','','','',1,1,'2026-04-05 16:50:43.199','user',0,0,0,NULL,0,1,1,512,10240,100,'','2026-03-06 15:49:57.953',0,'','','','','',0,0),(3,'44ff4a62-46c3-4c6a-a8e9-fe8b559bac22','2026-03-06 19:23:14.508','2026-03-07 09:30:57.376',NULL,'byzhenyu','$2a$10$K67b4BN4.JMkoa0jsV/MEOpnJbfwNkdA3v247E7H6WB7gjZ28I1ta','','','','','','',0,3,'2026-04-06 19:23:14.508','user',632,0,0,'2026-04-01 00:00:00.000',0,1,1,128,512,2,'','2026-03-06 19:23:14.582',0,'','','','','',0,0),(4,'b3d86fdb-5cf4-46c0-a6b6-7c6c81656bab','2026-03-06 20:11:59.038','2026-03-10 22:13:25.970',NULL,'399107679@qq.com','$2a$10$hFtXgT3XpkbIxHjD2sCs1ew9t69njfWQv8NHF5blFdF7QsY4BJVoy','','','','','','',1,2,'2026-05-06 20:11:59.037','user',106,0,0,'2026-04-01 00:00:00.000',0,1,2,1024,2048,2,'','2026-03-07 23:46:27.800',0,'','','','','',0,0),(5,'cab8d3bb-cb29-4f57-b405-c39b1fa42cc5','2026-03-07 10:36:40.085','2026-03-09 11:53:59.900',NULL,'ypvps','$2a$10$vcLBkvhxq64RaDci9DdUXOLk1wan0U8A8qs3.Ds0Rib9o8idQJ8w.','ypvps','admin@qdmz.biters.edu.kg','13885852635','','','',1,1,'2026-04-07 10:36:40.085','user',0,0,0,'2026-04-01 00:00:00.000',0,1,1,128,512,100,'','2026-03-07 11:56:18.716',0,'','','','','',0,0),(6,'5e4cdd79-313e-4415-a0db-14fd95c99b93','2026-03-07 19:28:02.052','2026-03-07 19:34:16.999',NULL,'root','$2a$10$Oa2bMBEAwB.lAeEyRb4vO.w6/JKsSU28k3e9lLXDBD0Cs43OpZV.e','','','','','','',1,4,'2026-05-07 19:28:02.051','user',851,0,0,'2026-04-01 00:00:00.000',0,1,4,4096,4096,100,'','2026-03-07 19:28:02.127',0,'','','','','',0,0),(7,'0765eb45-c44d-4ba4-b4fd-c5265e59a758','2026-03-07 20:05:50.902','2026-03-07 20:14:31.152',NULL,'wlqb1981','$2a$10$IBV6vzDomVs6MBnh1HeaUunc0NgIc.zpWngJL3xU26LlNqOxut.o.','','','','','','',1,1,'2026-04-07 20:05:50.901','user',106,0,0,'2026-04-01 00:00:00.000',0,1,1,128,512,100,'','2026-03-07 20:07:23.990',0,'','','','','',0,0),(8,'d4a785a1-9892-4266-ba95-a3114b8c2340','2026-03-08 00:04:50.972','2026-03-08 00:04:51.047',NULL,'deuspamm','$2a$10$lVV7wiEx3ymySIwrlrNWHOFXtYgAguZZt95uK778TUZlrmzP1xFz2','','','','','','',1,1,'2026-03-09 00:04:50.971','user',0,0,0,'2026-04-01 00:00:00.000',0,1,1,128,512,100,'','2026-03-08 00:04:51.046',0,'','','','','',0,0),(9,'ec051e46-55ed-44f8-9fa4-78a1a558b3a8','2026-03-11 16:03:08.799','2026-03-11 16:05:04.568',NULL,'mytest','$2a$12$4NKJx566BsSk3P5/Z2t05.OXxAuOLFFnfn2qxjxhH/mrzadQhAWaK','mytest','qdmz@10000.edu.pl','','','','',1,1,'2026-03-12 16:03:08.799','user',0,0,0,'2026-04-01 00:00:00.000',0,1,1,128,512,100,'','2026-03-11 16:03:08.805',0,'','','','','',0,0),(10,'61199b8b-7df0-418a-a506-d53df5db336c','2026-03-12 11:40:57.825','2026-03-12 12:11:38.157',NULL,'qdmz','$2a$10$DmgQpkXZd/ZA9EbZooOxQu0op/4fwdaf4d0c6TJ1dKn7UclCI1ZUy','qdmz','','','','','https://avatars.githubusercontent.com/u/6871334?v=4',1,1,'2026-04-11 11:45:29.280','user',106,0,0,NULL,0,1,1,128,512,100,'','2026-03-12 11:40:57.829',1,'6.871334e+06','qdmz','','https://avatars.githubusercontent.com/u/6871334?v=4','{\"avatar_url\":\"https://avatars.githubusercontent.com/u/6871334?v=4\",\"bio\":null,\"blog\":\"https://www.ypvps.com\",\"company\":null,\"created_at\":\"2014-03-06T09:35:41Z\",\"email\":null,\"events_url\":\"https://api.github.com/users/qdmz/events{/privacy}\",\"followers\":0,\"followers_url\":\"https://api.github.com/users/qdmz/followers\",\"following\":0,\"following_url\":\"https://api.github.com/users/qdmz/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/qdmz/gists{/gist_id}\",\"gravatar_id\":\"\",\"hireable\":null,\"html_url\":\"https://github.com/qdmz\",\"id\":6871334,\"location\":null,\"login\":\"qdmz\",\"name\":null,\"node_id\":\"MDQ6VXNlcjY4NzEzMzQ=\",\"notification_email\":null,\"organizations_url\":\"https://api.github.com/users/qdmz/orgs\",\"public_gists\":1,\"public_repos\":25,\"received_events_url\":\"https://api.github.com/users/qdmz/received_events\",\"repos_url\":\"https://api.github.com/users/qdmz/repos\",\"site_admin\":false,\"starred_url\":\"https://api.github.com/users/qdmz/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/qdmz/subscriptions\",\"twitter_username\":\"qdmzhost\",\"type\":\"User\",\"updated_at\":\"2026-03-02T06:08:34Z\",\"url\":\"https://api.github.com/users/qdmz\",\"user_view_type\":\"public\"}',0,0);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `providers`
--

DROP TABLE IF EXISTS `providers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `providers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `endpoint` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `port_ip` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ssh_port` bigint DEFAULT '22',
  `username` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ssh_key` text COLLATE utf8mb4_unicode_ci,
  `token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `config` text COLLATE utf8mb4_unicode_ci,
  `status` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  `region` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `country` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `country_code` varchar(8) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `city` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `version` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `container_enabled` tinyint(1) DEFAULT '1',
  `virtual_machine_enabled` tinyint(1) DEFAULT '0',
  `supported_types` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `allow_claim` tinyint(1) DEFAULT '1',
  `ipv4_port_mapping_method` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'device_proxy',
  `ipv6_port_mapping_method` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'device_proxy',
  `used_quota` bigint DEFAULT '0',
  `total_quota` bigint DEFAULT '0',
  `architecture` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'amd64',
  `expires_at` datetime(3) DEFAULT NULL,
  `is_frozen` tinyint(1) DEFAULT '0',
  `storage_pool` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'local',
  `storage_pool_path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `cert_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `key_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ca_cert_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cert_fingerprint` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `api_status` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'unknown',
  `ssh_status` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'unknown',
  `last_api_check` datetime(3) DEFAULT NULL,
  `last_ssh_check` datetime(3) DEFAULT NULL,
  `auth_config` text COLLATE utf8mb4_unicode_ci,
  `config_version` bigint DEFAULT '0',
  `auto_configured` tinyint(1) DEFAULT '0',
  `last_config_update` datetime(3) DEFAULT NULL,
  `config_backup_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cert_content` text COLLATE utf8mb4_unicode_ci,
  `key_content` text COLLATE utf8mb4_unicode_ci,
  `token_content` text COLLATE utf8mb4_unicode_ci,
  `node_cpu_cores` bigint DEFAULT '0',
  `node_memory_total` bigint DEFAULT '0',
  `node_disk_total` bigint DEFAULT '0',
  `allow_concurrent_tasks` tinyint(1) DEFAULT '0',
  `max_concurrent_tasks` bigint DEFAULT '1',
  `ssh_connect_timeout` bigint DEFAULT '30',
  `ssh_execute_timeout` bigint DEFAULT '300',
  `task_poll_interval` bigint DEFAULT '60',
  `enable_task_polling` tinyint(1) DEFAULT '1',
  `execution_rule` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'auto',
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
  `network_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'nat_ipv4',
  `default_inbound_bandwidth` bigint DEFAULT '300',
  `default_outbound_bandwidth` bigint DEFAULT '300',
  `max_inbound_bandwidth` bigint DEFAULT '1000',
  `max_outbound_bandwidth` bigint DEFAULT '1000',
  `enable_traffic_control` tinyint(1) DEFAULT '0',
  `max_traffic` bigint DEFAULT '1048576',
  `traffic_limited` tinyint(1) DEFAULT '0',
  `traffic_reset_at` datetime(3) DEFAULT NULL,
  `traffic_count_mode` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'both',
  `traffic_multiplier` double DEFAULT '1',
  `traffic_stats_mode` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'light',
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
  `level_limits` text COLLATE utf8mb4_unicode_ci,
  `host_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `container_privileged` tinyint(1) DEFAULT '0',
  `container_allow_nesting` tinyint(1) DEFAULT '0',
  `container_enable_lxcfs` tinyint(1) DEFAULT '1',
  `container_cpu_allowance` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT '100%',
  `container_memory_swap` tinyint(1) DEFAULT '1',
  `container_max_processes` bigint DEFAULT '0',
  `container_disk_io_limit` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trusted_fingerprint` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_providers_uuid` (`uuid`),
  UNIQUE KEY `idx_providers_name` (`name`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_region` (`region`),
  KEY `idx_providers_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `providers`
--

LOCK TABLES `providers` WRITE;
/*!40000 ALTER TABLE `providers` DISABLE KEYS */;
INSERT INTO `providers` VALUES (1,'bd0fa1e1-82c3-43ce-8e8d-25b7cfd93b66','2026-03-06 16:21:13.956','2026-03-13 11:39:17.629','heyun','proxmox','154.12.84.134','',22,'root','thanks12A#','','','','active','hk','中国香港','HK','','9.1.6',1,1,'',1,'iptables','native',0,0,'amd64','2029-03-28 08:00:00.000',0,'local','','','','','','unknown','online','2026-03-13 11:39:15.545','2026-03-13 11:39:15.545','',0,0,NULL,'','','','',16,15989,40960,0,1,30,300,60,1,'auto',20,10,1,1,1,1,1,1,5,33300,65535,33315,'nat_ipv4_ipv6',300,300,1000,1000,0,1048576,0,'2026-04-01 00:00:00.000','both',1,'light',90,10,90,10,1800,10,2,1600,5320,2,1,1,'2026-03-12 21:00:11.527','2026-03-12 21:05:11.527',14,14389,3,'{\"1\":{\"max-instances\":1,\"max-resources\":{\"bandwidth\":100,\"cpu\":1,\"disk\":10240,\"memory\":512},\"max-traffic\":102400},\"2\":{\"max-instances\":3,\"max-resources\":{\"bandwidth\":200,\"cpu\":2,\"disk\":20480,\"memory\":1024},\"max-traffic\":204800},\"3\":{\"max-instances\":5,\"max-resources\":{\"bandwidth\":500,\"cpu\":4,\"disk\":40960,\"memory\":2048},\"max-traffic\":307200},\"4\":{\"max-instances\":10,\"max-resources\":{\"bandwidth\":1000,\"cpu\":8,\"disk\":81920,\"memory\":4096},\"max-traffic\":409600},\"5\":{\"max-instances\":20,\"max-resources\":{\"bandwidth\":2000,\"cpu\":16,\"disk\":163840,\"memory\":8192},\"max-traffic\":512000}}','',0,0,1,'100%',1,0,'',''),(2,'fa00a303-0753-49b3-961a-111f8d717bc6','2026-03-06 19:14:03.477','2026-03-13 11:39:16.364','heyunus','proxmox','38.165.47.49','',22,'root','thanks12A#','','','','active','usa','美国','US','洛杉矶','9.1.6',1,1,'',1,'iptables','native',0,0,'amd64','2029-03-22 08:00:00.000',0,'local','','','','','','unknown','online','2026-03-13 11:39:15.544','2026-03-13 11:39:15.544','',0,0,NULL,'','','','',4,3921,40960,0,1,30,300,60,1,'auto',10,5,1,1,1,1,1,1,5,26500,65535,26535,'nat_ipv4_ipv6',300,300,1000,1000,0,1048576,0,'2026-04-01 00:00:00.000','both',1,'light',90,10,90,10,1800,10,2,2304,10240,5,1,1,'2026-03-12 21:00:11.532','2026-03-12 21:05:11.532',2,1617,6,'{\"1\":{\"max-instances\":1,\"max-resources\":{\"bandwidth\":100,\"cpu\":1,\"disk\":10240,\"memory\":512},\"max-traffic\":102400},\"2\":{\"max-instances\":3,\"max-resources\":{\"bandwidth\":200,\"cpu\":2,\"disk\":20480,\"memory\":1024},\"max-traffic\":204800},\"3\":{\"max-instances\":5,\"max-resources\":{\"bandwidth\":500,\"cpu\":4,\"disk\":40960,\"memory\":2048},\"max-traffic\":307200},\"4\":{\"max-instances\":10,\"max-resources\":{\"bandwidth\":1000,\"cpu\":8,\"disk\":81920,\"memory\":4096},\"max-traffic\":409600},\"5\":{\"max-instances\":20,\"max-resources\":{\"bandwidth\":2000,\"cpu\":16,\"disk\":163840,\"memory\":8192},\"max-traffic\":512000}}','',0,0,1,'100%',1,0,'',''),(3,'55a93e3d-1a04-4c77-8d10-f74989df81ce','2026-03-06 19:17:24.470','2026-03-13 11:39:17.364','heyhk3','proxmox','156.233.233.6','',22,'root','thanks12A#','','','','active','hk','中国香港','HK','','9.1.6',1,1,'',1,'iptables','native',0,0,'amd64','2029-03-26 08:00:00.000',0,'local','','','','','','unknown','online','2026-03-13 11:39:15.545','2026-03-13 11:39:15.545','',0,0,NULL,'','','','',8,9987,80896,0,1,30,300,60,1,'auto',20,10,1,1,1,1,1,1,5,37800,65535,37820,'nat_ipv4_ipv6',300,300,1000,1000,0,1048576,0,'2026-04-01 00:00:00.000','both',1,'light',90,10,90,10,1800,10,3,5888,8192,3,1,1,'2026-03-12 21:00:11.537','2026-03-12 21:05:11.537',5,4099,4,'{\"1\":{\"max-instances\":1,\"max-resources\":{\"bandwidth\":100,\"cpu\":1,\"disk\":10240,\"memory\":512},\"max-traffic\":102400},\"2\":{\"max-instances\":3,\"max-resources\":{\"bandwidth\":200,\"cpu\":2,\"disk\":20480,\"memory\":1024},\"max-traffic\":204800},\"3\":{\"max-instances\":5,\"max-resources\":{\"bandwidth\":500,\"cpu\":4,\"disk\":40960,\"memory\":2048},\"max-traffic\":307200},\"4\":{\"max-instances\":10,\"max-resources\":{\"bandwidth\":1000,\"cpu\":8,\"disk\":81920,\"memory\":4096},\"max-traffic\":409600},\"5\":{\"max-instances\":20,\"max-resources\":{\"bandwidth\":2000,\"cpu\":16,\"disk\":163840,\"memory\":8192},\"max-traffic\":512000}}','',0,0,1,'100%',1,0,'','');
/*!40000 ALTER TABLE `providers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `instances`
--

DROP TABLE IF EXISTS `instances`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `instances` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `provider` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `provider_id` bigint unsigned NOT NULL,
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `image` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_type` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'container',
  `cpu` bigint DEFAULT '1',
  `memory` bigint DEFAULT '512',
  `disk` bigint DEFAULT '10240',
  `bandwidth` bigint DEFAULT '10',
  `network` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `private_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `public_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ipv6_address` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `public_ipv6` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ssh_port` bigint DEFAULT '22',
  `port_range_start` bigint DEFAULT NULL,
  `port_range_end` bigint DEFAULT NULL,
  `username` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `os_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `region` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `max_traffic` bigint DEFAULT '0',
  `traffic_limited` tinyint(1) DEFAULT '0',
  `traffic_limit_reason` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `pmacct_interface_v4` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pmacct_interface_v6` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `expired_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_instances_uuid` (`uuid`),
  UNIQUE KEY `idx_instance_name_provider` (`name`,`provider_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_provider_name` (`provider`),
  KEY `idx_provider_id` (`provider_id`),
  KEY `idx_provider_status` (`provider_id`,`status`),
  KEY `idx_status` (`status`),
  KEY `idx_instance_type` (`instance_type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_user_status` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `instances`
--

LOCK TABLES `instances` WRITE;
/*!40000 ALTER TABLE `instances` DISABLE KEYS */;
INSERT INTO `instances` VALUES (1,'225e25b2-d989-4957-ad4d-1ba6a8297d57','2026-03-06 16:42:26.790','2026-03-10 15:11:02.136',NULL,'heyun-388c','heyun',1,'running','alpine-3.23-64_cloud','container',1,64,200,1,'','172.16.1.16','154.12.84.134','2001:0470:1f04:010a:0100:0000:0000:114','2001:0470:1f04:010a:0100:0000:0000:114',22,0,0,'root','j9vopl9ijtsa','alpine','',0,0,'','veth114i0','','2026-04-05 16:50:43.199',2),(2,'71be3971-d7ef-465f-93bf-5ab2a4c3921c','2026-03-06 19:36:15.357','2026-03-07 11:39:15.190','2026-03-07 11:39:15.191','heyunus-d9c8-old-1772854755','heyunus',2,'resetting','debian13','vm',3,1536,3072,3,'','172.16.1.9','38.165.47.49','2001:0470:1f04:020e:0100:0000:0000:107','2001:0470:1f04:020e:0100:0000:0000:107',22,0,0,'root','d0bd6s5dco3r','debian','',0,0,'','tap107i0','','2026-04-06 20:11:59.037',4),(3,'04776351-1eda-4a85-b5e2-8adc03bcd95e','2026-03-06 23:03:00.356','2026-03-11 01:05:22.524',NULL,'heyunus-c85c','heyunus',2,'running','alpine-3.23-64_cloud','container',1,128,512,2,'','172.16.1.10','38.165.47.49','2001:0470:1f04:020e:0100:0000:0000:108','2001:0470:1f04:020e:0100:0000:0000:108',22,0,0,'root','y3aql2a59c57','alpine','',0,0,'','veth108i0','','2026-04-05 16:50:43.199',2),(4,'b8512407-5e16-4fd8-b462-a75b2134601f','2026-03-07 10:13:34.683','2026-03-10 22:13:25.969',NULL,'heyhk3-c3ae','heyhk3',3,'running','debian13','vm',3,1536,3072,300,'','172.16.1.14','156.233.233.6','2001:0470:1f0e:0380:0100:0000:0000:112','2001:0470:1f0e:0380:0100:0000:0000:112',22,0,0,'root','va0qmomd80jg','debian','',0,0,'','tap112i0','','2026-05-06 20:11:59.037',4),(5,'6364f5f1-0dec-437c-9d6d-a0e479616b37','2026-03-07 10:38:44.683','2026-03-07 10:41:53.352',NULL,'heyhk3-f9e0','heyhk3',3,'running','alpine-3.23-64_cloud','container',1,128,512,100,'','172.16.1.15','156.233.233.6','2001:0470:1f0e:0380:0100:0000:0000:113','2001:0470:1f0e:0380:0100:0000:0000:113',22,0,0,'root','pfagg0a8xii2','alpine','',0,0,'','veth113i0','','2026-04-07 10:36:40.085',5),(6,'9a3eb51f-2de8-4f70-8dd3-8602a17a0d58','2026-03-07 11:39:15.191','2026-03-09 11:52:25.946','2026-03-09 11:53:59.901','heyunus-d9c8','heyunus',2,'deleting','debian13','vm',3,1536,3072,3,'','','38.165.47.49','','',22,0,0,'','','debian','',0,0,'','','','2026-04-07 10:36:40.085',5),(7,'ef55469a-a554-42f4-9e8b-4af378d7c7a9','2026-03-07 19:31:00.192','2026-03-07 19:34:16.998',NULL,'heyhk3-4a60','heyhk3',3,'running','alpine-3.23-64_cloud','container',4,4096,4096,500,'','172.16.1.16','156.233.233.6','2001:0470:1f0e:0380:0100:0000:0000:114','2001:0470:1f0e:0380:0100:0000:0000:114',22,0,0,'root','bwa7gvdsu1ma','alpine','',0,0,'','veth114i0','','2026-05-07 19:28:02.051',6),(8,'627da961-09e2-4808-8d24-2340712d4989','2026-03-07 20:11:25.189','2026-03-07 20:14:31.150',NULL,'heyunus-2441','heyunus',2,'running','alpine-3.23-64_cloud','container',1,128,512,100,'','172.16.1.11','38.165.47.49','2001:0470:1f04:020e:0100:0000:0000:109','2001:0470:1f04:020e:0100:0000:0000:109',22,0,0,'root','wgd2hzrwuz6k','alpine','',0,0,'','veth109i0','','2026-04-07 20:05:50.901',7),(9,'cbc29b45-2a02-4d09-add8-fb068c5a830c','2026-03-10 23:31:55.190','2026-03-11 01:14:30.280',NULL,'heyun-1808','heyun',1,'running','debian13','vm',2,1024,3072,200,'','172.16.1.18','154.12.84.134','2001:0470:1f04:010a:0100:0000:0000:116','2001:0470:1f04:010a:0100:0000:0000:116',22,0,0,'root','ue4aiy8grmms','debian','',0,0,'','tap116i0','','2026-04-05 16:50:43.199',2),(10,'2c88fd61-d536-4679-8329-1a28f4431ca6','2026-03-11 00:47:05.189','2026-03-11 00:56:39.691',NULL,'heyunus-f9d6','heyunus',2,'running','alpine-3.23-64_cloud','container',2,512,2048,200,'','172.16.1.4','38.165.47.49','2001:0470:1f04:020e:0100:0000:0000:102','2001:0470:1f04:020e:0100:0000:0000:102',22,0,0,'root','o59y9dp9ld0n','alpine','',0,0,'','veth102i0','','2026-04-07 10:36:40.085',5),(11,'8330a40e-cdf7-4701-b526-25b096f98642','2026-03-11 00:57:45.188','2026-03-11 01:03:12.924',NULL,'heyunus-1f58','heyunus',2,'running','debian13','vm',2,512,3072,200,'','172.16.1.9','38.165.47.49','2001:0470:1f04:020e:0100:0000:0000:107','2001:0470:1f04:020e:0100:0000:0000:107',22,0,0,'root','ix5j2xk79kid','debian','',0,0,'','tap107i0','','2026-05-06 20:11:59.037',4),(12,'226ab4b4-13da-4a87-9c2d-d93f5a4afc61','2026-03-11 01:06:20.188','2026-03-11 01:12:23.600',NULL,'heyunus-a9cd','heyunus',2,'running','debian-11-64_cloud','container',1,512,2048,200,'','172.16.1.12','38.165.47.49','2001:0470:1f04:020e:0100:0000:0000:110','2001:0470:1f04:020e:0100:0000:0000:110',22,0,0,'root','b1luee11c6t8','debian','',0,0,'','veth110i0','','2026-05-06 20:11:59.037',4),(13,'945d2eea-6e68-4c70-80c3-df8286e95b5a','2026-03-11 01:15:25.188','2026-03-11 01:21:24.078',NULL,'heyunus-930f','heyunus',2,'running','debian-11-64_cloud','container',1,512,2048,200,'','172.16.1.13','38.165.47.49','2001:0470:1f04:020e:0100:0000:0000:111','2001:0470:1f04:020e:0100:0000:0000:111',22,0,0,'root','ekfq06xc1y3g','debian','',0,0,'','veth111i0','','2026-05-06 20:11:59.037',4),(14,'fa1535db-5803-408c-a804-ee9c7f5f659b','2026-03-11 22:44:26.259','2026-03-11 22:49:02.212',NULL,'heyun-8550','heyun',1,'running','debian-11-64_cloud','container',1,512,2048,100,'','172.16.1.19','154.12.84.134','2001:0470:1f04:010a:0100:0000:0000:117','2001:0470:1f04:010a:0100:0000:0000:117',22,0,0,'root','p3mzy8slkh4q','debian','',0,0,'','veth117i0','','2026-06-04 16:41:56.325',1),(15,'85eb8bd2-fc54-4da2-84b2-ff6bac1b48c2','2026-03-12 12:09:01.262','2026-03-12 12:11:38.155',NULL,'heyhk3-cb16','heyhk3',3,'running','alpine-3.23-64_cloud','container',1,128,512,100,'','172.16.1.2','156.233.233.6','2001:0470:1f0e:0380:0100:0000:0000:100','2001:0470:1f0e:0380:0100:0000:0000:100',22,0,0,'root','xs7ddwtha0ol','alpine','',0,0,'','veth100i0','','2026-04-11 11:45:29.280',10);
/*!40000 ALTER TABLE `instances` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ports`
--

DROP TABLE IF EXISTS `ports`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ports` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `instance_id` bigint unsigned DEFAULT NULL,
  `provider_id` bigint unsigned DEFAULT NULL,
  `host_port` bigint NOT NULL,
  `host_port_end` bigint DEFAULT '0',
  `guest_port` bigint NOT NULL,
  `guest_port_end` bigint DEFAULT '0',
  `port_count` bigint DEFAULT '1',
  `protocol` varchar(8) COLLATE utf8mb4_unicode_ci DEFAULT 'both',
  `status` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  `description` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_ssh` tinyint(1) DEFAULT '0',
  `is_automatic` tinyint(1) DEFAULT '1',
  `port_type` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'range_mapped',
  `ipv6_enabled` tinyint(1) DEFAULT '0',
  `ipv6_address` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mapping_method` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'native',
  PRIMARY KEY (`id`),
  KEY `idx_ports_deleted_at` (`deleted_at`),
  KEY `idx_instance_ssh` (`instance_id`,`is_ssh`),
  KEY `idx_instance_status` (`instance_id`,`status`),
  KEY `idx_provider_id` (`provider_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ports`
--

LOCK TABLES `ports` WRITE;
/*!40000 ALTER TABLE `ports` DISABLE KEYS */;
INSERT INTO `ports` VALUES (1,'2026-03-06 16:42:33.235','2026-03-06 16:42:33.235',NULL,1,1,33300,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(2,'2026-03-06 16:42:33.237','2026-03-06 16:42:33.237',NULL,1,1,33301,0,33301,0,1,'both','active','端口33301',0,1,'range_mapped',1,'','native'),(3,'2026-03-06 16:42:33.237','2026-03-06 16:42:33.237',NULL,1,1,33302,0,33302,0,1,'both','active','端口33302',0,1,'range_mapped',1,'','native'),(4,'2026-03-06 16:42:33.237','2026-03-06 16:42:33.237',NULL,1,1,33303,0,33303,0,1,'both','active','端口33303',0,1,'range_mapped',1,'','native'),(5,'2026-03-06 16:42:33.237','2026-03-06 16:42:33.237',NULL,1,1,33304,0,33304,0,1,'both','active','端口33304',0,1,'range_mapped',1,'','native'),(6,'2026-03-06 19:36:18.811','2026-03-06 19:36:18.811',NULL,2,2,26500,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(7,'2026-03-06 19:36:18.812','2026-03-06 19:36:18.812',NULL,2,2,26501,0,26501,0,1,'both','active','端口26501',0,1,'range_mapped',1,'','native'),(8,'2026-03-06 19:36:18.812','2026-03-06 19:36:18.812',NULL,2,2,26502,0,26502,0,1,'both','active','端口26502',0,1,'range_mapped',1,'','native'),(9,'2026-03-06 19:36:18.812','2026-03-06 19:36:18.812',NULL,2,2,26503,0,26503,0,1,'both','active','端口26503',0,1,'range_mapped',1,'','native'),(10,'2026-03-06 19:36:18.812','2026-03-06 19:36:18.812',NULL,2,2,26504,0,26504,0,1,'both','active','端口26504',0,1,'range_mapped',1,'','native'),(11,'2026-03-06 23:03:00.641','2026-03-06 23:03:00.641',NULL,3,2,26505,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(12,'2026-03-06 23:03:00.642','2026-03-06 23:03:00.642',NULL,3,2,26506,0,26506,0,1,'both','active','端口26506',0,1,'range_mapped',1,'','native'),(13,'2026-03-06 23:03:00.642','2026-03-06 23:03:00.642',NULL,3,2,26507,0,26507,0,1,'both','active','端口26507',0,1,'range_mapped',1,'','native'),(14,'2026-03-06 23:03:00.642','2026-03-06 23:03:00.642',NULL,3,2,26508,0,26508,0,1,'both','active','端口26508',0,1,'range_mapped',1,'','native'),(15,'2026-03-06 23:03:00.642','2026-03-06 23:03:00.642',NULL,3,2,26509,0,26509,0,1,'both','active','端口26509',0,1,'range_mapped',1,'','native'),(16,'2026-03-07 10:13:40.218','2026-03-07 10:13:40.218',NULL,4,3,37800,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(17,'2026-03-07 10:13:40.218','2026-03-07 10:13:40.218',NULL,4,3,37801,0,37801,0,1,'both','active','端口37801',0,1,'range_mapped',1,'','native'),(18,'2026-03-07 10:13:40.218','2026-03-07 10:13:40.218',NULL,4,3,37802,0,37802,0,1,'both','active','端口37802',0,1,'range_mapped',1,'','native'),(19,'2026-03-07 10:13:40.218','2026-03-07 10:13:40.218',NULL,4,3,37803,0,37803,0,1,'both','active','端口37803',0,1,'range_mapped',1,'','native'),(20,'2026-03-07 10:13:40.218','2026-03-07 10:13:40.218',NULL,4,3,37804,0,37804,0,1,'both','active','端口37804',0,1,'range_mapped',1,'','native'),(21,'2026-03-07 10:38:46.931','2026-03-07 10:38:46.931',NULL,5,3,37805,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(22,'2026-03-07 10:38:46.932','2026-03-07 10:38:46.932',NULL,5,3,37806,0,37806,0,1,'both','active','端口37806',0,1,'range_mapped',1,'','native'),(23,'2026-03-07 10:38:46.932','2026-03-07 10:38:46.932',NULL,5,3,37807,0,37807,0,1,'both','active','端口37807',0,1,'range_mapped',1,'','native'),(24,'2026-03-07 10:38:46.932','2026-03-07 10:38:46.932',NULL,5,3,37808,0,37808,0,1,'both','active','端口37808',0,1,'range_mapped',1,'','native'),(25,'2026-03-07 10:38:46.932','2026-03-07 10:38:46.932',NULL,5,3,37809,0,37809,0,1,'both','active','端口37809',0,1,'range_mapped',1,'','native'),(26,'2026-03-07 19:31:05.929','2026-03-07 19:31:05.929',NULL,7,3,37810,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(27,'2026-03-07 19:31:05.931','2026-03-07 19:31:05.931',NULL,7,3,37811,0,37811,0,1,'both','active','端口37811',0,1,'range_mapped',1,'','native'),(28,'2026-03-07 19:31:05.931','2026-03-07 19:31:05.931',NULL,7,3,37812,0,37812,0,1,'both','active','端口37812',0,1,'range_mapped',1,'','native'),(29,'2026-03-07 19:31:05.931','2026-03-07 19:31:05.931',NULL,7,3,37813,0,37813,0,1,'both','active','端口37813',0,1,'range_mapped',1,'','native'),(30,'2026-03-07 19:31:05.931','2026-03-07 19:31:05.931',NULL,7,3,37814,0,37814,0,1,'both','active','端口37814',0,1,'range_mapped',1,'','native'),(31,'2026-03-07 20:11:25.472','2026-03-07 20:11:25.472',NULL,8,2,26510,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(32,'2026-03-07 20:11:25.473','2026-03-07 20:11:25.473',NULL,8,2,26511,0,26511,0,1,'both','active','端口26511',0,1,'range_mapped',1,'','native'),(33,'2026-03-07 20:11:25.473','2026-03-07 20:11:25.473',NULL,8,2,26512,0,26512,0,1,'both','active','端口26512',0,1,'range_mapped',1,'','native'),(34,'2026-03-07 20:11:25.473','2026-03-07 20:11:25.473',NULL,8,2,26513,0,26513,0,1,'both','active','端口26513',0,1,'range_mapped',1,'','native'),(35,'2026-03-07 20:11:25.473','2026-03-07 20:11:25.473',NULL,8,2,26514,0,26514,0,1,'both','active','端口26514',0,1,'range_mapped',1,'','native'),(36,'2026-03-10 23:31:57.550','2026-03-10 23:31:57.550',NULL,9,1,33305,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(37,'2026-03-10 23:31:57.551','2026-03-10 23:31:57.551',NULL,9,1,33306,0,33306,0,1,'both','active','端口33306',0,1,'range_mapped',1,'','native'),(38,'2026-03-10 23:31:57.551','2026-03-10 23:31:57.551',NULL,9,1,33307,0,33307,0,1,'both','active','端口33307',0,1,'range_mapped',1,'','native'),(39,'2026-03-10 23:31:57.551','2026-03-10 23:31:57.551',NULL,9,1,33308,0,33308,0,1,'both','active','端口33308',0,1,'range_mapped',1,'','native'),(40,'2026-03-10 23:31:57.551','2026-03-10 23:31:57.551',NULL,9,1,33309,0,33309,0,1,'both','active','端口33309',0,1,'range_mapped',1,'','native'),(41,'2026-03-11 00:47:05.530','2026-03-11 00:47:05.530',NULL,10,2,26515,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(42,'2026-03-11 00:47:05.531','2026-03-11 00:47:05.531',NULL,10,2,26516,0,26516,0,1,'both','active','端口26516',0,1,'range_mapped',1,'','native'),(43,'2026-03-11 00:47:05.531','2026-03-11 00:47:05.531',NULL,10,2,26517,0,26517,0,1,'both','active','端口26517',0,1,'range_mapped',1,'','native'),(44,'2026-03-11 00:47:05.531','2026-03-11 00:47:05.531',NULL,10,2,26518,0,26518,0,1,'both','active','端口26518',0,1,'range_mapped',1,'','native'),(45,'2026-03-11 00:47:05.531','2026-03-11 00:47:05.531',NULL,10,2,26519,0,26519,0,1,'both','active','端口26519',0,1,'range_mapped',1,'','native'),(46,'2026-03-11 00:57:45.528','2026-03-11 00:57:45.528',NULL,11,2,26520,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(47,'2026-03-11 00:57:45.529','2026-03-11 00:57:45.529',NULL,11,2,26521,0,26521,0,1,'both','active','端口26521',0,1,'range_mapped',1,'','native'),(48,'2026-03-11 00:57:45.529','2026-03-11 00:57:45.529',NULL,11,2,26522,0,26522,0,1,'both','active','端口26522',0,1,'range_mapped',1,'','native'),(49,'2026-03-11 00:57:45.529','2026-03-11 00:57:45.529',NULL,11,2,26523,0,26523,0,1,'both','active','端口26523',0,1,'range_mapped',1,'','native'),(50,'2026-03-11 00:57:45.529','2026-03-11 00:57:45.529',NULL,11,2,26524,0,26524,0,1,'both','active','端口26524',0,1,'range_mapped',1,'','native'),(51,'2026-03-11 01:06:20.496','2026-03-11 01:06:20.496',NULL,12,2,26525,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(52,'2026-03-11 01:06:20.497','2026-03-11 01:06:20.497',NULL,12,2,26526,0,26526,0,1,'both','active','端口26526',0,1,'range_mapped',1,'','native'),(53,'2026-03-11 01:06:20.497','2026-03-11 01:06:20.497',NULL,12,2,26527,0,26527,0,1,'both','active','端口26527',0,1,'range_mapped',1,'','native'),(54,'2026-03-11 01:06:20.497','2026-03-11 01:06:20.497',NULL,12,2,26528,0,26528,0,1,'both','active','端口26528',0,1,'range_mapped',1,'','native'),(55,'2026-03-11 01:06:20.497','2026-03-11 01:06:20.497',NULL,12,2,26529,0,26529,0,1,'both','active','端口26529',0,1,'range_mapped',1,'','native'),(56,'2026-03-11 01:15:25.536','2026-03-11 01:15:25.536',NULL,13,2,26530,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(57,'2026-03-11 01:15:25.537','2026-03-11 01:15:25.537',NULL,13,2,26531,0,26531,0,1,'both','active','端口26531',0,1,'range_mapped',1,'','native'),(58,'2026-03-11 01:15:25.537','2026-03-11 01:15:25.537',NULL,13,2,26532,0,26532,0,1,'both','active','端口26532',0,1,'range_mapped',1,'','native'),(59,'2026-03-11 01:15:25.537','2026-03-11 01:15:25.537',NULL,13,2,26533,0,26533,0,1,'both','active','端口26533',0,1,'range_mapped',1,'','native'),(60,'2026-03-11 01:15:25.537','2026-03-11 01:15:25.537',NULL,13,2,26534,0,26534,0,1,'both','active','端口26534',0,1,'range_mapped',1,'','native'),(61,'2026-03-11 22:44:33.164','2026-03-11 22:44:33.164',NULL,14,1,33310,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(62,'2026-03-11 22:44:33.165','2026-03-11 22:44:33.165',NULL,14,1,33311,0,33311,0,1,'both','active','端口33311',0,1,'range_mapped',1,'','native'),(63,'2026-03-11 22:44:33.165','2026-03-11 22:44:33.165',NULL,14,1,33312,0,33312,0,1,'both','active','端口33312',0,1,'range_mapped',1,'','native'),(64,'2026-03-11 22:44:33.165','2026-03-11 22:44:33.165',NULL,14,1,33313,0,33313,0,1,'both','active','端口33313',0,1,'range_mapped',1,'','native'),(65,'2026-03-11 22:44:33.165','2026-03-11 22:44:33.165',NULL,14,1,33314,0,33314,0,1,'both','active','端口33314',0,1,'range_mapped',1,'','native'),(66,'2026-03-12 12:09:07.139','2026-03-12 12:09:07.139',NULL,15,3,37815,0,22,0,1,'both','active','SSH',1,1,'range_mapped',1,'','native'),(67,'2026-03-12 12:09:07.140','2026-03-12 12:09:07.140',NULL,15,3,37816,0,37816,0,1,'both','active','端口37816',0,1,'range_mapped',1,'','native'),(68,'2026-03-12 12:09:07.140','2026-03-12 12:09:07.140',NULL,15,3,37817,0,37817,0,1,'both','active','端口37817',0,1,'range_mapped',1,'','native'),(69,'2026-03-12 12:09:07.140','2026-03-12 12:09:07.140',NULL,15,3,37818,0,37818,0,1,'both','active','端口37818',0,1,'range_mapped',1,'','native'),(70,'2026-03-12 12:09:07.140','2026-03-12 12:09:07.140',NULL,15,3,37819,0,37819,0,1,'both','active','端口37819',0,1,'range_mapped',1,'','native');
/*!40000 ALTER TABLE `ports` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单号',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `product_id` bigint unsigned DEFAULT NULL COMMENT '产品ID',
  `amount` bigint NOT NULL COMMENT '订单金额(分)',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单状态',
  `payment_method` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '支付方式',
  `payment_time` datetime(3) DEFAULT NULL COMMENT '支付时间',
  `paid_amount` bigint DEFAULT '0' COMMENT '实付金额',
  `product_data` json DEFAULT NULL COMMENT '产品快照',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `expire_at` datetime(3) DEFAULT NULL COMMENT '订单过期时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_orders_order_no` (`order_no`),
  KEY `idx_user_order` (`user_id`),
  KEY `idx_orders_product_id` (`product_id`),
  KEY `idx_orders_status` (`status`),
  CONSTRAINT `fk_orders_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
INSERT INTO `orders` VALUES (1,'1772785321857924',1,NULL,100,'pending','epay',NULL,0,'{}','','2026-03-06 16:52:01.692','2026-03-06 16:22:01.692','2026-03-06 16:22:01.692'),(2,'1772785343678722',1,NULL,100,'pending','epay',NULL,0,'{}','','2026-03-06 16:52:23.968','2026-03-06 16:22:23.968','2026-03-06 16:22:23.968'),(3,'1772785415117379',1,1,100,'pending','epay',NULL,0,'{\"id\": 1, \"cpu\": 1, \"disk\": 512, \"name\": \"入门套餐\", \"level\": 1, \"price\": 100, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 2, \"description\": \"适合个人用户和小型项目使用\", \"maxInstances\": 1}','','2026-03-06 16:53:35.512','2026-03-06 16:23:35.512','2026-03-06 16:23:35.512'),(4,'1772786091598467',1,NULL,100,'pending','epay',NULL,0,'{}','','2026-03-06 17:04:51.085','2026-03-06 16:34:51.085','2026-03-06 16:34:51.085'),(5,'1772786303356312',1,1,100,'pending','epay',NULL,0,'{\"id\": 1, \"cpu\": 1, \"disk\": 512, \"name\": \"入门套餐\", \"level\": 1, \"price\": 100, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 2, \"description\": \"适合个人用户和小型项目使用\", \"maxInstances\": 1}','','2026-03-06 17:08:23.548','2026-03-06 16:38:23.548','2026-03-06 16:38:23.548'),(6,'1772786516066954',1,7,0,'paid','balance','2026-03-06 16:41:56.326',0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 2, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-06 17:11:56.325','2026-03-06 16:41:56.326','2026-03-06 16:41:56.326'),(7,'1772787043661535',2,7,0,'paid','balance','2026-03-06 16:50:43.200',0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 2, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-06 17:20:43.199','2026-03-06 16:50:43.200','2026-03-06 16:50:43.200'),(8,'1772788160423478',1,NULL,100,'pending','epay',NULL,0,'{}','','2026-03-06 17:39:20.995','2026-03-06 17:09:20.996','2026-03-06 17:09:20.996'),(9,'1772796507751942',3,2,200,'paid','balance','2026-03-06 19:28:27.945',200,'{\"id\": 2, \"cpu\": 2, \"disk\": 2048, \"name\": \"中级套餐\", \"level\": 2, \"price\": 200, \"memory\": 1024, \"period\": 30, \"traffic\": 204800, \"bandwidth\": 2, \"description\": \"适合小型团队和中型项目使用\", \"maxInstances\": 1}','','2026-03-06 19:58:27.944','2026-03-06 19:28:27.945','2026-03-06 19:28:27.945'),(10,'1772799287466059',4,7,0,'pending','epay',NULL,0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 2, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-06 20:44:47.277','2026-03-06 20:14:47.278','2026-03-06 20:14:47.278'),(11,'1772799299497544',4,1,100,'pending','alipay',NULL,0,'{\"id\": 1, \"cpu\": 1, \"disk\": 512, \"name\": \"入门套餐\", \"level\": 1, \"price\": 100, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 2, \"description\": \"适合个人用户和小型项目使用\", \"maxInstances\": 1}','','2026-03-06 20:44:59.050','2026-03-06 20:14:59.050','2026-03-06 20:14:59.050'),(12,'1772799322273808',4,1,100,'expired','epay',NULL,0,'{\"id\": 1, \"cpu\": 1, \"disk\": 512, \"name\": \"入门套餐\", \"level\": 1, \"price\": 100, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 2, \"description\": \"适合个人用户和小型项目使用\", \"maxInstances\": 1}','','2026-03-06 20:45:22.248','2026-03-06 20:15:22.248','2026-03-06 21:18:17.092'),(13,'1772809346576203',4,7,0,'paid','balance','2026-03-06 23:02:26.305',0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 2, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-06 23:32:26.304','2026-03-06 23:02:26.305','2026-03-06 23:02:26.305'),(14,'1772847371510739',1,NULL,100,'paid','epay','2026-03-07 09:56:18.201',100,'{}','','2026-03-07 10:06:11.379','2026-03-07 09:36:11.379','2026-03-07 09:56:18.201'),(15,'1772848856204104',1,NULL,200,'paid','epay','2026-03-07 10:01:27.725',200,'{}','','2026-03-07 10:30:56.569','2026-03-07 10:00:56.569','2026-03-07 10:01:27.725'),(16,'1772849126504258',1,2,200,'paid','epay','2026-03-07 10:05:41.829',200,'{\"id\": 2, \"cpu\": 2, \"disk\": 2048, \"name\": \"中级套餐\", \"level\": 2, \"price\": 200, \"memory\": 1024, \"period\": 30, \"traffic\": 204800, \"bandwidth\": 200, \"description\": \"适合小型团队和中型项目使用\", \"maxInstances\": 1}','','2026-03-07 10:35:26.459','2026-03-07 10:05:26.459','2026-03-07 10:05:41.829'),(17,'1772849469490162',1,3,300,'paid','balance','2026-03-07 10:11:09.768',300,'{\"id\": 3, \"cpu\": 3, \"disk\": 3072, \"name\": \"高级套餐\", \"level\": 3, \"price\": 300, \"memory\": 1536, \"period\": 30, \"traffic\": 307200, \"bandwidth\": 300, \"description\": \"适合中型团队和大型项目使用\", \"maxInstances\": 1}','','2026-03-07 10:41:09.767','2026-03-07 10:11:09.768','2026-03-07 10:11:09.768'),(18,'1772851078828917',5,7,0,'paid','balance','2026-03-07 10:37:58.517',0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 100, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-07 11:07:58.516','2026-03-07 10:37:58.517','2026-03-07 10:37:58.517'),(19,'1772882918256081',6,7,0,'pending','alipay',NULL,0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 100, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-07 19:58:38.566','2026-03-07 19:28:38.566','2026-03-07 19:28:38.566'),(20,'1772882928131265',6,7,0,'paid','balance','2026-03-07 19:28:48.205',0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 100, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-07 19:58:48.199','2026-03-07 19:28:48.205','2026-03-07 19:28:48.205'),(21,'1772882973316931',6,4,500,'paid','epay','2026-03-07 19:29:49.661',500,'{\"id\": 4, \"cpu\": 4, \"disk\": 4096, \"name\": \"超级套餐\", \"level\": 4, \"price\": 500, \"memory\": 4096, \"period\": 30, \"traffic\": 409600, \"bandwidth\": 500, \"description\": \"适合大型团队和企业级项目使用\", \"maxInstances\": 1}','','2026-03-07 19:59:33.234','2026-03-07 19:29:33.234','2026-03-07 19:29:49.661'),(22,'1772885266175316',7,7,0,'pending','alipay',NULL,0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 100, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-07 20:37:46.580','2026-03-07 20:07:46.580','2026-03-07 20:07:46.580'),(23,'1772885416303058',7,7,0,'paid','balance','2026-03-07 20:10:16.547',0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 100, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-07 20:40:16.544','2026-03-07 20:10:16.547','2026-03-07 20:10:16.547'),(24,'1772899538699987',8,7,0,'pending','alipay',NULL,0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 100, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-08 00:35:38.189','2026-03-08 00:05:38.189','2026-03-08 00:05:38.189'),(25,'1772899605027596',8,7,0,'pending','alipay',NULL,0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 100, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-08 00:36:45.846','2026-03-08 00:06:45.846','2026-03-08 00:06:45.846'),(26,'1773151928836287',4,2,300,'paid','epay','2026-03-10 22:13:25.967',300,'{\"id\": 2, \"cpu\": 2, \"disk\": 2048, \"name\": \"中级套餐\", \"level\": 2, \"price\": 300, \"memory\": 1024, \"period\": 30, \"traffic\": 204800, \"bandwidth\": 200, \"description\": \"适合小型团队和中型项目使用\", \"maxInstances\": 1}','','2026-03-10 22:42:08.585','2026-03-10 22:12:08.585','2026-03-10 22:13:25.967'),(27,'1773216332224365',9,NULL,100,'pending','epay',NULL,0,'{}','','2026-03-11 16:35:32.741','2026-03-11 16:05:32.741','2026-03-11 16:05:32.741'),(28,'1773216352464342',9,1,9900,'pending','epay',NULL,0,'{\"id\": 1, \"cpu\": 1, \"disk\": 1024, \"name\": \"入门级套餐\", \"level\": 1, \"price\": 9900, \"memory\": 350, \"period\": 30, \"traffic\": 102400, \"bandwidth\": 100, \"description\": \"适合个人用户和小型项目使用\", \"maxInstances\": 1}','','2026-03-11 16:35:52.243','2026-03-11 16:05:52.243','2026-03-11 16:05:52.243'),(29,'1773287129769996',10,7,0,'paid','balance','2026-03-12 11:45:29.283',0,'{\"id\": 7, \"cpu\": 1, \"disk\": 512, \"name\": \"免费套餐\", \"level\": 1, \"price\": 0, \"memory\": 128, \"period\": 30, \"traffic\": 1024, \"bandwidth\": 100, \"description\": \"免费套餐\", \"maxInstances\": 1}','','2026-03-12 12:15:29.280','2026-03-12 11:45:29.283','2026-03-12 11:45:29.283');
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_wallets`
--

DROP TABLE IF EXISTS `user_wallets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_wallets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `balance` bigint NOT NULL DEFAULT '0' COMMENT '余额(分)',
  `frozen` bigint NOT NULL DEFAULT '0' COMMENT '冻结金额',
  `total_recharge` bigint NOT NULL DEFAULT '0' COMMENT '累计充值',
  `total_expense` bigint NOT NULL DEFAULT '0' COMMENT '累计消费',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_wallets_user_id` (`user_id`),
  KEY `idx_user_wallet` (`user_id`),
  CONSTRAINT `fk_user_wallets_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_wallets`
--

LOCK TABLES `user_wallets` WRITE;
/*!40000 ALTER TABLE `user_wallets` DISABLE KEYS */;
INSERT INTO `user_wallets` VALUES (1,1,0,0,300,300,'2026-03-06 16:21:54.841','2026-03-07 10:11:09.767'),(2,2,0,0,0,0,'2026-03-06 16:50:43.199','2026-03-06 16:50:43.200'),(3,3,0,0,200,200,'2026-03-06 19:23:26.686','2026-03-06 19:28:27.944'),(4,4,100,0,100,0,'2026-03-06 23:01:48.898','2026-03-06 23:02:26.304'),(5,5,0,0,0,0,'2026-03-07 10:37:06.366','2026-03-07 10:37:58.517'),(6,6,0,0,0,0,'2026-03-07 19:28:48.201','2026-03-07 19:28:48.203'),(7,7,0,0,0,0,'2026-03-07 20:10:16.545','2026-03-07 20:10:16.546'),(8,9,0,0,0,0,'2026-03-11 16:03:40.128','2026-03-11 16:03:40.128'),(9,10,0,0,0,0,'2026-03-12 11:45:29.281','2026-03-12 11:45:29.282');
/*!40000 ALTER TABLE `user_wallets` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-03-13 16:33:18
