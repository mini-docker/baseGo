
-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_agency`;
CREATE TABLE `red_agency` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `account` varchar(64) NOT NULL,
  `password` varchar(64) NOT NULL,
  `is_online` TINYINT(1) DEFAULT 2,
  `is_admin` TINYINT(1) DEFAULT 1,
  `status` TINYINT(1) DEFAULT 2,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  `edit_time` int(11) NOT NULL,
  `white_ip_address` varchar(255) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_admin`;
CREATE TABLE `red_system_admin` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` varchar(64) NOT NULL,
  `account` varchar(64) NOT NULL,
  `password` varchar(64) NOT NULL,
  `is_online` TINYINT(1) DEFAULT 2,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  `last_login_time` int(11) NOT NULL,
  `last_ip` varchar(255) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- /system/agencyCode
-- ----------------------------
-- 将添加的超管 手动 添加到数据库
-- red_agency ==> red_system_admin
-- 
-- ----------------------------





-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_packet_site`;
CREATE TABLE `red_packet_site` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `site_name` varchar(64) NOT NULL,
  `status` TINYINT(1) DEFAULT 1,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_menu`;
CREATE TABLE `red_system_menu` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `name` varchar(64) NOT NULL,
  `route` varchar(64) NOT NULL,
  `icon` varchar(64) NOT NULL,
  `level` TINYINT(1) DEFAULT 1,
  `status` TINYINT(1) DEFAULT 1,
  `is_show` TINYINT(1) DEFAULT 1,
  `sort` int(11) NOT NULL,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) DEFAULT 0,
  `update_time` int(11) DEFAULT 0,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_role_menu`;
CREATE TABLE `red_system_role_menu` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `menu_id` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;



-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_role`;
CREATE TABLE `red_system_role` (
  `Id` int(11) unsigned NOT NULL,
  `role_name` int(11) NOT NULL,
  `is_default` varchar(64) DEFAULT 2,
  `status` TINYINT(1) DEFAULT 1,
  `remark` varchar(64) NOT NULL,
  `edit_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  `update_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;






-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_line`;
CREATE TABLE `red_system_line` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `line_name` varchar(64) NOT NULL,
  `limit_cost` int(11) NOT NULL,
  `meal_id` int(11) NOT NULL,
  `domain` varchar(64) NOT NULL,
  `status` TINYINT(1) DEFAULT 1,
  `trans_type` TINYINT(1) DEFAULT 1,
  `api_url` varchar(64) NOT NULL,
  `md5key` varchar(255) NOT NULL,
  `rsa_pub_key` varchar(1024) NOT NULL,
  `rsa_pri_key` varchar(1024) NOT NULL,
  `create_time` int(11) NOT NULL,
  `edit_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_line_meal`;
CREATE TABLE `red_system_line_meal` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `meal_name` varchar(64) NOT NULL,
  `nn_royalty` int(11) NOT NULL,
  `sl_royalty` int(11) NOT NULL,
  `create_time` int(11) NOT NULL,
  `edit_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;








-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_role`;
CREATE TABLE `red_system_role` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_name` varchar(64) NOT NULL,
  `is_default` TINYINT(1) DEFAULT 2,
  `status` TINYINT(1) DEFAULT 1,
  `remark` varchar(255) NOT NULL,
  `edit_time` int(11) NOT NULL,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_role_menu`;
CREATE TABLE `red_system_role_menu` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `menu_id` int(11) DEFAULT 2,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;










