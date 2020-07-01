
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
  `level` TINYINT(1) DEFAULT 0,
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


-- 线路提成 根据 统计数据 
-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
-- DROP TABLE IF EXISTS `red_order_record`;
-- CREATE TABLE `red_order_record` (
--   `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
--   `line_id` varchar(64) NOT NULL,
--   `nn_royalty` FLOAT(10,2) NOT NULL,
--   `sl_royalty` FLOAT(10,2) NOT NULL,
--   PRIMARY KEY (`Id`)
-- ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_user`;
CREATE TABLE `red_user` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `account` varchar(64) NOT NULL,
  `password` varchar(64) NOT NULL,
  `is_online` TINYINT(1) DEFAULT 2,
  `balance` FLOAT(10,2) NOT NULL,
  `ip` varchar(64) NOT NULL,
  `status` TINYINT(1) DEFAULT 1,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  `edit_time` int(11) NOT NULL,
  `capital` FLOAT(10,2) NOT NULL,
  `last_login_ip` varchar(255) NOT NULL,
  `last_login_time` int(11) NOT NULL,
  `is_robot` TINYINT(1) DEFAULT 2,
  `is_group_owner` TINYINT(1) DEFAULT 2,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



-- ----------------------------
-- Table structure for cate
-- ----------------------------
DROP TABLE IF EXISTS `red_system_game`;
CREATE TABLE `red_system_game` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `game_name` VARCHAR(64) NOT NULL,
  `game_type` int(11) NOT NULL,
  `status` TINYINT(1) DEFAULT 1,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- // 线路 站点Id 群名称 三级联动 start --

-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_room`;
CREATE TABLE `red_room` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `room_name` varchar(64) NOT NULL,
  `game_type` int(11) NOT NULL,
  `max_money` FLOAT(10,2) NOT NULL,
  `min_money` FLOAT(10,2) NOT NULL,
  `game_play` int(11) NOT NULL,
  `odds` FLOAT(10,2) NOT NULL,
  `red_num` int(11) NOT NULL,
  `red_min_num` int(11) NOT NULL,
  `royalty` FLOAT(10,2) NOT NULL,
  `royalty_money` FLOAT(10,2) NOT NULL,
  `game_time` int(11) NOT NULL,
  `room_sort` int(11) NOT NULL,
  `status` TINYINT(1) DEFAULT 1,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  `room_type` int(11) NOT NULL,
  `free_from_death` TINYINT(1) DEFAULT 2,
  `robot_send_packet` TINYINT(1) DEFAULT 2,
  `robot_send_packet_time` int(11) NOT NULL,
  `robot_grab_packet` TINYINT(1) DEFAULT 2,
  `room_no` int(11) NOT NULL,
  `robot_id` int(11) NOT NULL,
  `control_kill` TINYINT(1) DEFAULT 1,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_order_record`;
CREATE TABLE `red_order_record` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `user_id` int(11) NOT NULL,
  `account` varchar(64) NOT NULL,
  `red_sender` varchar(64) NOT NULL,
  `game_type` int(11) NOT NULL,
  `game_play` int(11) NOT NULL,
  `room_id` int(11) NOT NULL,
  `room_name` varchar(64) NOT NULL,
  `order_no` varchar(64) NOT NULL,
  `red_id` int(11) NOT NULL,
  `red_money` FLOAT(10,2) NOT NULL,
  `red_num` int(11) NOT NULL,
  `receive_money` FLOAT(10,2) NOT NULL,
  `royalty` FLOAT(10,2) NOT NULL,
  `royalty_money` FLOAT(10,2) NOT NULL,
  `money` FLOAT(10,2) NOT NULL,
  `real_money` FLOAT(10,2) NOT NULL,
  `game_time` int(11) NOT NULL,
  `receive_time` int(11) NOT NULL,
  `red_start_time` int(11) NOT NULL,
  `status` TINYINT(1) NOT NULL,
  `extra` varchar(255) NOT NULL,
  `is_robot` TINYINT(1) DEFAULT 2,
  `is_free_death` TINYINT(1) DEFAULT 2,
  `robot_win` FLOAT(10,2) NOT NULL,
  `valid_bet` FLOAT(10,2) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_packet`;
CREATE TABLE `red_packet` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `red_envelope_amount` FLOAT(10,2) NOT NULL,
  `red_envelopes_num` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `account` varchar(64) NOT NULL,
  `create_time` int(11) NOT NULL,
  `delete_time` int(11) NOT NULL,
  `red_type` int(11) NOT NULL,
  `red_play` int(11) NOT NULL,
  `room_id` int(11) NOT NULL,
  `room_name` varchar(64) NOT NULL,
  `status` TINYINT(1) NOT NULL,
  `mine` int(11) NOT NULL,
  `capital` FLOAT(10,2) NOT NULL,
  `money` FLOAT(10,2) NOT NULL,
  `real_money` FLOAT(10,2) NOT NULL,
  `royalty_money` FLOAT(10,2) NOT NULL,
  `return_money` FLOAT(10,2) NOT NULL,
  `is_auto` TINYINT(1) NOT NULL,
  `auto_time` int(11) NOT NULL,
  `end_time` int(11) NOT NULL,
  `is_robot` TINYINT(1) DEFAULT 2,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- // 线路 站点Id 群名称 三级联动 end --


-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_post`;
CREATE TABLE `red_post` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `title` varchar(64) NOT NULL,
  `start_time` int(11) NOT NULL,
  `end_time` int(11) NOT NULL,
  `sort` int(11) NOT NULL,
  `status` TINYINT(1) DEFAULT 1,
  `delete_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_active_picture`;
CREATE TABLE `red_active_picture` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `active_name` varchar(64) NOT NULL,
  `picture` varchar(64) NOT NULL,
  `start_time` int(11) NOT NULL,
  `end_time` int(11) NOT NULL,
  `sort` int(11) NOT NULL,
  `status` TINYINT(1) DEFAULT 1,
  `delete_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_order_statistical`;
CREATE TABLE `red_order_statistical` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `statistical_date` varchar(255) NOT NULL,
  `valid_bet` FLOAT(10,2) NOT NULL,
  `red_num` int(11) NOT NULL,
  `order_num` int(11) NOT NULL,
  `royalty_money` FLOAT(10,2) NOT NULL,
  `free_death_win` FLOAT(10,2) NOT NULL,
  `robot_win` FLOAT(10,2) NOT NULL,
  `total_win` FLOAT(10,2) NOT NULL,
  `game_type` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_log`;
CREATE TABLE `red_log` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `log_type` TINYINT(1) NOT NULL,
  `remark` varchar(255)NOT NULL,
  `creator` varchar(64) NOT NULL,
  `creator_id` int(11) NOT NULL,
  `creator_ip` varchar(255)NOT NULL,
  `create_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;




-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_message_history`;
CREATE TABLE `red_message_history` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `msg_type` TINYINT(1) NOT NULL,
  `msg_content` varchar(1024)NOT NULL,
  `sender_id` int(11) NOT NULL,
  `sender_name` varchar(64) NOT NULL,
  `status` TINYINT(1) NOT NULL,
  `send_time` int(11) NOT NULL,
  `room_id` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_member_cash_record`;
CREATE TABLE `red_member_cash_record` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64) NOT NULL,
  `agency_id` varchar(64) NOT NULL,
  `order_no` varchar(64) NOT NULL,
  `game_type` int(11) NOT NULL,
  `game_name` varchar(64)NOT NULL,
  `flow_type` int(11) NOT NULL,
  `money` FLOAT(10,2) NOT NULL,
  `remark` varchar(255) NOT NULL,
  `create_time` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `account` varchar(64)NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_post_content`;
CREATE TABLE `red_post_content` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pid` int(11) NOT NULL,
  `content` varchar(1024) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;





-- -- ----------------------------
-- -- Table structure for cate
-- -- ----------------------------
DROP TABLE IF EXISTS `red_packet_collect`;
CREATE TABLE `red_packet_collect` (
  `Id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `line_id` varchar(64)NOT NULL,
  `agency_id` varchar(64)NOT NULL,
  `settlement_info` varchar(1024)NOT NULL,
  `collect_status` TINYINT(1) DEFAULT 1,
  `create_time` int(11) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;






