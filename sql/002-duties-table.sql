CREATE TABLE duties (
  id integer DEFAULT nextval('sales_day_id_seq'::regclass) PRIMARY KEY,
  date date NOT NULL,
  candy_quantity integer NOT NULL
);