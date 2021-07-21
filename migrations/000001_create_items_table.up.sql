CREATE TABLE IF NOT EXISTS items(
    id text NOT NULL,
    owner_id text NOT NULL,
    name text NOT NULL,
    status text NOT NULL,
    description text NOT NULL,
    total_quantity INT NOT NULL,
    locked_quantity INT NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    primary key(id)
)

