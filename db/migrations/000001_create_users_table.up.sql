CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	name varchar,
	email varchar,
	username varchar,
	password varchar,
	CONSTRAINT email_unique UNIQUE (email),
	CONSTRAINT username_unique UNIQUE (username)
);