CREATE TABLE IF NOT EXISTS chats (
    uuid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chat_memebers (
    id bigserial PRIMARY KEY,
    chat_uuid uuid REFERENCES chats(uuid),
    user_id bigserial REFERENCES users(id),
    is_owner boolean DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id bigserial PRIMARY KEY,
    text text NOT NULL DEFAULT '',
    sent_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    chat_member_id bigint REFERENCES chat_memebers(id)
);
