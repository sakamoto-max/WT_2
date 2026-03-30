CREATE TABLE outbox(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_service TEXT NOT NULL,
    task TEXT NOT NULL,
    status TEXT NOT NULL,
    payload JSONB NOT NULL,
    CREATED_AT TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    number_of_tries INTEGER,
    CONSTRAINT check_status CHECK(status IN ('completed', 'not_completed', 'pending')) 
)

