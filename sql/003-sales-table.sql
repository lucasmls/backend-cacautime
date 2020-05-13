CREATE TABLE sales (
  id SERIAL PRIMARY KEY,
  customer_id integer REFERENCES customers(id) ON DELETE CASCADE ON UPDATE CASCADE,
  duty_id integer REFERENCES duties(id) ON DELETE CASCADE ON UPDATE CASCADE
);

COMMENT ON CONSTRAINT customer_fk ON sales IS 'The customer who bought the candy';
COMMENT ON CONSTRAINT duty_fk ON sales IS 'In which duty the candy was bought';