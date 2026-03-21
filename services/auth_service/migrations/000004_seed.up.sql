INSERT INTO roles(role)
VALUES
    ('user'),
    ('admin'),
    ('super_admin');

INSERT INTO users(name, email, role_id, hashed_pass, created_at)
VALUES
	('test1', 'test1@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
	('test2', 'test2@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
	('test3', 'test3@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
	('test4', 'test4@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
	('test5', 'test5@gmail.com', 1,  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW());
