CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    username citext UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS chats (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(128) DEFAULT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chat_members (
    id bigserial PRIMARY KEY,
    chat_id uuid REFERENCES chats(id),
    user_id bigint REFERENCES users(id),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    is_owner boolean DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id bigserial PRIMARY KEY,
    text varchar(2048) DEFAULT NULL,
    chat_id uuid REFERENCES chats(id),
    user_id bigint REFERENCES users(id),
    sent_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);
