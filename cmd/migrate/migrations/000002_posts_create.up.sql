CREATE table if not EXISTS posts(
   id bigserial PRIMARY KEY,
   title text not null,
   user_id bigint not null,
   content text not null,
   created_at timestamp(0) with time zone not null default now()
);