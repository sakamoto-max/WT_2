INSERT INTO body_parts(name)
VALUES 
	('Arms'),
	('Back'),
	('Cardio'),
	('Chest'),
	('Core'),
	('Full Body'),
	('Legs'),
	('Olympic'),
	('Other'),
	('Shoulders');

INSERT INTO equipment(name)
VALUES 
	('Barbell'),
	('Dumbbell'),
	('Machine'),
	('Other'),
	('Body Weight'),
	('Assited Body Weight'),
	('Reps Only'),
	('Cardio'),
	('Duration');

INSERT INTO exercises (name, body_part_id, rest_time_in_seconds, equipment_id, created_at)
VALUES
    ('Barbell Bench Press', 4, 120, 1, NOW()),
    ('Incline Dumbbell Press', 4, 120, 2, NOW()),
    ('Push Ups', 4, 120, 5, NOW()),
    ('Pull Ups', 2, 120, 5, NOW()),
    ('Lat Pulldown', 2, 120, 3, NOW()),
    ('Seated Cable Row', 2, 120, 3, NOW()),
    ('Barbell Squats', 7, 120, 1, NOW()),
    ('Leg Press', 7, 120, 3, NOW()),
    ('Lunges', 7, 120, 5, NOW()),
    ('Deadlift', 2, 120, 1, NOW()),
    ('Shoulder Press', 10, 120, 2, NOW()),
    ('Lateral Raises', 10, 120, 2, NOW()),
    ('Front Raises', 10, 120, 2, NOW()),
    ('Bicep Curls', 1, 120, 2, NOW()),
    ('Tricep Pushdown', 1, 120, 3, NOW()),
    ('Hammer Curls', 1, 120, 2, NOW()),
    ('Crunches', 5, 120, 5, NOW()),
    ('Plank', 5, 120, 9, NOW()),
    ('Hanging Leg Raise', 5, 120, 5, NOW()),
    ('Jump Rope', 3, 120, 8, NOW()),
    ('Treadmill Run', 3, 120, 8, NOW()),
    ('Cycling', 3, 120, 8, NOW()),
    ('Clean and Jerk', 8, 120, 1, NOW()),
    ('Snatch', 8, 120, 1, NOW()),
    ('Burpees', 6, 120, 5, NOW()),
    ('Mountain Climbers', 6, 120, 5, NOW()),
    ('Kettlebell Swings', 6, 120, 4, NOW()),
    ('Machine Chest Fly', 4, 120, 3, NOW()),
    ('Leg Curl Machine', 7, 120, 3, NOW()),
    ('Rowing Machine', 3, 120, 8, NOW());