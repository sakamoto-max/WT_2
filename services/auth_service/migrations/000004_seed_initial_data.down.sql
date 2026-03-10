DELETE FROM user_roles
WHERE user_id IN (1, 2, 3, 4, 5);

DELETE FROM users
WHERE name IN ('test1', 'test2', 'test3', 'test4', 'test5');

DELETE FROM user_roles
WHERE role IN ('user', 'admin');