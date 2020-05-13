-- Table Definition ----------------------------------------------
CREATE TABLE customers (
  id SERIAL PRIMARY KEY,
  name character varying(40) NOT NULL,
  phone character varying(11) NOT NULL,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  updated_at timestamp without time zone NOT NULL DEFAULT now()
);

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX customers_pkey ON customers(id int4_ops);

-- Triggers -------------------------------------------------------
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();