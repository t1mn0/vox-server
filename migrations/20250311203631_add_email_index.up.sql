ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);

CREATE INDEX idx_users_email ON users (email);