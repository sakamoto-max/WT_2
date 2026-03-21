CREATE TABLE outbox(
    id SERIAL PRIMARY KEY,
    target_service TEXT NOT NULL,
    task TEXT NOT NULL,
    status TEXT NOT NULL,
    payload JSONB NOT NULL,
    CREATED_AT TIMESTAMPTZ NOT NULL,
    CONSTRAINT check_status CHECK(status IN ('completed', 'not_completed')) 
)