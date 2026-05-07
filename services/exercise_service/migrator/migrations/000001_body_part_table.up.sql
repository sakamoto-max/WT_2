CREATE TABLE body_parts(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 	name TEXT NOT NULL UNIQUE
);