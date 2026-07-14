/*
 Navicat Premium Dump SQL

 Source Server         : 本地
 Source Server Type    : MySQL
 Source Server Version : 80029 (8.0.29)
 Source Host           : localhost:3306
 Source Schema         : chat

 Target Server Type    : MySQL
 Target Server Version : 80029 (8.0.29)
 File Encoding         : 65001

 Date: 14/07/2026 09:42:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `conf_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `conf_key` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `conf_value` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `user_id` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `conf_key`(`conf_key` ASC) USING BTREE,
  INDEX `user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of config
-- ----------------------------
INSERT INTO `config` VALUES (1, 'Announcement', 'AllNotice', 'Open source customer support system at your service', 'agent');
INSERT INTO `config` VALUES (2, 'Offline Message', 'OfflineMessage', 'I am currently offline and will reply to you later!', 'agent');
INSERT INTO `config` VALUES (3, 'Welcome Message', 'WelcomeMessage', 'How may I help you?', 'agent');
INSERT INTO `config` VALUES (4, 'Email Address (SMTP)', 'NoticeEmailSmtp', '', 'agent');
INSERT INTO `config` VALUES (5, 'Email Account', 'NoticeEmailAddress', '', 'agent');
INSERT INTO `config` VALUES (6, 'Email Password (SMTP)', 'NoticeEmailPassword', '', 'agent');

-- ----------------------------
-- Table structure for ipblack
-- ----------------------------
DROP TABLE IF EXISTS `ipblack`;
CREATE TABLE `ipblack`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `ip` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `kefu_id` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `ip`(`ip` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ipblack
-- ----------------------------

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `kefu_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `visitor_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `content` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `mes_type` enum('kefu','visitor') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'visitor',
  `status` enum('read','unread') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'unread',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `kefu_id`(`kefu_id` ASC) USING BTREE,
  INDEX `visitor_id`(`visitor_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of message
-- ----------------------------
INSERT INTO `message` VALUES (1, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', '顶顶顶顶', '2026-07-13 13:26:00', '2026-07-13 13:26:00', NULL, 'visitor', 'unread');
INSERT INTO `message` VALUES (2, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 'I am currently offline and will reply to you later!', '2026-07-13 13:26:01', '2026-07-13 13:26:01', NULL, 'kefu', 'unread');
INSERT INTO `message` VALUES (3, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', '顶顶顶顶', '2026-07-13 13:26:09', '2026-07-13 13:26:09', NULL, 'visitor', 'unread');
INSERT INTO `message` VALUES (4, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 'I am currently offline and will reply to you later!', '2026-07-13 13:26:10', '2026-07-13 13:26:10', NULL, 'kefu', 'unread');
INSERT INTO `message` VALUES (5, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 'ssssssssss', '2026-07-13 13:32:10', '2026-07-13 13:32:10', NULL, 'kefu', 'unread');
INSERT INTO `message` VALUES (6, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 'vvvvvvvvvvvv', '2026-07-13 13:32:36', '2026-07-13 13:32:36', NULL, 'visitor', 'unread');
INSERT INTO `message` VALUES (7, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 'img[/static/upload/2026July/b3253839336313138222a0c3a88decc1.jpg]', '2026-07-13 13:32:59', '2026-07-13 13:32:59', NULL, 'visitor', 'unread');
INSERT INTO `message` VALUES (8, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 'eeee', '2026-07-13 14:54:19', '2026-07-13 14:54:19', NULL, 'kefu', 'unread');
INSERT INTO `message` VALUES (9, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 'eeeeee', '2026-07-13 14:56:02', '2026-07-13 14:56:02', NULL, 'kefu', 'unread');
INSERT INTO `message` VALUES (10, 'agent', '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 'img[/static/upload/2026July/9a45c68b6e6d0b8922acea7ed0ede4fd.jpg]', '2026-07-13 14:56:09', '2026-07-13 14:56:09', NULL, 'kefu', 'unread');

-- ----------------------------
-- Table structure for reply_group
-- ----------------------------
DROP TABLE IF EXISTS `reply_group`;
CREATE TABLE `reply_group`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `user_id` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of reply_group
-- ----------------------------
INSERT INTO `reply_group` VALUES (1, 'Frequently Asked Questions', 'agent');

-- ----------------------------
-- Table structure for reply_item
-- ----------------------------
DROP TABLE IF EXISTS `reply_item`;
CREATE TABLE `reply_item`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `content` varchar(1024) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `group_id` int NOT NULL DEFAULT 0,
  `user_id` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `item_name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_id`(`user_id` ASC) USING BTREE,
  INDEX `group_id`(`group_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of reply_item
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `password` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `nickname` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `avator` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'agent', 'b33aed8f3134996703dc39f9a7c95783', 'Open Source LiveChat Support', '2020-06-27 19:32:41', '2020-07-04 09:32:20', NULL, '/static/images/4.jpg');

-- ----------------------------
-- Table structure for visitor
-- ----------------------------
DROP TABLE IF EXISTS `visitor`;
CREATE TABLE `visitor`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `avator` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `source_ip` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `to_id` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `visitor_id` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 0,
  `refer` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `last_message` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `city` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `client_ip` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `extra` varchar(2048) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `visitor_id`(`visitor_id` ASC) USING BTREE,
  INDEX `to_id`(`to_id` ASC) USING BTREE,
  INDEX `idx_update`(`updated_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of visitor
-- ----------------------------
INSERT INTO `visitor` VALUES (1, 'Guest', '/static/images/2.png', '127.0.0.1', 'agent', '2026-07-13 13:24:50', '2026-07-13 16:10:24', NULL, '4a93586b-f1b6-40bf-ab94-744c8a1b481b', 0, 'Direct access', 'img[/static/upload/2026July/9a45c68b6e6d0b8922acea7ed0ede4fd.jpg]', '', '127.0.0.1', '');

SET FOREIGN_KEY_CHECKS = 1;
