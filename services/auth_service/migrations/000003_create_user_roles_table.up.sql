CREATE TABLE user_roles(
	user_id SERIAL NOT NULL,
	role_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
)