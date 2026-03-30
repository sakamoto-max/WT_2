CREATE TABLE exercises(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name TEXT NOT NULL,
	created_by INTEGER,
	body_part_id UUID REFERENCES body_parts(id),
	equipment_id UUID REFERENCES equipment(id)
);
