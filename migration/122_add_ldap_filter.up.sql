INSERT INTO `ko_system_setting` (`created_at`, `updated_at`, `id`, `key`, `value`, `tab`)
VALUES
	(date_add(now(), interval 8 HOUR), date_add(now(), interval 8 HOUR),  UUID(), 'ldap_mapping', '{\"Name\":\"cn\",\"Email\":\"mail\"}', 'LDAP');
