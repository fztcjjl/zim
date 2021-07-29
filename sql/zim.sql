-- --------------------------------------------------------
-- 主机:                           172.16.13.4
-- 服务器版本:                        5.7.25 - MySQL Community Server (GPL)
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  9.2.0.4947
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- 导出  表 zim.im_group 结构
CREATE TABLE IF NOT EXISTS `im_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '系统编号',
  `owner` varchar(50) NOT NULL DEFAULT '' COMMENT '群主',
  `group_id` varchar(50) NOT NULL DEFAULT '' COMMENT '群ID',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '群类型',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '群名称',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `group_id_deleted_at` (`group_id`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- 数据导出被取消选择。


-- 导出  表 zim.im_group_member 结构
CREATE TABLE IF NOT EXISTS `im_group_member` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '系统编号',
  `group_id` varchar(50) NOT NULL DEFAULT '' COMMENT '群ID',
  `member` varchar(50) NOT NULL DEFAULT '' COMMENT '成员ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` bigint(20) NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='群成员';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_recv 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '系统编号',
  `msg_id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '扩展',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` bigint(20) NOT NULL DEFAULT '0' COMMENT '删除时间',
  `sender` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `receiver` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `at_user_list` varchar(5000) NOT NULL DEFAULT '' COMMENT '@用户列表',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_send 结构
CREATE TABLE IF NOT EXISTS `im_msg_send` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '系统编号',
  `msg_id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '扩展',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` bigint(20) NOT NULL DEFAULT '0' COMMENT '删除时间',
  `sender` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `at_user_list` varchar(5000) NOT NULL DEFAULT '' COMMENT '@用户列表',
  `read_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息读取时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息发件箱';

-- 数据导出被取消选择。
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
