CREATE TABLE IF NOT EXISTS assistants
(
    id uuid primary key,
    name varchar(128) not null default '',
    description text,
    category_id uuid references categories(id) on delete cascade,
    model varchar(256) not null default '',
    system_prompt text not null,
    ex_user_prompt text,
    is_activated boolean not null default true,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz not null default current_timestamp
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_assistants_search_trgm ON assistants USING gin ((name || ' ' || description) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_assistants_category_id ON assistants(category_id);
