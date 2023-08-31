CREATE TABLE m_apps (
   id VARCHAR(64) PRIMARY KEY UNIQUE NOT NULL,
   name VARCHAR(128) NOT NULL,
   description TEXT NOT NULL,
   created_at DECIMAL NOT NULL,
   created_by VARCHAR(64),
   updated_at DECIMAL NOT NULL,
   updated_by VARCHAR(64),
   deleted_at DECIMAL,
   deleted_by VARCHAR(64)
);

INSERT INTO auth.m_apps
(id, "name", description, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by)
VALUES('5410801c-faaf-4776-95be-56472e044820', 'kakeibo', 'aplikasi manajemen keuangan', 1684768516, NULL, 1684768516, NULL, NULL, NULL);
