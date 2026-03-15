CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	role_id int NOT NULL,
	hashed_pass TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL, 
	updated_at TIMESTAMPTZ,
	FOREIGN KEY (role_id) REFERENCES roles(id)
)

-- seed
-- INSERT INTO users(name, email, role_id, hashed_pass, created_at)
-- VALUES
-- 	('test1', 'test1@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
-- 	('test2', 'test2@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
-- 	('test3', 'test3@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
-- 	('test4', 'test4@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
-- 	('test5', 'test5@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW());
