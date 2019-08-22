
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied


create table configurations(
	id serial not null,
	created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  min_value_buy numeric(12,2),
);
ALTER TABLE ONLY configurations ADD CONSTRAINT configurations_pkey PRIMARY KEY (id);

CREATE TABLE users(
	id serial not null,
	created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  deleted_at timestamp with time zone,
  last_login timestamp with time zone,
  hashed_password bytea NOT NULL,
  password varchar(255),
  image varchar(255),
  phone varchar(255),
	name varchar(255) not null,
	username varchar(50) not null,
	email varchar(100) not null,
	balance numeric(12,2) default 0,
	admin boolean not null default false,
	ban boolean not null default false,
);

ALTER TABLE ONLY users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
CREATE UNIQUE INDEX idx_users_email ON users USING btree (email);
CREATE UNIQUE INDEX idx_users_username ON users USING btree (username);
CREATE UNIQUE INDEX idx_lower_case_username ON users ((lower(username)));


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE configurations;
DROP TABLE users;