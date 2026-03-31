DROP TRIGGER IF EXISTS activity_notify_trigger ON activities;
DROP FUNCTION IF EXISTS notify_activity();
DROP FUNCTION IF EXISTS next_card_number(UUID);
DROP TABLE IF EXISTS activities;
