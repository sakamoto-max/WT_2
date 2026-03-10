
INSERT INTO roles(role)
VALUES 
	('user'),
	('admin');



INSERT INTO users(name, email, hashed_pass, created_at)
VALUES
	('test1', 'test1@gmail.com', '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
	('test2', 'test2@gmail.com', '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
	('test3', 'test3@gmail.com', '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
	('test4', 'test4@gmail.com', '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW()),
	('test5', 'test5@gmail.com', '$2a$10$iHlzHAzGLPIck0s6oTg4geq26ImcQzfG9az6Jqr6Dlai./we0E9WG', NOW());


INSERT INTO user_roles(user_id, role_id) 
VALUES 	
	(1, 1),
	(2, 1),
	(3, 1),
	(4, 1),
	(5, 1);