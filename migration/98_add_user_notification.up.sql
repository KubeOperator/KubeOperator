INSERT INTO `ko_user_notification_config`(`id`, `vars`, `type`, `updated_at`, `created_at`, `user_id`)
select  MD5(UUID()),  '{\"DING_TALK\":\"DISABLE\",\"EMAIL\":\"DISABLE\",\"LOCAL\":\"ENABLE\",\"WORK_WEIXIN\":\"DISABLE\"}', 'CLUSTER', date_add(now(), interval 8 HOUR) ,date_add(now(), interval 8 HOUR),u.id from ko_user as u where u.id not in (select user_id from ko_user_notification_config where type = 'CLUSTER');




INSERT INTO `ko`.`ko_user_notification_config`(`id`, `vars`, `type`, `updated_at`, `created_at`, `user_id`)
select  MD5(UUID()),  '{\"DING_TALK\":\"DISABLE\",\"EMAIL\":\"DISABLE\",\"LOCAL\":\"ENABLE\",\"WORK_WEIXIN\":\"DISABLE\"}', 'SYSTEM', date_add(now(), interval 8 HOUR) ,date_add(now(), interval 8 HOUR),u.id from ko_user as u where u.id not in (select user_id from ko_user_notification_config where type = 'SYSTEM');