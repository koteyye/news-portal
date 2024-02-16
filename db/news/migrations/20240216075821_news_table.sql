-- +goose Up
-- +goose StatementBegin
create type news_state as enum ('DRAFT', 'PUBLISHED')

create table if not exists files (
    id uuid not null default gen_random_uuid() primary key unique,
    file_name varchar(1024) not null,
    bucket_name varchar(1024) not null,
    mime_type varchar(512) not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp
)

create table if not exists news (
    id uuid not null default gen_random_uuid() primary key unique,
    title varchar(1024) not null,
    author uuid,
    description text,
    content_id uuid,
    preview_id uuid,
    state news_state not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp,
    user_created uuid not null,
    user_updated uuid,
    user_deleted uuid,
    foreign key (content_id) references files(id),
    foreign key (preview_id) references files(id)
)

create table if not exists likes (
    id uuid not null default gen_random_uuid() primary key unique,
    news_id uuid not null,
    liker uuid not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    is_active boolean,
    foreign key (news_id) references news(id)
)

create table if not exists comments (
    id uuid not null default gen_random_uuid() primary key unique,
    news_id uuid not null,
    author uuid not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    is_active boolean,
    foreign key (news_id) references news(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table comments

drop table likes

drop table news

drop table files

drop type news_state
-- +goose StatementEnd
