create type gender_enum as enum('male', 'female', 'undefined');

CREATE TABLE m_users (
   id VARCHAR(64) NOT NULL UNIQUE PRIMARY KEY,
   fullname VARCHAR(32) NOT NULL,
   gender gender_enum default 'undefined',
   image VARCHAR(255) default 'default-male.png',
   email VARCHAR(55) NOT NULL,
   username VARCHAR(22) NOT NULL UNIQUE,
   password VARCHAR(255) NOT NULL,
   phone_number VARCHAR(14) UNIQUE,
   email_verified_at BOOLEAN NOT NULL DEFAULT FALSE,
   created_at DECIMAL NOT NULL,
   created_by VARCHAR(64),
   updated_at DECIMAL NOT NULL,
   updated_by VARCHAR(64),
   deleted_at DECIMAL,
   deleted_by VARCHAR(64)
);
