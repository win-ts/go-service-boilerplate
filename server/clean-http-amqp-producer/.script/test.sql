CREATE DATABASE IF NOT EXISTS `test`;
USE `test`;

DROP TABLE IF EXISTS `tbl_test`;

CREATE TABLE `tbl_test` (
	`id` INT auto_increment NOT NULL,
	`message` varchar(100) NULL,
	CONSTRAINT `test_pk` PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `tbl_test` (`message`) VALUES
    ('test1'),
    ('test2'),
    ('test3');
