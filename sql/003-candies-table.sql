-- Table Definition ----------------------------------------------
CREATE TABLE candies (
  id SERIAL PRIMARY KEY,
  name text NOT NULL,
  price integer NOT NULL,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  updated_at timestamp without time zone NOT NULL DEFAULT now()
);

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX candies_pkey ON candies(id int4_ops);

-- Triggers -------------------------------------------------------
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON candies
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();