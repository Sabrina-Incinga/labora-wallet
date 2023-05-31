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
    track_type varchar(50) NOT NULL,
    request_status varchar(50) NOT NULL,
    CONSTRAINT customer_id FOREIGN KEY (customer_id)
        REFERENCES public.customer (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE IF EXISTS public.wallet_tracker
    OWNER to postgres;

ALTER TABLE IF EXISTS public.wallet
DROP CONSTRAINT IF EXISTS only_positive_balance_allowed;

ALTER TABLE IF EXISTS public.wallet
ADD CONSTRAINT only_positive_balance_allowed CHECK (balance >= 0);

ALTER TABLE IF EXISTS public.wallet
ALTER COLUMN creation_date TYPE timestamp;

ALTER TABLE IF EXISTS public.wallet_tracker
ALTER COLUMN record_date TYPE timestamp;

CREATE TABLE IF NOT EXISTS public.wallet_movement
(
    id SERIAL PRIMARY KEY,
    sender_wallet_id integer NOT NULL,
    receiver_wallet_id integer,
    movement_date timestamp NOT NULL,
    movement_type varchar(10) NOT NULL,
    amount double precision NOT NULL,
    CONSTRAINT sender_wallet_id FOREIGN KEY (sender_wallet_id)
        REFERENCES public.wallet (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);