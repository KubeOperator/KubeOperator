INSERT INTO `ko_system_setting` (`created_at`, `updated_at`, `id`, `key`, `value`, `tab`)
VALUES
	( date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR), UUID(), 'size_limit', '1000', 'LDAP');

INSERT INTO `ko_system_setting` (`created_at`, `updated_at`, `id`, `key`, `value`, `tab`)
VALUES
	( date_add(now(), interval 8 HOUR),  date_add(now(), interval 8 HOUR), UUID(), 'time_limit', '30', 'LDAP');
