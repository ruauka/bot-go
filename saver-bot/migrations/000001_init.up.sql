CREATE TABLE IF NOT EXISTS event
(
    id serial PRIMARY KEY,
    date varchar(255),
    type varchar(255),
    username varchar(255),
    telega_id varchar(255)
);