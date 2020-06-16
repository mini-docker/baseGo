-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_admin`;
CREATE TABLE `red_system_admin` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(64) DEFAULT '',
  `password` varchar(64) NOT NULL,
  `role_id` varchar(64) NOT NULL,
  `is_online` TINYINT(1) DEFAULT 2,
  `last_ip` varchar(64) NOT NULL,
  `last_login_time` int(11) NOT NULL,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_menu`;
CREATE TABLE `red_system_menu` (
  `Id` int(11) unsigned NOT NULL,
  `parent_id` int(11) DEFAULT NULL,
  `name` varchar(64) NOT NULL,
  `route` varchar(64) NOT NULL,
  `icon` varchar(64) NOT NULL,
  `level` varchar(64) NOT NULL,
  `status` TINYINT(1) DEFAULT 2,
  `is_show` TINYINT(1) DEFAULT 1,
  `sort` varchar(64) NOT NULL,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  `update_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;


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
  `md5key` varchar(64) NOT NULL,
  `rsa_pub_key` varchar(64) NOT NULL,
  `rsa_pri_key` varchar(64) NOT NULL,
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




