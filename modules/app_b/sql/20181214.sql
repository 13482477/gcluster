CREATE TABLE `ecs` (
  `id` bigint AUTO_INCREMENT,
  `created_time` timestamp NULL,
  `updated_time` timestamp NULL,
  `ext` varchar(255),`platform` int,
  `instance_id` varchar(255),
  `instance_name` varchar(255),
  `public_ip` varchar(255),
  `private_ip` varchar(255),
  `os_type` varchar(255),
  `cpu` int,
  `mem` int,
  `disk_size` varchar(255),
  `instance_type` varchar(255),
  `zone_id` varchar(255),
  `region_id` varchar(255),
  `grafana_url` varchar(255),
  `account_id` bigint,
  `account_name` varchar(255),
  `status` int DEFAULT 0,
  `group_id` bigint ,
  PRIMARY KEY (`id`)
);

CREATE UNIQUE INDEX ecs_index ON `ecs`(`platform`, instance_id);