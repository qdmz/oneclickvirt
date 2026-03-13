-- 简化的迁移脚本
-- 只导入必要的字段

-- 先删除新数据库中可能存在的旧数据，避免冲突
DELETE FROM ports;
DELETE FROM instances;
DELETE FROM tasks;
DELETE FROM providers;
DELETE FROM users;

-- 导入用户数据（只导入必要字段）
INSERT INTO `users` (`id`, `uuid`, `username`, `email`, `password`, `user_type`, `level`, `status`, `nickname`, `phone`, `created_at`, `updated_at`)
VALUES 
(1,'6939f1fd-a925-4d7a-8958-8ee19a9aee7a','admin','admin@ypvps.com','$2a$10$nyAOEDpaP9GH01njj4ymDuPJn4tZA0ibGQT3TQARAQyWw5PDYq/O6','admin',5,1,'admin','','2026-03-06 15:37:07.000','2026-03-12 21:08:17.339'),
(2,'61a9e805-157f-41e6-adae-d6e0bc036e2d','user','qdmz@vip.qq.com','$2a$10$mH153cHmiB.q50yDQLSlKulSG2FDojHHWUenf8hBI8wWivJlj8OXK','user',1,1,'user','','2026-03-06 15:37:07.000','2026-03-12 20:30:03.892'),
(3,'44ff4a62-46c3-4c6a-a8e9-fe8b559bac22','byzhenyu','','$2a$10$K67b4BN4.JMkoa0jsV/MEOpnJbfwNkdA3v247E7H6WB7gjZ28I1ta','user',3,0,'','','2026-03-06 19:23:14.508','2026-03-07 09:30:57.376'),
(4,'b3d86fdb-5cf4-46c0-a6b6-7c6c81656bab','399107679@qq.com','','$2a$10$hFtXgT3XpkbIxHjD2sCs1ew9t69njfWQv8NHF5blFdF7QsY4BJVoy','user',2,1,'','','2026-03-06 20:11:59.038','2026-03-10 22:13:25.970'),
(5,'cab8d3bb-cb29-4f57-b405-c39b1fa42cc5','ypvps','admin@qdmz.biters.edu.kg','$2a$10$vcLBkvhxq64RaDci9DdUXOLk1wan0U8A8qs3.Ds0Rib9o8idQJ8w.','user',1,1,'ypvps','13885852635','2026-03-07 10:36:40.085','2026-03-09 11:53:59.900');

-- 导入节点数据（只导入必要字段）
INSERT INTO `providers` (`id`, `uuid`, `name`, `type`, `endpoint`, `ssh_port`, `username`, `password`, `status`, `region`, `country`, `created_at`, `updated_at`)
VALUES 
(1,'bd0fa1e1-82c3-43ce-8e8d-25b7cfd93b66','heyun','proxmox','154.12.84.134',22,'root','thanks12A#','active','hk','中国香港','2026-03-06 16:21:13.956','2026-03-13 11:39:17.629'),
(2,'fa00a303-0753-49b3-961a-111f8d717bc6','heyunus','proxmox','38.165.47.49',22,'root','thanks12A#','active','usa','美国','2026-03-06 19:14:03.477','2026-03-13 11:39:16.364');

-- 完成迁移
SELECT '数据迁移完成' AS message;