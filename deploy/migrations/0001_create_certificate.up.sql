CREATE TABLE public.certificates (
    certificate_id UUID NOT NULL,
    artist_id UUID NOT NULL,
    artwork_type VARCHAR NOT NULL,
    title VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT (now() at time zone 'utc') NOT NULL,
    modified_at TIMESTAMP DEFAULT (now() at time zone 'utc') NOT NULL,
    deleted BOOLEAN DEFAULT FALSE NOT NULL,
    CONSTRAINT certificate_id_pk PRIMARY KEY (certificate_id)
);