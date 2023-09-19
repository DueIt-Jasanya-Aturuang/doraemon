CREATE TABLE m_tokens
(
    id            SERIAL primary key,
    user_id       VARCHAR(64),
    app_id        VARCHAR(64),
    remember_me   BOOLEAN NOT NULL DEFAULT FALSE,
    access_token  VARCHAR(255),
    refresh_token VARCHAR(255),
    CONSTRAINT fk_m_user
        FOREIGN KEY (user_id)
            REFERENCES m_users (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE,
    CONSTRAINT fk_m_app
        FOREIGN KEY (app_id)
            REFERENCES m_apps (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
);