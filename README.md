# envir_collector_go
Data collector for Envir energy monitor ([Envi kit](http://www.currentcost.net/Monitor%20Details.html). This replaces the [python version](https://github.com/clinstid/energydash) I wrote a while back.

# Database schema
Unlike the original version from my [energydash](https://github.com/clinstid/energydash#data-collection) collector, I'm now using postgres as my database because the collected data fits very nicely into a traditional RDBMS. Mongo was fun to play with, but it definitely was not the right choice for this kind of data.

The data is currently stored in a single table called `energydash` (yeah, I kept the name):

```sql
-- Table: public.energydash

-- DROP TABLE public.energydash;

CREATE TABLE public.energydash
(
  src character varying,
  dsb integer,
  "time" timestamp with time zone NOT NULL,
  tmprf double precision,
  sensor integer,
  device_id character varying,
  ch1_watts integer,
  ch2_watts integer,
  ch3_watts integer,
  CONSTRAINT energydash_pkey PRIMARY KEY ("time")
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.energydash
  OWNER TO postgres;
GRANT ALL ON TABLE public.energydash TO postgres;
GRANT ALL ON TABLE public.energydash TO energydash;
```
