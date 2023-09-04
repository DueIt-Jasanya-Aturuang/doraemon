CREATE TABLE m_tokens
(
    id              VARCHAR(64) NOT NULL UNIQUE PRIMARY KEY,
    user_id         VARCHAR(64),
    app_id          VARCHAR(64),
    token           TEXT,
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