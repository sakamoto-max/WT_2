INSERT INTO body_parts(name)
VALUES 
	('arms'),
	('back'),
	('cardio'),
	('chest'),
	('core'),
	('full_body'),
	('legs'),
	('olympic'),
	('other'),
	('shoulders');


INSERT INTO equipment(name)
VALUES 
	('barbell'),
	('dumbbell'),
	('machine'),
	('other'),
	('body_weight'),
	('assited_body_weight'),
	('reps_only'),
	('cardio'),
	('duration');


INSERT INTO exercises (name, created_by, body_part_id, equipment_id)
VALUES
    ('barbell_bench_press', NULL,
        (SELECT id FROM body_parts WHERE name = 'chest'),
        (SELECT id FROM equipment WHERE name = 'barbell')
    ),
    ('incline_dumbbell_press', NULL,
        (SELECT id FROM body_parts WHERE name = 'chest'),
        (SELECT id FROM equipment WHERE name = 'dumbbell')
    ),
    ('push_ups', NULL,
        (SELECT id FROM body_parts WHERE name = 'chest'),
        (SELECT id FROM equipment WHERE name = 'body_weight')
    ),
    ('pull_ups', NULL,
        (SELECT id FROM body_parts WHERE name = 'back'),
        (SELECT id FROM equipment WHERE name = 'body_weight')
    ),
    ('lat_pulldown', NULL,
        (SELECT id FROM body_parts WHERE name = 'back'),
        (SELECT id FROM equipment WHERE name = 'machine')
    ),
    ('seated_cable_row', NULL,
        (SELECT id FROM body_parts WHERE name = 'back'),
        (SELECT id FROM equipment WHERE name = 'machine')
    ),
    ('barbell_squats', NULL,
        (SELECT id FROM body_parts WHERE name = 'legs'),
        (SELECT id FROM equipment WHERE name = 'barbell')
    ),
    ('leg_press', NULL,
        (SELECT id FROM body_parts WHERE name = 'legs'),
        (SELECT id FROM equipment WHERE name = 'machine')
    ),
    ('lunges', NULL,
        (SELECT id FROM body_parts WHERE name = 'legs'),
        (SELECT id FROM equipment WHERE name = 'body_weight')
    ),
    ('deadlift', NULL,
        (SELECT id FROM body_parts WHERE name = 'back'),
        (SELECT id FROM equipment WHERE name = 'barbell')
    ),
    ('shoulder_press', NULL,
        (SELECT id FROM body_parts WHERE name = 'shoulders'),
        (SELECT id FROM equipment WHERE name = 'dumbbell')
    ),
    ('lateral_raises', NULL,
        (SELECT id FROM body_parts WHERE name = 'shoulders'),
        (SELECT id FROM equipment WHERE name = 'dumbbell')
    ),
    ('front_raises', NULL,
        (SELECT id FROM body_parts WHERE name = 'shoulders'),
        (SELECT id FROM equipment WHERE name = 'dumbbell')
    ),
    ('bicep_curls', NULL,
        (SELECT id FROM body_parts WHERE name = 'arms'),
        (SELECT id FROM equipment WHERE name = 'dumbbell')
    ),
    ('tricep_pushdown', NULL,
        (SELECT id FROM body_parts WHERE name = 'arms'),
        (SELECT id FROM equipment WHERE name = 'machine')
    ),
    ('hammer_curls', NULL,
        (SELECT id FROM body_parts WHERE name = 'arms'),
        (SELECT id FROM equipment WHERE name = 'dumbbell')
    ),
    ('crunches', NULL,
        (SELECT id FROM body_parts WHERE name = 'core'),
        (SELECT id FROM equipment WHERE name = 'body_weight')
    ),
    ('plank', NULL,
        (SELECT id FROM body_parts WHERE name = 'core'),
        (SELECT id FROM equipment WHERE name = 'duration')
    ),
    ('hanging_leg_raise', NULL,
        (SELECT id FROM body_parts WHERE name = 'core'),
        (SELECT id FROM equipment WHERE name = 'body_weight')
    ),
    ('jump_rope', NULL,
        (SELECT id FROM body_parts WHERE name = 'cardio'),
        (SELECT id FROM equipment WHERE name = 'cardio')
    ),
    ('treadmill_run', NULL,
        (SELECT id FROM body_parts WHERE name = 'cardio'),
        (SELECT id FROM equipment WHERE name = 'cardio')
    ),
    ('cycling', NULL,
        (SELECT id FROM body_parts WHERE name = 'cardio'),
        (SELECT id FROM equipment WHERE name = 'cardio')
    ),
    ('clean_and_jerk', NULL,
        (SELECT id FROM body_parts WHERE name = 'olympic'),
        (SELECT id FROM equipment WHERE name = 'barbell')
    ),
    ('snatch', NULL,
        (SELECT id FROM body_parts WHERE name = 'olympic'),
        (SELECT id FROM equipment WHERE name = 'barbell')
    ),
    ('burpees', NULL,
        (SELECT id FROM body_parts WHERE name = 'full_body'),
        (SELECT id FROM equipment WHERE name = 'body_weight')
    ),
    ('mountain_climbers', NULL,
        (SELECT id FROM body_parts WHERE name = 'full_body'),
        (SELECT id FROM equipment WHERE name = 'body_weight')
    ),
    ('kettlebell_swings', NULL,
        (SELECT id FROM body_parts WHERE name = 'full_body'),
        (SELECT id FROM equipment WHERE name = 'other')
    ),
    ('machine_chest_fly', NULL,
        (SELECT id FROM body_parts WHERE name = 'chest'),
        (SELECT id FROM equipment WHERE name = 'machine')
    ),
    ('leg_curl_machine', NULL,
        (SELECT id FROM body_parts WHERE name = 'legs'),
        (SELECT id FROM equipment WHERE name = 'machine')
    ),
    ('rowing_machine', NULL,
        (SELECT id FROM body_parts WHERE name = 'cardio'),
        (SELECT id FROM equipment WHERE name = 'cardio')
    );
