CREATE TABLE activities (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    actor_id     UUID NOT NULL REFERENCES users(id),
    action       TEXT NOT NULL,
    target_type  TEXT NOT NULL,
    target_id    UUID NOT NULL,
    project_id   UUID REFERENCES projects(id),
    board_id     UUID REFERENCES boards(id),
    data         JSONB NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION next_card_number(p_project_id UUID)
RETURNS INTEGER AS $$
DECLARE v_number INTEGER;
BEGIN
    UPDATE projects SET card_counter = card_counter + 1
    WHERE id = p_project_id
    RETURNING card_counter INTO v_number;
    RETURN v_number;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION notify_activity() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.board_id IS NOT NULL THEN
        PERFORM pg_notify('board_events',
            json_build_object(
                'event_id', NEW.id::TEXT,
                'board_id', NEW.board_id::TEXT,
                'action', NEW.action,
                'target_type', NEW.target_type,
                'target_id', NEW.target_id::TEXT,
                'actor_id', NEW.actor_id::TEXT,
                'data', NEW.data
            )::TEXT);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER activity_notify_trigger
    AFTER INSERT ON activities
    FOR EACH ROW EXECUTE FUNCTION notify_activity();
