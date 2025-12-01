-- Users table for player authentication and profiles
CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   password VARCHAR (100) NOT NULL,
   email VARCHAR (50) UNIQUE NOT NULL
);

-- Index for faster email lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

