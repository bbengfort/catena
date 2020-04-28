-- NOTE: this schema file is always run for any migration command because the migration
-- commands require this table. This is not true for other migrations, which are only
-- run on demand and whose state are stored in this table.
-- migrate: up

CREATE TABLE IF NOT EXISTS migrations (
    "revision" integer NOT NULL,
    "name" varchar(128) NOT NULL,
    "active" boolean NOT NULL DEFAULT false,
    "applied" TIMESTAMP WITH TIME ZONE,
    "created" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("revision")
) WITHOUT OIDS;

COMMENT ON TABLE "migrations" IS 'Manages the state of database by enabling migrations and rollbacks';
COMMENT ON COLUMN "migrations"."revision" IS 'The revision id parsed from the filename of the migration';
COMMENT ON COLUMN "migrations"."name" IS 'The name of the migration parsed from the filename of the migration';
COMMENT ON COLUMN "migrations"."active" IS 'If the migration has been applied, set to false on rollbacks or if not applied';
COMMENT ON COLUMN "migrations"."applied" IS 'Timestamp when the migration was applied, null if rolledback or not applied';
COMMENT ON COLUMN "migrations"."created" IS 'Timestamp when the migration was created';

-- NOTE: the down migration is run to complete reset the state of migrations if
-- something has gone completely sideways.
-- migrate: down

DROP TABLE IF EXISTS migrations CASCADE;