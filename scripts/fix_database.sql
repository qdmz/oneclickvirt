-- OneClickVirt 数据库修复脚本
-- 修复以下问题：
-- 1. redemption_codes 表缺少 used_count 字段
-- 2. product_purchases 表的外键约束字段类型不匹配
-- 3. system_images 表的 provider_type 和 architecture 字段需要默认值

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 修复 1: redemption_codes 表添加 used_count 字段
-- ============================================
ALTER TABLE `redemption_codes` ADD COLUMN IF NOT EXISTS `used_count` int(11) NOT NULL DEFAULT 0 COMMENT '已使用次数' AFTER `max_uses`;

-- 更新现有数据，将 uses 字段的值复制到 used_count
UPDATE `redemption_codes` SET `used_count` = `uses` WHERE `used_count` = 0;

-- ============================================
-- 修复 2: product_purchases 表的外键约束
-- 确保 user_id 和 product_id 字段类型与 users 和 products 表的 id 字段一致
-- ============================================
-- 先移除可能存在的外键约束
ALTER TABLE `product_purchases` DROP FOREIGN KEY IF EXISTS `fk_product_purchases_user`;
ALTER TABLE `product_purchases` DROP FOREIGN KEY IF EXISTS `fk_product_purchases_product`;

-- 确保 user_id 和 product_id 为 bigint unsigned 类型（与 users 和 products 表一致）
ALTER TABLE `product_purchases` MODIFY COLUMN `user_id` bigint unsigned NOT NULL COMMENT '用户ID';
ALTER TABLE `product_purchases` MODIFY COLUMN `product_id` bigint unsigned NOT NULL COMMENT '产品ID';
ALTER TABLE `product_purchases` MODIFY COLUMN `order_id` bigint unsigned DEFAULT NULL COMMENT '订单ID';

-- 添加外键约束
ALTER TABLE `product_purchases` ADD CONSTRAINT `fk_product_purchases_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`);
ALTER TABLE `product_purchases` ADD CONSTRAINT `fk_product_purchases_product` FOREIGN KEY (`product_id`) REFERENCES `products`(`id`);

-- ============================================
-- 修复 3: system_images 表添加默认值
-- ============================================
ALTER TABLE `system_images` MODIFY COLUMN `provider_type` varchar(32) NOT NULL DEFAULT 'docker' COMMENT '支持的Provider类型';
ALTER TABLE `system_images` MODIFY COLUMN `architecture` varchar(16) NOT NULL DEFAULT 'amd64' COMMENT 'CPU架构';

-- 更新现有空值的记录
UPDATE `system_images` SET `provider_type` = 'docker' WHERE `provider_type` = '' OR `provider_type` IS NULL;
UPDATE `system_images` SET `architecture` = 'amd64' WHERE `architecture` = '' OR `architecture` IS NULL;

-- ============================================
-- 修复 4: invite_codes 表确保 used_count 字段存在
-- ============================================
ALTER TABLE `invite_codes` ADD COLUMN IF NOT EXISTS `used_count` int(11) NOT NULL DEFAULT 0 COMMENT '已使用次数';

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 修复完成
-- ============================================
SELECT '数据库修复完成' AS result;