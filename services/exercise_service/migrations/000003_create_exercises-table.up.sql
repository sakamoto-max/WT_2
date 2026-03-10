CREATE TABLE exercises(
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	body_part_id INTEGER NOT NULL,
	rest_time_in_seconds INTEGER NOT NULL,
	equipment_id INTEGER NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ,
	FOREIGN KEY (body_part_id) REFERENCES body_parts(id),
	FOREIGN KEY (equipment_id) REFERENCES equipment(id)
)
