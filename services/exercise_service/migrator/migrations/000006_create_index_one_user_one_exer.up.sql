CREATE UNIQUE INDEX one_user_one_exercise 
ON exercises(created_by, name, body_part_id, equipment_id);