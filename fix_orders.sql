-- 修复订单表结构
USE oneclickvirt;

-- 删除旧的订单表
DROP TABLE IF EXISTS orders;

-- 创建新的订单表
CREATE TABLE `orders` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL COMMENT 'UUID',
  `user_id` int(11) unsigned NOT NULL COMMENT '用户ID',
  `product_id` int(11) unsigned NOT NULL COMMENT '产品ID',
  `instance_id` int(11) unsigned DEFAULT NULL COMMENT '实例ID',
  `order_no` varchar(64) NOT NULL COMMENT '订单号',
  `amount` decimal(10,2) NOT NULL COMMENT '订单金额',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '订单状态',
  `payment_method` varchar(20) DEFAULT NULL COMMENT '支付方式',
  `payment_id` varchar(64) DEFAULT NULL COMMENT '支付ID',
  `paid_at` datetime(3) DEFAULT NULL COMMENT '支付时间',
  `provisioned_at` datetime(3) DEFAULT NULL COMMENT '开通时间',
  `expired_at` datetime(3) DEFAULT NULL COMMENT '过期时间',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_data` text DEFAULT NULL COMMENT '产品数据（JSON）',
  `expire_at` datetime(3) DEFAULT NULL COMMENT '过期时间（兼容）',
  `payment_time` datetime(3) DEFAULT NULL COMMENT '支付时间（兼容）',
  `paid_amount` decimal(10,2) DEFAULT '0.00' COMMENT '已支付金额',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_instance_id` (`instance_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 删除旧的支付记录表
DROP TABLE IF EXISTS payment_records;

-- 创建新的支付记录表
CREATE TABLE `payment_records` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL COMMENT 'UUID',
  `order_id` int(11) unsigned NOT NULL COMMENT '订单ID',
  `user_id` int(11) unsigned NOT NULL COMMENT '用户ID',
  `payment_no` varchar(64) NOT NULL COMMENT '支付号',
  `amount` decimal(10,2) NOT NULL COMMENT '支付金额',
  `payment_method` varchar(20) DEFAULT NULL COMMENT '支付方式',
  `payment_status` varchar(20) DEFAULT NULL COMMENT '支付状态',
  `payment_time` datetime(3) DEFAULT NULL COMMENT '支付时间',
  `notify_data` text DEFAULT NULL COMMENT '回调数据',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `type` varchar(20) DEFAULT NULL COMMENT '支付类型',
  `transaction_id` varchar(64) DEFAULT NULL COMMENT '交易ID',
  `status` varchar(20) DEFAULT NULL COMMENT '状态（兼容）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`),
  UNIQUE KEY `idx_payment_no` (`payment_no`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_payment_status` (`payment_status`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';
