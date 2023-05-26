CREATE TABLE IF NOT EXISTS public."customer"
(
    id SERIAL PRIMARY KEY,
    first_name varchar(150) NOT NULL,
    last_name varchar(150) NOT NULL,
    national_identity_number varchar(20) NOT NULL,
    national_identity_type varchar(10) NOT NULL,
    country_id varchar(10) NOT NULL
);

ALTER TABLE IF EXISTS public."customer"
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.wallet
(
    id SERIAL PRIMARY KEY,
    customer_id integer NOT NULL,
    creation_date date NOT NULL,
    balance double precision NOT NULL,
    CONSTRAINT customer_id FOREIGN KEY (customer_id)
        REFERENCES public.customer (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE IF EXISTS public.wallet
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.wallet_tracker
(
    id SERIAL PRIMARY KEY,
    customer_id integer NOT NULL,
    record_date date NOT NULL,
    creation_status varchar(10) NOT NULL,
    CONSTRAINT customer_id FOREIGN KEY (customer_id)
        REFERENCES public.customer (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE IF EXISTS public.wallet_tracker
    OWNER to postgres;