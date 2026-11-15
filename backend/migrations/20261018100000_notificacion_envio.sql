-- Create "notificacion_envio" table (CU-14): avisos de cambio de estado al docente
CREATE TABLE "notificacion_envio" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tramite_id" uuid NOT NULL,
  "tramite_codigo" character varying(30) NOT NULL,
  "estado_tramite" character varying(12) NOT NULL,
  "destino" character varying(120) NOT NULL DEFAULT '',
  "nombres" character varying(80) NOT NULL DEFAULT '',
  "estado" character varying(9) NOT NULL DEFAULT 'PENDIENTE',
  "intentos" smallint NOT NULL DEFAULT 0,
  "error" text NULL,
  "creado_en" timestamptz NOT NULL DEFAULT now(),
  "enviado_en" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notificacion_envio_tramite" FOREIGN KEY ("tramite_id") REFERENCES "tramite" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_notificacion_estado" ON "notificacion_envio" ("estado", "creado_en");
CREATE INDEX "idx_notificacion_tramite" ON "notificacion_envio" ("tramite_id");
