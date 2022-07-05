CREATE TABLE IF NOT EXISTS "balance"
(
    id SERIAL UNIQUE NOT NULL,
    userId text UNIQUE,
    balance float CHECK (balance > 0),
    CONSTRAINT "pk_Address" PRIMARY KEY ("id")
);



