CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  user_id BIGINT UNIQUE NOT NULL,
  username VARCHAR(50) NOT NULL,
  display VARCHAR(50),
  first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_user_platform
    UNIQUE (user_id, username)
);

CREATE INDEX IF NOT EXISTS idx_user_id
    ON users (user_id);

CREATE TABLE IF NOT EXISTS channels (
  user_id INT UNIQUE NOT NULL
    REFERENCES users (user_id)
    ON DELETE CASCADE,
  joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_channel_id
    ON channels (user_id);
