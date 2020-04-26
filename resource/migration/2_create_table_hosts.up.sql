create table if not exists hosts
(
	id varchar(64) not null,
	name varchar(128) not null,
	ip varchar(128) not null,
	port int default 22 not null,
	constraint hosts_pk
		primary key (id)
);

create unique index hosts_ip_uindex
	on hosts (ip);

create unique index hosts_name_uindex
	on hosts (name);

