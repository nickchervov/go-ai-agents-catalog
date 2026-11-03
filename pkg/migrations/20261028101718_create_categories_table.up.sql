CREATE TABLE IF NOT EXISTS categories
(
    id uuid primary key,
    category varchar(128) unique not null default '',
    description text,
    created_at timestamptz not null default current_timestamp
);