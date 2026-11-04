CREATE TABLE IF NOT EXISTS assistant_run
(   
    id uuid primary key,
    assistant_id uuid references assistants(id) not null,
    user_id uuid references users(id) not null,
    user_prompt text not null default '',
    result text,
    status varchar(7) not null default 'pending' check(status IN ('pending','success','failed')),
    error text,
    created_at timestamptz not null default current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_assistant_run_status ON assistant_run(status);
CREATE INDEX IF NOT EXISTS idx_assistant_run_user_id ON assistant_run(user_id);
CREATE INDEX IF NOT EXISTS idx_assistant_run_created_at ON assistant_run(created_at);
CREATE INDEX IF NOT EXISTS idx_assistant_run_assistant_id ON assistant_run(assistant_id);