-- 创建数据库
CREATE DATABASE IF NOT EXISTS oauth2 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE oauth2;

-- 客户端信息表
CREATE TABLE IF NOT EXISTS `client` (
    `id` VARCHAR(64) NOT NULL COMMENT '客户端ID',
    `secret` VARCHAR(128) NOT NULL COMMENT '客户端密钥',
    `name` VARCHAR(100) NOT NULL COMMENT '应用名称',
    `redirect_url` VARCHAR(500) NOT NULL COMMENT '回调地址',
    `grant_type` VARCHAR(50) NOT NULL COMMENT '支持的授权模式',
    `scope` VARCHAR(200) NOT NULL COMMENT '请求的权限范围',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户端信息表';

-- 权限申请记录表
CREATE TABLE IF NOT EXISTS `authorization` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `client_id` VARCHAR(64) NOT NULL COMMENT '客户端ID',
    `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
    `scope` VARCHAR(200) NOT NULL COMMENT '请求的权限范围',
    `code` VARCHAR(128) NOT NULL COMMENT '授权码',
    `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending/approved/rejected',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_client_user` (`client_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限申请记录表';

-- 插入一些测试数据
INSERT INTO `client` (`id`, `secret`, `name`, `redirect_url`, `grant_type`, `scope`) VALUES
('trusted_client_001', 'trusted_secret_001', '可信应用1', 'http://localhost:3000/callback', 'authorization_code', 'userid profile'),
('trusted_client_002', 'trusted_secret_002', '可信应用2', 'http://localhost:3001/callback', 'authorization_code', 'userid'),
('test_client_001', 'test_secret_001', '测试应用1', 'http://localhost:3002/callback', 'authorization_code', 'userid profile'); 