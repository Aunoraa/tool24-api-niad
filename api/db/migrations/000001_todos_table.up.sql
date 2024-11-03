<<<<<<< HEAD
=======
-- +goose Up
>>>>>>> 47370263c8b7d4210ea76d37ab4be1150f31b32a
CREATE TABLE todo (
    id SERIAL VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    descripsion TEXT,
    done BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    done_at TIMESTAMP
);

-- add_done_at_column.up.sql
ALTER TABLE todo ADD COLUMN done_at TIMESTAMP;