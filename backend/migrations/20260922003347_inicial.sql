-- Create "rol" table
CREATE TABLE "rol" (
  "id" smallserial NOT NULL,
  "nombre" character varying(50) NOT NULL,
  "descripcion" character varying(200) NULL,
  "es_sistema" boolean NOT NULL DEFAULT false,
  "activo" boolean NOT NULL DEFAULT true,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_rol_nombre" UNIQUE ("nombre")
);
-- Create "usuario" table
CREATE TABLE "usuario" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "rol_id" smallint NOT NULL,
  "nombres" character varying(80) NOT NULL,
  "apellidos" character varying(80) NOT NULL,
  "correo" character varying(120) NOT NULL,
  "telefono" character varying(20) NULL,
  "hash_contrasena" character varying(100) NOT NULL,
  "intentos_fallidos" smallint NOT NULL DEFAULT 0,
  "bloqueado_hasta" timestamptz NULL,
  "version_sesion" integer NOT NULL DEFAULT 1,
  "activo" boolean NOT NULL DEFAULT true,
  "ultimo_acceso" timestamptz NULL,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  "actualizado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_usuario_rol" FOREIGN KEY ("rol_id") REFERENCES "rol" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "chk_usuario_intentos" CHECK (intentos_fallidos >= 0)
);
-- Create index "idx_usuario_rol" to table: "usuario"
CREATE INDEX "idx_usuario_rol" ON "usuario" ("rol_id");
-- Create "auditoria" table
CREATE TABLE "auditoria" (
  "id" bigserial NOT NULL,
  "usuario_id" uuid NULL,
  "accion" character varying(40) NOT NULL,
  "entidad" character varying(40) NOT NULL,
  "entidad_id" character varying(40) NULL,
  "detalle" jsonb NULL,
  "direccion_ip" inet NULL,
  "fecha" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_auditoria_usuario" FOREIGN KEY ("usuario_id") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- Create index "idx_auditoria_entidad" to table: "auditoria"
CREATE INDEX "idx_auditoria_entidad" ON "auditoria" ("entidad", "entidad_id");
-- Create index "idx_auditoria_usuario_fecha" to table: "auditoria"
CREATE INDEX "idx_auditoria_usuario_fecha" ON "auditoria" ("usuario_id", "fecha" DESC);
-- Create "permiso" table
CREATE TABLE "permiso" (
  "id" smallserial NOT NULL,
  "codigo" character varying(60) NOT NULL,
  "nombre" character varying(120) NOT NULL,
  "modulo" character varying(40) NOT NULL,
  "caso_uso" character varying(6) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_permiso_codigo" UNIQUE ("codigo")
);
-- Create "rol_permiso" table
CREATE TABLE "rol_permiso" (
  "rol_id" smallint NOT NULL,
  "permiso_id" smallint NOT NULL,
  PRIMARY KEY ("rol_id", "permiso_id"),
  CONSTRAINT "fk_rol_permiso_permiso" FOREIGN KEY ("permiso_id") REFERENCES "permiso" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_rol_permiso_rol" FOREIGN KEY ("rol_id") REFERENCES "rol" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "token_recuperacion" table
CREATE TABLE "token_recuperacion" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "usuario_id" uuid NOT NULL,
  "proposito" character varying(12) NOT NULL,
  "token_hash" character(64) NOT NULL,
  "expira_en" timestamptz NOT NULL,
  "usado_en" timestamptz NULL,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_token_recuperacion_token_hash" UNIQUE ("token_hash"),
  CONSTRAINT "fk_token_recuperacion_usuario" FOREIGN KEY ("usuario_id") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_token_proposito" CHECK ((proposito)::text = ANY ((ARRAY['RECUPERACION'::character varying, 'ACTIVACION'::character varying])::text[]))
);
