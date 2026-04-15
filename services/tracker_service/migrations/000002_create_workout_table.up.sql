CREATE TABLE workout (
	tracker_id UUID NOT NULL,
	exercise_id UUID NOT NULL,
	set_number INTEGER NOT NULL,
	weight INTEGER,
	reps INTEGER,
	FOREIGN KEY (tracker_id) REFERENCES tracker(id)
)

