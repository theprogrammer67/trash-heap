CREATE TABLE IF NOT EXISTS events (
	id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
	payload	JSONB,
 	created_at timestamp NOT NULL DEFAULT now()   
)

CREATE OR REPLACE FUNCTION on_insert_event()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $$
BEGIN
    PERFORM pg_notify('event', row_to_json(NEW)::text);
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
--     PERFORM pg_notify('test', new.name);
	INSERT INTO events (payload)
	VALUES
		(json_build_object('type', 'insert_user', 'data', row_to_json(NEW))::JSONB);
	
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER event_user
AFTER INSERT ON users
FOR EACH ROW EXECUTE FUNCTION on_insert_user();
