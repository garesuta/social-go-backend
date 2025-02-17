create table if not EXISTS user_invitations
(
    token bytea PRIMARY KEY,
    user_id BIGINT not NULL
)