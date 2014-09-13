/* mysql
  CREATE TABLE `timer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `place_id` varchar(100) NOT NULL,
  `timer_ms` int(11) NOT NULL,
  `ip` int(10) unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)
*/

--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: timer; Type: TABLE; Schema: public; Owner: sethammons; Tablespace:
--

CREATE TABLE timer (
    id integer NOT NULL,
    place_id character varying(100) NOT NULL,
    time_ms integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    ip inet
);


ALTER TABLE public.timer OWNER TO sethammons;

--
-- Name: timer_id_seq; Type: SEQUENCE; Schema: public; Owner: sethammons
--

CREATE SEQUENCE timer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.timer_id_seq OWNER TO sethammons;

--
-- Name: timer_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sethammons
--

ALTER SEQUENCE timer_id_seq OWNED BY timer.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: sethammons
--

ALTER TABLE ONLY timer ALTER COLUMN id SET DEFAULT nextval('timer_id_seq'::regclass);


--
-- Name: timer_pkey; Type: CONSTRAINT; Schema: public; Owner: sethammons; Tablespace:
--

ALTER TABLE ONLY timer
    ADD CONSTRAINT timer_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--