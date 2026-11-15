-- Create "categoria" table
CREATE TABLE "categoria" (
  "id" smallserial NOT NULL,
  "nombre" character varying(60) NOT NULL,
  "anios_antiguedad" smallint NOT NULL DEFAULT 0,
  "activa" boolean NOT NULL DEFAULT true,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_categoria_nombre" UNIQUE ("nombre"),
  CONSTRAINT "chk_categoria_anios" CHECK (anios_antiguedad >= 0)
);
-- Create "categoria_requisito" table
CREATE TABLE "categoria_requisito" (
  "id" serial NOT NULL,
  "categoria_id" smallint NOT NULL,
  "nombre" character varying(120) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categoria_requisito_categoria" FOREIGN KEY ("categoria_id") REFERENCES "categoria" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "ux_categoria_requisito" to table: "categoria_requisito"
CREATE UNIQUE INDEX "ux_categoria_requisito" ON "categoria_requisito" ("categoria_id", "nombre");
-- Create "docente" table
CREATE TABLE "docente" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "cedula" character varying(20) NOT NULL,
  "nombres" character varying(80) NOT NULL,
  "apellidos" character varying(80) NOT NULL,
  "correo" character varying(120) NULL,
  "telefono" character varying(20) NULL,
  "categoria_id" smallint NULL,
  "activo" boolean NOT NULL DEFAULT true,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  "actualizado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_docente_cedula" UNIQUE ("cedula"),
  CONSTRAINT "fk_docente_categoria" FOREIGN KEY ("categoria_id") REFERENCES "categoria" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- Create "tramite" table
CREATE TABLE "tramite" (
  "id" uuid NOT NULL,
  "codigo" character varying(16) NOT NULL,
  "docente_id" uuid NOT NULL,
  "categoria_id" smallint NOT NULL,
  "tipo" character varying(14) NOT NULL,
  "estado" character varying(12) NOT NULL,
  "origen" character varying(6) NOT NULL,
  "pendiente_validacion" boolean NOT NULL DEFAULT false,
  "analista_id" uuid NULL,
  "radicado_por" uuid NULL,
  "resultado" character varying(12) NULL,
  "dictamen" text NULL,
  "dictaminado_en" timestamptz NULL,
  "archivado_en" timestamptz NULL,
  "ubicacion_legajo" character varying(160) NULL,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  "actualizado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_tramite_codigo" UNIQUE ("codigo"),
  CONSTRAINT "fk_tramite_analista" FOREIGN KEY ("analista_id") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_tramite_categoria" FOREIGN KEY ("categoria_id") REFERENCES "categoria" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_tramite_docente" FOREIGN KEY ("docente_id") REFERENCES "docente" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "chk_tramite_estado" CHECK ((estado)::text = ANY ((ARRAY['INGRESADO'::character varying, 'ASIGNADO'::character varying, 'EN_REVISION'::character varying, 'OBSERVADO'::character varying, 'DICTAMINADO'::character varying, 'ARCHIVADO'::character varying])::text[])),
  CONSTRAINT "chk_tramite_origen" CHECK ((origen)::text = ANY ((ARRAY['STAFF'::character varying, 'PORTAL'::character varying])::text[])),
  CONSTRAINT "chk_tramite_resultado" CHECK ((resultado)::text = ANY ((ARRAY['FAVORABLE'::character varying, 'DESFAVORABLE'::character varying])::text[])),
  CONSTRAINT "chk_tramite_tipo" CHECK ((tipo)::text = ANY ((ARRAY['CATEGORIZACION'::character varying, 'REGULARIZACION'::character varying])::text[]))
);
-- Create index "idx_tramite_analista" to table: "tramite"
CREATE INDEX "idx_tramite_analista" ON "tramite" ("analista_id");
-- Create index "idx_tramite_docente" to table: "tramite"
CREATE INDEX "idx_tramite_docente" ON "tramite" ("docente_id");
-- Create index "idx_tramite_estado" to table: "tramite"
CREATE INDEX "idx_tramite_estado" ON "tramite" ("estado");
-- Create "observacion" table
CREATE TABLE "observacion" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tramite_id" uuid NOT NULL,
  "documento" character varying(120) NOT NULL,
  "motivo" text NOT NULL,
  "plazo" date NOT NULL,
  "resuelta" boolean NOT NULL DEFAULT false,
  "resuelta_en" timestamptz NULL,
  "creada_por" uuid NOT NULL,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_observacion_autor" FOREIGN KEY ("creada_por") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_observacion_tramite" FOREIGN KEY ("tramite_id") REFERENCES "tramite" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_observacion_tramite" to table: "observacion"
CREATE INDEX "idx_observacion_tramite" ON "observacion" ("tramite_id");
-- Create "tramite_documento" table
CREATE TABLE "tramite_documento" (
  "id" uuid NOT NULL,
  "tramite_id" uuid NOT NULL,
  "requisito_id" integer NULL,
  "nombre_original" character varying(200) NOT NULL,
  "tipo_mime" character varying(80) NOT NULL,
  "tamano" bigint NOT NULL,
  "hash_sha256" character(64) NOT NULL,
  "bucket" character varying(100) NOT NULL,
  "clave_objeto" character varying(300) NOT NULL,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tramite_documento_requisito" FOREIGN KEY ("requisito_id") REFERENCES "categoria_requisito" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_tramite_documento_tramite" FOREIGN KEY ("tramite_id") REFERENCES "tramite" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_tramite_documento_tramite" to table: "tramite_documento"
CREATE INDEX "idx_tramite_documento_tramite" ON "tramite_documento" ("tramite_id");
-- Create "tramite_historial" table
CREATE TABLE "tramite_historial" (
  "id" bigserial NOT NULL,
  "tramite_id" uuid NOT NULL,
  "estado" character varying(12) NOT NULL,
  "estado_anterior" character varying(12) NULL,
  "usuario_id" uuid NULL,
  "detalle" text NULL,
  "fecha" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tramite_historial_tramite" FOREIGN KEY ("tramite_id") REFERENCES "tramite" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_tramite_historial_usuario" FOREIGN KEY ("usuario_id") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- Create index "idx_historial_tramite" to table: "tramite_historial"
CREATE INDEX "idx_historial_tramite" ON "tramite_historial" ("tramite_id");
-- Historial de estados de solo anexión (reutiliza fn_solo_anexion de ajustes_manuales)
CREATE TRIGGER tg_tramite_historial_inmutable BEFORE UPDATE OR DELETE ON tramite_historial
  FOR EACH ROW EXECUTE FUNCTION fn_solo_anexion();
CREATE TRIGGER tg_tramite_historial_truncate BEFORE TRUNCATE ON tramite_historial
  FOR EACH STATEMENT EXECUTE FUNCTION fn_solo_anexion();

-- actualizado_en automático
CREATE TRIGGER tg_docente_actualizado BEFORE UPDATE ON docente
  FOR EACH ROW EXECUTE FUNCTION fn_actualizado_en();
CREATE TRIGGER tg_tramite_actualizado BEFORE UPDATE ON tramite
  FOR EACH ROW EXECUTE FUNCTION fn_actualizado_en();
