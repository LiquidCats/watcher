create table if not exists heights(
    blockchain varchar(15) primary key,
    height_value bigint default 0
);

create table if not exists blocks_bitcoin(
    height bigint primary key,
    hash varchar not null,
    previous varchar not null,
    is_confirmed bool default false
);

create table if not exists blocks_ethereum(
    height bigint primary key,
    hash varchar not null,
    previous varchar not null,
    is_confirmed bool default false
);

insert into heights (blockchain, height_value) VALUES ('ethereum', 4961961);
insert into heights (blockchain, height_value) VALUES ('bitcoin',  822910);