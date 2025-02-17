DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'social') THEN
      PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE social');
   END IF;
END
$$;

\c social

CREATE EXTENSION IF NOT EXISTS citext;