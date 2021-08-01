CREATE TABLE IF NOT EXISTS items(
    id text NOT NULL,
    owner_id text NOT NULL,
    name text NOT NULL,
    status text NOT NULL,
    description text NOT NULL,
    total_quantity INT NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    primary key(id)
);

CREATE TABLE IF NOT EXISTS item_locks(
    item_id text NOT NULL,
    locked_by text NOT NULL,
    quantity INT NOT NULL,
    PRIMARY KEY (item_id, locked_by),
    CONSTRAINT item_id_fk FOREIGN KEY (item_id)
        REFERENCES items (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

