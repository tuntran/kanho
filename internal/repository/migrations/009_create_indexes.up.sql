CREATE INDEX idx_cards_column_position ON cards(column_id, position);
CREATE INDEX idx_cards_board ON cards(board_id) WHERE archived_at IS NULL;
CREATE INDEX idx_cards_due_date ON cards(due_date) WHERE due_date IS NOT NULL AND archived_at IS NULL;
CREATE INDEX idx_cards_search ON cards USING GIN (search_vector);
CREATE INDEX idx_activities_board ON activities(board_id, created_at DESC);
CREATE INDEX idx_activities_workspace ON activities(workspace_id, created_at DESC);
