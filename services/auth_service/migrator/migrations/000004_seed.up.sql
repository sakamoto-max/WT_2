INSERT INTO roles(role)
VALUES
    ('user'),
    ('admin'),
    ('super_admin');

INSERT INTO users(name, email, role_id, hashed_pass)
VALUES
	('test1', 'test1@gmail.com', (SELECT ID FROM ROLES WHERE role = 'user'),  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG'),
	('test2', 'test2@gmail.com',  (SELECT ID FROM ROLES WHERE role = 'user'),  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG'),
	('test3', 'test3@gmail.com',  (SELECT ID FROM ROLES WHERE role = 'user'),  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG'),
	('test4', 'test4@gmail.com',  (SELECT ID FROM ROLES WHERE role = 'user'),  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG'),
	('test5', 'test5@gmail.com',  (SELECT ID FROM ROLES WHERE role = 'user'),  '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG');
