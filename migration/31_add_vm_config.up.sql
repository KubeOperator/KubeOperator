CREATE TABLE  IF NOT EXISTS`ko_vm_config` (
  `id` varchar(64) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `cpu` int DEFAULT NULL,
  `memory` int DEFAULT NULL,
  `disk` int DEFAULT NULL,
  `provider` varchar(64) DEFAULT NULL,
   `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);
INSERT INTO `ko`.`ko_vm_config`(`id`, `name`, `cpu`, `memory`, `disk`, `provider`, `created_at`, `updated_at`) VALUES ('15c899a0-b16e-4b6f-8b60-163418776d0c', '2xlarge', 32, 128, 50, 'vSphere', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_vm_config`(`id`, `name`, `cpu`, `memory`, `disk`, `provider`, `created_at`, `updated_at`) VALUES ('88d8c1b1-bff3-45b9-b7d5-fef2bcfaa896', '4xlarge', 64, 256, 50, 'vSphere', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_vm_config`(`id`, `name`, `cpu`, `memory`, `disk`, `provider`, `created_at`, `updated_at`) VALUES ('9b757b84-a5f7-4a79-ab81-258a6dbcdcec', 'medium', 4, 16, 50, 'vSphere', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_vm_config`(`id`, `name`, `cpu`, `memory`, `disk`, `provider`, `created_at`, `updated_at`) VALUES ('ad28892c-3ca1-4fac-8ae9-f749d3493582', 'large', 8, 32, 50, 'vSphere', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_vm_config`(`id`, `name`, `cpu`, `memory`, `disk`, `provider`, `created_at`, `updated_at`) VALUES ('f075feb5-c34b-464e-8661-ab1d625f2083', 'xlarge', 16, 64, 50, 'vSphere', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));
INSERT INTO `ko`.`ko_vm_config`(`id`, `name`, `cpu`, `memory`, `disk`, `provider`, `created_at`, `updated_at`) VALUES ('f17e0c2d-0b67-4f5b-bcea-e7b63afa58d1', 'small', 2, 8, 50, 'vSphere', date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR));