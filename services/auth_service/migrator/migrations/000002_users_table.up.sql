CREATE TABLE users(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	role_id UUID NOT NULL,
	hashed_pass TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, 
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE INDEX email_ind ON users(email);
CREATE UNIQUE INDEX one_user_one_email ON users(name, email);
