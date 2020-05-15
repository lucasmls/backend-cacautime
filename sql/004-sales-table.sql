-- Table Definition ----------------------------------------------
CREATE TABLE sales (
  id SERIAL PRIMARY KEY,
  customer_id integer REFERENCES customers(id) ON DELETE CASCADE ON UPDATE CASCADE,
  duty_id integer REFERENCES duties(id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  updated_at timestamp without time zone NOT NULL DEFAULT now(),
  candy_id integer NOT NULL REFERENCES candies(id) ON DELETE CASCADE ON UPDATE CASCADE,
  status text NOT NULL DEFAULT 'paid'::text,
  payment_method text NOT NULL DEFAULT 'transfer'::text
);

-- Comments -------------------------------------------------------
COMMENT ON COLUMN sales.status IS 'paid/not_paid';
COMMENT ON COLUMN sales.payment_method IS 'money/transfer/scheduled';
COMMENT ON CONSTRAINT customer_fk ON sales IS 'The customer who bought the candy';
COMMENT ON CONSTRAINT duty_fk ON sales IS 'In which duty the candy was bought';
COMMENT ON CONSTRAINT candy_fk ON sales IS 'The candy that was sold';

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX sales_pkey ON sales(id int4_ops);

-- Triggers -------------------------------------------------------
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON sales
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();