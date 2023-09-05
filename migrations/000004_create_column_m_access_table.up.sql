CREATE TABLE m_access(
    id SERIAL primary key,
    role_id INT,
    user_id VARCHAR(64),
    app_id VARCHAR(64),
    access_endpoint JSONB,
    created_at DECIMAL NOT NULL,
    created_by VARCHAR(64),
    updated_at DECIMAL NOT NULL,
    updated_by VARCHAR(64),
    deleted_at DECIMAL,
    deleted_by VARCHAR(64),
    CONSTRAINT fk_m_role
      FOREIGN KEY(role_id)
	   REFERENCES m_roles(id)
       ON DELETE CASCADE
       ON UPDATE CASCADE,
    CONSTRAINT fk_m_user
      FOREIGN KEY(user_id)
	   REFERENCES m_users(id)
       ON DELETE CASCADE
       ON UPDATE CASCADE,
    CONSTRAINT fk_m_app
      FOREIGN KEY(app_id)
	   REFERENCES m_apps(id)
       ON DELETE CASCADE
       ON UPDATE CASCADE
);
