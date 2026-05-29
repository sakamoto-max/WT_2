CREATE TABLE PUSH_TO_QUEUE_FAILED(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	target_db_id UUID NOT NULL,
	target_service TEXT NOT NULL,
	number_of_tries INTEGER NOT NULL,
	status TEXT NOT NULL,
	reason TEXT NOT NULL
)