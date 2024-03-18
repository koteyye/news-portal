-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id uuid not null default gen_random_uuid() primary key unique,
    login varchar(512) not null unique,
    hashed_password varchar(512) not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp
);
comment on table users is 'пользователи';

create table if not exists avatar (
    id uuid not null default gen_random_uuid() primary key unique,
    file_name varchar(1024) not null,
    bucket_name varchar(1024) not null,
    mime_type varchar(1024) not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp
);
comment on table avatar is 'аватар пользователя';

create table if not exists profile (
    id uuid not null default gen_random_uuid() primary key unique,
    user_id uuid not null,
    username varchar(512) not null,
    first_name varchar(512),
    last_name varchar(512),
    sur_name varchar(512),
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp,
    avatar_id uuid,
    foreign key (user_id) references users(id),
    foreign key (avatar_id) references avatar(id)
);
comment on table profile is 'профиль пользователя';

create table if not exists roles (
    id uuid not null default gen_random_uuid() primary key unique,
    role_name varchar(256)
);
comment on table roles is 'роли';

create table if not exists user_roles (
    user_id uuid not null,
    role_id uuid not null,
    foreign key (user_id) references users(id),
    foreign key (role_id) references roles(id)
);
comment on table user_roles is 'связь m2m пользователя и роли';

insert into roles (role_name) values ('reader'), ('writer');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table user_roles;

drop table roles;

drop table profile;

drop table avatar;

drop table users;
-- +goose StatementEnd
