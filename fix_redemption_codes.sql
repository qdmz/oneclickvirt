-- 修复兑换码表结构
USE oneclickvirt;

-- 删除旧的兑换码表
DROP TABLE IF EXISTS redemption_codes;

-- 创建新的兑换码表
CREATE TABLE `redemption_codes` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(32) NOT NULL COMMENT '兑换码',
  `type` varchar(20) NOT NULL COMMENT '兑换码类型',
  `value` bigint(20) NOT NULL DEFAULT '0' COMMENT '金额(分)或等级数',
  `product_id` int(11) unsigned DEFAULT NULL COMMENT '产品ID',
  `max_uses` int(11) NOT NULL DEFAULT '1' COMMENT '最大使用次数',
  `used_count` int(11) NOT NULL DEFAULT '0' COMMENT '已使用次数',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `expire_at` datetime(3) DEFAULT NULL COMMENT '过期时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_type` (`type`),
  KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='兑换码表';

-- 删除旧的兑换码使用记录表
DROP TABLE IF EXISTS redemption_code_usages;

-- 创建新的兑换码使用记录表
CREATE TABLE `redemption_code_usages` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `code_id` int(11) unsigned NOT NULL COMMENT '兑换码ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `reward` json DEFAULT NULL COMMENT '奖励详情',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_code_id` (`code_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='兑换码使用记录表';
