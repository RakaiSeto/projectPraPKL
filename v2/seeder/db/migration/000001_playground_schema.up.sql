CREATE TABLE public.product (
   id SERIAL primary key,
   name VARCHAR NOT NULL,
   description VARCHAR,
   price integer NOT NULL
);
GRANT UPDATE, TRUNCATE, REFERENCES, INSERT, DELETE, TRIGGER, SELECT ON TABLE public.product TO productschema;

CREATE TABLE public.user (
   id SERIAL PRIMARY KEY,
   uname VARCHAR,
   email VARCHAR,
   password VARCHAR,
   role VARCHAR
);
GRANT UPDATE, TRUNCATE, REFERENCES, INSERT, DELETE, TRIGGER, SELECT ON TABLE public.user TO userschema;

CREATE TABLE public.order (
   id SERIAL PRIMARY KEY DEFAULT,
   userid integer not null,
   productid integer not null,
   quantity integer not null,
   totalprice integer default 0
);
GRANT UPDATE, TRUNCATE, REFERENCES, INSERT, DELETE, TRIGGER, SELECT ON TABLE public.order TO orderschema;
