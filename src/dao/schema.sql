CREATE DATABASE IF NOT EXISTS `crawler`;

CREATE TABLE IF NOT EXISTS `tasks` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(64) NOT NULL COMMENT '任务名', 
    `seeds` varchar(4096) NOT NULL COMMENT '种子url(多个逗号分隔)', 
    `desci` varchar(256) DEFAULT NULL COMMENT '任务描述', 
    `status` tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '状态: 0未执行, 1失败, 2成功',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unq_name` (`name`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `task_items` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `task_id` int(11) NOT NULL COMMENT '任务ID',
    `name` varchar(64) NOT NULL COMMENT '人名', 
    `desci` varchar(256) DEFAULT NULL COMMENT '图片描述', 
    `url` varchar(512) NOT NULL COMMENT '图片外链', 
    `filepath` varchar(256) DEFAULT NULL COMMENT '图片的本地保存地址', 
    `digest` varchar(32) DEFAULT NULL COMMENT '图片内容的hash签名，可用于去重',
    `status` tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '状态: 0未下载, 1下载失败, 2保存失败,3成功',
    `effective` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '图片有效性: 0未判断, 1有人脸, 2无人脸',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unq_name_url` (`name`, `url`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;
