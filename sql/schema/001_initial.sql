-- +goose Up
CREATE TABLE users (
   id UUID PRIMARY KEY,
   username VARCHAR(50) NOT NULL,
   email VARCHAR(50) NOT NULL,
   hashed_password VARCHAR(60) NOT NULL,
   first_name VARCHAR(50) NOT NULL,
   last_name VARCHAR(50) NOT NULL,
   phone_number VARCHAR(20) NULL,
   date_of_birth DATE NULL,
   gender VARCHAR(10) NULL,
   shipping_address VARCHAR(100) NULL,
   billing_address VARCHAR(100) NULL,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   last_login TIMESTAMP NULL,
   account_status VARCHAR(10) NOT NULL DEFAULT 'active',
   user_role VARCHAR(10) NOT NULL DEFAULT 'customer',
   profile_picture VARCHAR(100) NULL,
   two_factor_auth BOOLEAN NOT NULL DEFAULT FALSE,
   UNIQUE (username),
   UNIQUE (email),
   CHECK (account_status IN ('active', 'inactive', 'suspended', 'deleted')),
   CHECK (user_role IN ('user', 'admin', 'superadmin'))
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX refresh_tokens_user_id_idx ON refresh_tokens (user_id);
CREATE INDEX refresh_tokens_token_idx ON refresh_tokens (token);
CREATE INDEX refresh_tokens_expires_at_idx ON refresh_tokens (expires_at);
CREATE INDEX refresh_tokens_revoked_at_idx ON refresh_tokens (revoked_at);

-- +goose Down
DROP TABLE refresh_tokens;
DROP TABLE users;