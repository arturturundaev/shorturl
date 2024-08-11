CREATE TABLE IF NOT EXISTS url
(
    id uuid NOT NULL,
    url_full character varying NOT NULL,
    url_short character varying NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT name UNIQUE (url_full)
);