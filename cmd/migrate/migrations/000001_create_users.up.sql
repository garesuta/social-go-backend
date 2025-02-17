CREATE table if not EXISTS users(
   id bigserial PRIMARY KEY,
   email citext unique not null,
   username varchar(255) unique not null,
   password bytea not null,
   created_at timestamp(0) with time zone not null default now()
);