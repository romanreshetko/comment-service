CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    review_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    comment_text TEXT,
    status TEXT NOT NULL CHECK (status IN ('published', 'moderating', 'blocked', 'reported', 'blocked_reported', 'undefined', 'moderation_error')),
    created_at TIMESTAMP,
    prev_comment_id BIGINT NULL
);

CREATE INDEX idx_review_comments
ON comments (review_id)
WHERE status = 'published';