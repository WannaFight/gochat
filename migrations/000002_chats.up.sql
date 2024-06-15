CREATE TABLE IF NOT EXISTS chats (
    uuid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(128) NOT NULL DEFAULT '',
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chat_members (
    id bigserial PRIMARY KEY,
    chat_uuid uuid REFERENCES chats(uuid),
    user_id bigserial REFERENCES users(id),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    is_owner boolean DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id bigserial PRIMARY KEY,
    text varchar(2048) NOT NULL DEFAULT '',
    sent_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    chat_member_id bigint REFERENCES chat_members(id)
);
