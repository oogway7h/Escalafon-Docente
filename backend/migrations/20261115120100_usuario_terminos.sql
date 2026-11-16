-- Modify "usuario" table: acuerdo de usuario aceptado (versión y fecha)
ALTER TABLE "usuario" ADD COLUMN "terminos_version" smallint NOT NULL DEFAULT 0, ADD COLUMN "terminos_aceptados_en" timestamptz NULL;
