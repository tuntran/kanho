CREATE TABLE cards (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id       UUID NOT NULL REFERENCES projects(id),
    board_id         UUID NOT NULL REFERENCES boards(id),
    column_id        UUID NOT NULL REFERENCES columns(id),
    card_number      INTEGER NOT NULL,
    title            TEXT NOT NULL,
    description      TEXT,
    description_json JSONB,
    priority         TEXT NOT NULL DEFAULT 'none'
                         CHECK (priority IN ('none','low','medium','high','urgent')),
    position         TEXT NOT NULL COLLATE "C",
    start_date       TIMESTAMPTZ,
    due_date         TIMESTAMPTZ,
    custom_fields    JSONB NOT NULL DEFAULT '{}',
    archived_at      TIMESTAMPTZ,
    created_by       UUID NOT NULL REFERENCES users(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(project_id, card_number),
    UNIQUE(column_id, position)
);

ALTER TABLE cards ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(description, '')), 'B')
) STORED;
