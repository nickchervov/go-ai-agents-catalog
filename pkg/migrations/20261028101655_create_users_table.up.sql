CREATE TABLE IF NOT EXISTS users
(
    id uuid PRIMARY KEY,
    email varchar(128) unique not null default '',
    role varchar(5) not null check(role IN ('admin','user')),
    created_at timestamptz not null default current_timestamp
);

create index idx_user_role on users(role);