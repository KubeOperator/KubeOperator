create table if not exists credential
(
	id varchar(64) not null,
	name varchar(128) not null,
	user varchar(128) default 'root' not null,
	password varchar(256) default '' null,
	private_key text null,
	constraint credential_pk
		primary key (id)
);

create unique index credential_name_uindex
	on credential (name);

