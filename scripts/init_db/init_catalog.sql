CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name varchar(250) NOT NULL,
    description varchar(250) NOT NULL,
    price numeric NOT NULL,
    price_code varchar(3) NOT NULL
)