create table if not exists users
(
	id varchar(64) not null,
	name varchar(128) not null,
	password varchar(256) not null,
	constraint users_pk
		primary key (id)
);

create unique index users_name_uindex
	on users (name);
