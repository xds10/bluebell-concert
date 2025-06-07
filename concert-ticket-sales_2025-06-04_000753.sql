-- MySQL dump 10.13  Distrib 8.0.39, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: concert-ticket-sales
-- ------------------------------------------------------
-- Server version	8.0.40

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `concert`
--

DROP TABLE IF EXISTS `concert`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `concert` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '演唱会ID',
  `title` varchar(100) DEFAULT NULL COMMENT '演唱会名称',
  `venue_id` int DEFAULT NULL COMMENT '场馆名称',
  `start_time` datetime DEFAULT NULL COMMENT '演出开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '演出结束时间',
  `release_time` datetime DEFAULT NULL COMMENT '放票时间',
  `status` tinyint(1) DEFAULT '0' COMMENT '状态(0:未开始,1:进行中,2:已结束)',
  PRIMARY KEY (`id`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='演唱会排班表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `concert_id` bigint unsigned NOT NULL COMMENT '演唱会ID',
  `seat_id` int NOT NULL COMMENT '座位ID',
  `price` decimal(10,2) NOT NULL COMMENT '订单金额',
  `status` int NOT NULL DEFAULT '1' COMMENT '订单状态(1:待支付,2:已支付,3:已取消,4:已过期)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_concert_id` (`concert_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_orders_concert_id` FOREIGN KEY (`concert_id`) REFERENCES `concert` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `seat`
--

DROP TABLE IF EXISTS `seat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `seat` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '座位ID',
  `place_id` bigint unsigned NOT NULL COMMENT '场馆ID',
  `place` varchar(20) NOT NULL COMMENT '场馆名字',
  `section` varchar(20) NOT NULL COMMENT '区域',
  `seat_row` tinyint NOT NULL COMMENT '排号',
  `seat_no` smallint NOT NULL COMMENT '座位号',
  `price` decimal(10,2) NOT NULL COMMENT '票价',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_concert_section_row_seat` (`section`,`seat_row`,`seat_no`),
  KEY `idx_place_id` (`place_id`),
  KEY `idx_section` (`section`),
  KEY `idx_price` (`price`)
) ENGINE=InnoDB AUTO_INCREMENT=61 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='场馆座位表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `seat_record`
--

DROP TABLE IF EXISTS `seat_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `seat_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `concert_id` bigint unsigned NOT NULL COMMENT '演唱会ID',
  `seat_id` int NOT NULL COMMENT '座位信息',
  `is_selected` int DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_concert_id` (`concert_id`)
) ENGINE=InnoDB AUTO_INCREMENT=301 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='座位选择表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(128) NOT NULL COMMENT '密码',
  `is_blacklisted` tinyint(1) DEFAULT '0' COMMENT '黑名单标记(0:正常,1:封禁)',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'concert-ticket-sales'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-06-04  0:07:55
