-- Table Definition ----------------------------------------------
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name character varying(40) NOT NULL,
  email character varying(40) NOT NULL,
  password character varying(100) NOT NULL,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  updated_at timestamp without time zone NOT NULL DEFAULT now()
);

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX users_pkey ON users(id int4_ops);

-- Triggers -------------------------------------------------------
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();