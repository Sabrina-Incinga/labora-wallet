CREATE TABLE IF NOT EXISTS public."customer"
(
    id SERIAL PRIMARY KEY,
    first_name varchar(150) NOT NULL,
    last_name varchar(150) NOT NULL,
    national_identity_number varchar(20) NOT NULL,
    national_identity_type varchar(10) NOT NULL,
    country_id varchar(10) NOT NULL,
    CONSTRAINT unique_national_identity UNIQUE (national_identity_number, national_identity_type, country_id )

);

CREATE INDEX IF NOT EXISTS idx_unique_national_identity ON public."customer" (national_identity_number, national_identity_type, country_id);

ALTER TABLE IF EXISTS public."customer"
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.wallet
(
    id SERIAL PRIMARY KEY,
    customer_id integer NOT NULL,
    wallet_number varchar(22) CHECK (length(wallet_number) = 22) NOT NULL,
    creation_date date NOT NULL,
    balance double precision NOT NULL,
    CONSTRAINT customer_id FOREIGN KEY (customer_id)
        REFERENCES public.customer (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT unique_wallet_number UNIQUE (wallet_number),
    CONSTRAINT unique_customer_id UNIQUE (customer_id)
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