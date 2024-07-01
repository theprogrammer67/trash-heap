DROP TRIGGER IF EXISTS user_insert ON users;
DROP TRIGGER IF EXISTS event_insert ON events;
DROP FUNCTION IF EXISTS on_insert_user;
DROP FUNCTION IF EXISTS on_insert_event;
DROP TABLE IF EXISTS events;

CREATE TABLE IF NOT EXISTS events (
	id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
	event_type TEXT,
    event_data TEXT,
 	created_at timestamp NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION on_insert_event()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $$
BEGIN
    PERFORM pg_notify('event', json_build_object(
            'id', NEW.id,
            'event_type', NEW.event_type,
            'event_data', NEW.event_data
        )::text);
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER event_insert
AFTER INSERT ON events
FOR EACH ROW EXECUTE FUNCTION on_insert_event();

CREATE OR REPLACE FUNCTION on_insert_user()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $$
BEGIN
    INSERT INTO events (event_type, event_data)
    VALUES
        ('insert_user', row_to_json(NEW)::text);

    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER user_insert
AFTER INSERT ON users
FOR EACH ROW EXECUTE FUNCTION on_insert_user();
