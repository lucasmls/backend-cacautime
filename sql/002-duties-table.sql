-- Table Definition ----------------------------------------------
CREATE TABLE duties (
  id SERIAL PRIMARY KEY,
  date date NOT NULL,
  candy_quantity integer NOT NULL,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  updated_at timestamp without time zone NOT NULL DEFAULT now()
);

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX sales_day_pkey ON duties(id int4_ops);

-- Triggers -------------------------------------------------------
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON duties
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();