ALTER TABLE events DROP COLUMN IF EXISTS tournament;
ALTER TABLE events DROP COLUMN IF EXISTS started_at;
DELETE FROM event_user_roles WHERE role = 'substitute';
ALTER TABLE event_user_roles DROP CONSTRAINT event_user_roles_role_check;
ALTER TABLE event_user_roles ADD CONSTRAINT event_user_roles_role_check
    CHECK (role IN ('player', 'referee', 'observer'));
