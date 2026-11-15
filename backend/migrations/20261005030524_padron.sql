-- Create "facultad" table
CREATE TABLE "facultad" (
  "id" smallserial NOT NULL,
  "nombre" character varying(160) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_facultad_nombre" UNIQUE ("nombre")
);
-- Modify "docente" table
ALTER TABLE "docente" ALTER COLUMN "cedula" DROP NOT NULL, ADD COLUMN "codigo" character varying(20) NULL, ADD COLUMN "facultad_id" smallint NULL, ADD CONSTRAINT "uni_docente_codigo" UNIQUE ("codigo"), ADD CONSTRAINT "fk_docente_facultad" FOREIGN KEY ("facultad_id") REFERENCES "facultad" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT;
-- Create "legajo" table
CREATE TABLE "legajo" (
  "id" serial NOT NULL,
  "docente_id" uuid NOT NULL,
  "gaveta" character varying(10) NULL,
  "orden" integer NULL,
  "anexos" character varying(80) NULL,
  "observacion" character varying(200) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_legajo_docente" FOREIGN KEY ("docente_id") REFERENCES "docente" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_legajo_docente" to table: "legajo"
CREATE INDEX "idx_legajo_docente" ON "legajo" ("docente_id");
-- Create index "idx_legajo_gaveta" to table: "legajo"
CREATE INDEX "idx_legajo_gaveta" ON "legajo" ("gaveta");
