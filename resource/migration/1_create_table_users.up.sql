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

insert into users values ('4b03397c-acd5-44a1-bd93-4625e388bc9b', 'admin','$2a$10$CLrxkpxsakpnQib2gSbfheg7B04jGWdTYhFW');


