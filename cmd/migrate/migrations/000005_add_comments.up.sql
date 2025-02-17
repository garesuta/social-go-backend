create table if not EXISTS comments (
    id bigserial primary key,
    post_id bigserial not NULL,
    user_id bigserial not NULL,
    content text not null,
    created_at timestamp(0) with time zone not null default now()
);