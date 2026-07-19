-- 产品表新字段迁移脚本
-- 为 products 表添加 VM 产品相关字段
-- 执行: mysql -uroot -p [数据库名] < scripts/add_product_vm_fields.sql

SET FOREIGN_KEY_CHECKS = 0;

-- 添加 VM 类型字段
ALTER TABLE `products` 
ADD COLUMN IF NOT EXISTS `vm_type` varchar(20) NOT NULL DEFAULT '' COMMENT 'VM类型(vm/container/proxmox)' AFTER `sold_count`;

-- 添加默认操作系统字段
ALTER TABLE `products` 
ADD COLUMN IF NOT EXISTS `os` varchar(50) NOT NULL DEFAULT '' COMMENT '默认操作系统' AFTER `vm_type`;

-- 添加自动开通字段
ALTER TABLE `products` 
ADD COLUMN IF NOT EXISTS `auto_provision` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否支付后自动开通' AFTER `os`;

-- 添加最低等级字段
ALTER TABLE `products` 
ADD COLUMN IF NOT EXISTS `required_level` int NOT NULL DEFAULT 1 COMMENT '购买所需最低等级' AFTER `auto_provision`;

-- 更新现有产品数据，设置自动开通为1
UPDATE `products` SET `auto_provision` = 1 WHERE `auto_provision` IS NULL;

SET FOREIGN_KEY_CHECKS = 1;

-- 显示变更结果
SELECT id, name, vm_type, os, auto_provision, required_level FROM `products`;
