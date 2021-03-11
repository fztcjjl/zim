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

-- 导出  表 zim.im_msg_recv_00 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv_00` (
  `id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '额外内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `to` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivered` tinyint(4) NOT NULL DEFAULT '0' COMMENT '送达状态[0:未送达;1:已送达]',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `seq` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息序号',
  `client_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端发送时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_recv_01 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv_01` (
  `id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '额外内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `to` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivered` tinyint(4) NOT NULL DEFAULT '0' COMMENT '送达状态[0:未送达;1:已送达]',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `seq` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息序号',
  `client_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端发送时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_recv_02 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv_02` (
  `id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '额外内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `to` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivered` tinyint(4) NOT NULL DEFAULT '0' COMMENT '送达状态[0:未送达;1:已送达]',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `seq` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息序号',
  `client_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端发送时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_recv_03 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv_03` (
  `id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '额外内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `to` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivered` tinyint(4) NOT NULL DEFAULT '0' COMMENT '送达状态[0:未送达;1:已送达]',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `seq` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息序号',
  `client_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端发送时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_recv_04 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv_04` (
  `id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '额外内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `to` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivered` tinyint(4) NOT NULL DEFAULT '0' COMMENT '送达状态[0:未送达;1:已送达]',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `seq` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息序号',
  `client_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端发送时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_recv_05 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv_05` (
  `id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '额外内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `to` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivered` tinyint(4) NOT NULL DEFAULT '0' COMMENT '送达状态[0:未送达;1:已送达]',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `seq` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息序号',
  `client_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端发送时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_recv_06 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv_06` (
  `id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '额外内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `to` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivered` tinyint(4) NOT NULL DEFAULT '0' COMMENT '送达状态[0:未送达;1:已送达]',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `seq` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息序号',
  `client_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端发送时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.im_msg_recv_07 结构
CREATE TABLE IF NOT EXISTS `im_msg_recv_07` (
  `id` bigint(20) NOT NULL COMMENT '消息ID',
  `conv_type` tinyint(4) NOT NULL COMMENT '会话类型[1:单聊;2:群聊]',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '额外内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者',
  `to` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivered` tinyint(4) NOT NULL DEFAULT '0' COMMENT '送达状态[0:未送达;1:已送达]',
  `target` varchar(50) NOT NULL DEFAULT '' COMMENT '目标',
  `seq` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息序号',
  `client_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '客户端发送时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='消息收件箱';

-- 数据导出被取消选择。


-- 导出  表 zim.seq 结构
CREATE TABLE IF NOT EXISTS `seq` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `obj_type` tinyint(4) NOT NULL DEFAULT '0',
  `obj_id` varchar(50) NOT NULL DEFAULT '0',
  `seq` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='序号生成器';

-- 数据导出被取消选择。
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
