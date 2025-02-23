--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6
-- Dumped by pg_dump version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE IF EXISTS test;
--
-- Name: test; Type: DATABASE; Schema: -; Owner: test
--

CREATE DATABASE test WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C' LC_CTYPE = 'ru_RU.UTF-8';


ALTER DATABASE test OWNER TO test;

\connect test

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: objects; Type: TABLE; Schema: public; Owner: test
--

DROP TABLE IF EXISTS public.objects;

CREATE TABLE public.objects (
    id uuid NOT NULL,
    o_string_1 text,
    o_bool boolean,
    o_float_64 double precision,
    o_uuid_2 uuid,
    o_time_1 timestamp with time zone,
    o_float_32 real DEFAULT 0 NOT NULL,
    o_int integer,
    o_null jsonb generated always as (null) stored,
    o_int_16 smallint DEFAULT 16 NOT NULL,
    o_time_2 timestamp with time zone,
    o_time_3 timestamp with time zone,
    o_time_4 timestamp with time zone,
    o_time_0 timestamp with time zone
);


ALTER TABLE public.objects OWNER TO test;

--
-- Data for Name: objects; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.objects (id, o_string_1, o_bool, o_float_64, o_uuid_2, o_time_1, o_float_32, o_int, o_null, o_int_16, o_time_2, o_time_3, o_time_4, o_time_0) FROM stdin;
\.


--
-- Name: objects objects_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.objects
    ADD CONSTRAINT objects_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

