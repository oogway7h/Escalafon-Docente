-- Índice para el orden por fecha de la consulta de auditoría (CU-06)
CREATE INDEX "idx_auditoria_fecha" ON "auditoria" ("fecha" DESC);

-- Legajo: código único y estado (CU-19)
ALTER TABLE "legajo" ADD COLUMN "codigo" character varying(20) NULL;
ALTER TABLE "legajo" ADD COLUMN "estado" character varying(10) NOT NULL DEFAULT 'EN_ARCHIVO';
UPDATE "legajo" SET "codigo" = 'LEG-' || lpad("id"::text, 5, '0');
ALTER TABLE "legajo" ALTER COLUMN "codigo" SET NOT NULL;
CREATE UNIQUE INDEX "ux_legajo_codigo" ON "legajo" ("codigo");
ALTER TABLE "legajo" ADD CONSTRAINT "chk_legajo_estado" CHECK ("estado" IN ('EN_ARCHIVO', 'PRESTADO'));

-- Los legajos nuevos (p. ej. los del importador) reciben código si no traen uno
CREATE OR REPLACE FUNCTION fn_legajo_codigo() RETURNS trigger AS $$
BEGIN
  IF NEW.codigo IS NULL OR NEW.codigo = '' THEN
    NEW.codigo := 'LEG-' || lpad(NEW.id::text, 5, '0');
  END IF;
  RETURN NEW;
END; $$ LANGUAGE plpgsql;
CREATE TRIGGER tg_legajo_codigo BEFORE INSERT ON "legajo"
  FOR EACH ROW EXECUTE FUNCTION fn_legajo_codigo();

-- Create "solicitud_carpeta" table (CU-20, CU-21)
CREATE TABLE "solicitud_carpeta" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "legajo_id" integer NOT NULL,
  "solicitante_id" uuid NOT NULL,
  "motivo" text NOT NULL,
  "estado" character varying(10) NOT NULL DEFAULT 'PENDIENTE',
  "respuesta" text NULL,
  "atendido_por" uuid NULL,
  "atendido_en" timestamptz NULL,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_solicitud_carpeta_legajo" FOREIGN KEY ("legajo_id") REFERENCES "legajo" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_solicitud_carpeta_solicitante" FOREIGN KEY ("solicitante_id") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
CREATE INDEX "idx_solicitud_legajo" ON "solicitud_carpeta" ("legajo_id");
CREATE INDEX "idx_solicitud_estado" ON "solicitud_carpeta" ("estado");
-- A lo sumo una solicitud en curso por legajo
CREATE UNIQUE INDEX "ux_solicitud_abierta" ON "solicitud_carpeta" ("legajo_id")
  WHERE "estado" IN ('PENDIENTE', 'ATENDIDA', 'DESPACHADA');

-- Create "movimiento_custodia" table (CU-22 al CU-24), de solo anexión
CREATE TABLE "movimiento_custodia" (
  "id" bigserial NOT NULL,
  "legajo_id" integer NOT NULL,
  "solicitud_id" uuid NOT NULL,
  "tipo" character varying(7) NOT NULL,
  "custodio_id" uuid NOT NULL,
  "fecha_limite" date NULL,
  "estado_material" character varying(11) NOT NULL,
  "observacion" text NULL,
  "registrado_por" uuid NOT NULL,
  "fecha" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_movimiento_custodia_legajo" FOREIGN KEY ("legajo_id") REFERENCES "legajo" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_movimiento_custodia_solicitud" FOREIGN KEY ("solicitud_id") REFERENCES "solicitud_carpeta" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_movimiento_custodia_custodio" FOREIGN KEY ("custodio_id") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_movimiento_custodia_registro" FOREIGN KEY ("registrado_por") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "chk_movimiento_tipo" CHECK ("tipo" IN ('SALIDA', 'ENTRADA'))
);
CREATE INDEX "idx_movimiento_legajo" ON "movimiento_custodia" ("legajo_id");
CREATE INDEX "idx_movimiento_fecha" ON "movimiento_custodia" ("fecha" DESC);
CREATE TRIGGER tg_movimiento_inmutable BEFORE UPDATE OR DELETE ON "movimiento_custodia"
  FOR EACH ROW EXECUTE FUNCTION fn_solo_anexion();
CREATE TRIGGER tg_movimiento_truncate BEFORE TRUNCATE ON "movimiento_custodia"
  FOR EACH STATEMENT EXECUTE FUNCTION fn_solo_anexion();
