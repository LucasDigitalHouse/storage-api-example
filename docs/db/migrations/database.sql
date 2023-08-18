-- DDL: Data Definition Language
DROP DATABASE IF EXISTS `storage_api_db`;

CREATE DATABASE `storage_api_db`;

USE `storage_api_db`;

CREATE TABLE `warehouses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `count` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `warehouse_id` int NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_warehouse_id` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;