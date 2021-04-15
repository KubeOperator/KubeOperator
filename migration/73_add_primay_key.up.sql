ALTER TABLE `ko`.`ko_cluster_member`
MODIFY COLUMN `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL AFTER `updated_at`,
ADD PRIMARY KEY (`id`);

ALTER TABLE `ko`.`ko_cluster_resource`
MODIFY COLUMN `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL FIRST,
ADD PRIMARY KEY (`id`);