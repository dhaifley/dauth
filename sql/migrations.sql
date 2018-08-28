-- ============================================================================
-- dauth migrations
-- Create the schema for the dauth database.
-- Author: David Haifley
-- Created: 2018-08-01
-- ============================================================================

-- Table: public.perm

-- DROP TABLE public.perm;

CREATE TABLE public.perm
(
    id bigint NOT NULL,
    service character varying(32) COLLATE pg_catalog."default" NOT NULL,
    name character varying(32) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT perm_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.perm
    OWNER to dauth;

-- Index: ix_perm_id

-- DROP INDEX public.ix_perm_id;

CREATE UNIQUE INDEX ix_perm_id
    ON public.perm USING btree
    (id)
    TABLESPACE pg_default;

ALTER TABLE public.perm
    CLUSTER ON ix_perm_id;

-- Index: ix_perm_name

-- DROP INDEX public.ix_perm_name;

CREATE INDEX ix_perm_name
    ON public.perm USING btree
    (name COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Index: ix_perm_service

-- DROP INDEX public.ix_perm_service;

CREATE INDEX ix_perm_service
    ON public.perm USING btree
    (service COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Index: ix_perm_service_name

-- DROP INDEX public.ix_perm_service_name;

CREATE UNIQUE INDEX ix_perm_service_name
    ON public.perm USING btree
    (service COLLATE pg_catalog."default", name COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Table: public.token

-- DROP TABLE public.token;

CREATE TABLE public.token
(
    id bigint NOT NULL,
    token character varying(255) COLLATE pg_catalog."default" NOT NULL,
    user_id bigint NOT NULL,
    created timestamp with time zone,
    expires timestamp with time zone,
    CONSTRAINT token_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.token
    OWNER to dauth;

-- Index: ix_token_created

-- DROP INDEX public.ix_token_created;

CREATE INDEX ix_token_created
    ON public.token USING btree
    (created)
    TABLESPACE pg_default;

-- Index: ix_token_expires

-- DROP INDEX public.ix_token_expires;

CREATE INDEX ix_token_expires
    ON public.token USING btree
    (expires)
    TABLESPACE pg_default;

-- Index: ix_token_id

-- DROP INDEX public.ix_token_id;

CREATE UNIQUE INDEX ix_token_id
    ON public.token USING btree
    (id)
    TABLESPACE pg_default;

ALTER TABLE public.token
    CLUSTER ON ix_token_id;

-- Index: ix_token_token

-- DROP INDEX public.ix_token_token;

CREATE UNIQUE INDEX ix_token_token
    ON public.token USING btree
    (token COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Index: ix_token_user_id

-- DROP INDEX public.ix_token_user_id;

CREATE INDEX ix_token_user_id
    ON public.token USING btree
    (user_id)
    TABLESPACE pg_default;

-- Table: public."user"

-- DROP TABLE public."user";

CREATE TABLE public."user"
(
    id bigint NOT NULL,
    "user" character varying(32) COLLATE pg_catalog."default" NOT NULL,
    pass character varying(128) COLLATE pg_catalog."default" NOT NULL,
    name character varying(64) COLLATE pg_catalog."default",
    email character varying(128) COLLATE pg_catalog."default",
    CONSTRAINT user_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public."user"
    OWNER to dauth;

-- Index: ix_user_email

-- DROP INDEX public.ix_user_email;

CREATE INDEX ix_user_email
    ON public."user" USING btree
    (email COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Index: ix_user_id

-- DROP INDEX public.ix_user_id;

CREATE UNIQUE INDEX ix_user_id
    ON public."user" USING btree
    (id)
    TABLESPACE pg_default;

ALTER TABLE public."user"
    CLUSTER ON ix_user_id;

-- Index: ix_user_name

-- DROP INDEX public.ix_user_name;

CREATE INDEX ix_user_name
    ON public."user" USING btree
    (name COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Index: ix_user_user

-- DROP INDEX public.ix_user_user;

CREATE UNIQUE INDEX ix_user_user
    ON public."user" USING btree
    ("user" COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Table: public.user_perm

-- DROP TABLE public.user_perm;

CREATE TABLE public.user_perm
(
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    perm_id bigint NOT NULL,
    CONSTRAINT user_perm_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.user_perm
    OWNER to dauth;

-- Index: ix_user_perm_id

-- DROP INDEX public.ix_user_perm_id;

CREATE UNIQUE INDEX ix_user_perm_id
    ON public.user_perm USING btree
    (id)
    TABLESPACE pg_default;

ALTER TABLE public.user_perm
    CLUSTER ON ix_user_perm_id;

-- Index: ix_user_perm_perm_id

-- DROP INDEX public.ix_user_perm_perm_id;

CREATE INDEX ix_user_perm_perm_id
    ON public.user_perm USING btree
    (perm_id)
    TABLESPACE pg_default;

-- Index: ix_user_perm_user_id

-- DROP INDEX public.ix_user_perm_user_id;

CREATE INDEX ix_user_perm_user_id
    ON public.user_perm USING btree
    (user_id)
    TABLESPACE pg_default;

-- Index: ix_user_perm_user_id_perm_id

-- DROP INDEX public.ix_user_perm_user_id_perm_id;

CREATE UNIQUE INDEX ix_user_perm_user_id_perm_id
    ON public.user_perm USING btree
    (user_id, perm_id)
    TABLESPACE pg_default;
