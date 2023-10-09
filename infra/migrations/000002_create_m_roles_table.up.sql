CREATE TABLE m_roles (
   id SERIAL PRIMARY KEY,
   name VARCHAR(128) NOT NULL,
   created_at DECIMAL NOT NULL,
   created_by VARCHAR(64),
   updated_at DECIMAL NOT NULL,
   updated_by VARCHAR(64),
   deleted_at DECIMAL,
   deleted_by VARCHAR(64)
);

INSERT INTO m_roles
("name", created_at, created_by, updated_at, updated_by, deleted_at, deleted_by)
VALUES('user_repository', 1684768468, NULL, 1684768468, NULL, NULL, NULL);