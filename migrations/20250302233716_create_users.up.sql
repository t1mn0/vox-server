CREATE TABLE users (
    login TEXT PRIMARY KEY,
    username TEXT NOT NULL, 
    email TEXT NOT NULL UNIQUE,
    encrypted_password TEXT NOT NULL 
)