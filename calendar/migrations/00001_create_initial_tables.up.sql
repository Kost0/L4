CREATE TABLE events (
    event_id UUID PRIMARY KEY,
    user_id UUID,
    date TIMESTAMP,
    event VARCHAR(200)
);

CREATE TABLE archive (
    event_id UUID PRIMARY KEY,
    user_id UUID,
    date TIMESTAMP,
    event VARCHAR(200)
);