-- Flujo real del Departamento: el trámite sale a la facultad (tramo sin seguimiento), regresa con la
-- resolución, se revisa de nuevo, se remite al Rectorado y Archivo lo archiva. Las carpetas solicitadas
-- se ligan al trámite. Los estados anteriores siguen siendo válidos: no hay datos que migrar.
ALTER TABLE "tramite" DROP CONSTRAINT "chk_tramite_estado";
ALTER TABLE "tramite"
  ALTER COLUMN "estado" TYPE character varying(20),
  ADD COLUMN "anio_desde" smallint NULL,
  ADD COLUMN "anio_hasta" smallint NULL,
  ADD COLUMN "resolucion_numero" character varying(40) NULL,
  ADD COLUMN "resolucion_fecha" date NULL,
  ADD COLUMN "remision_nota" character varying(120) NULL,
  ADD CONSTRAINT "chk_tramite_estado" CHECK (estado IN ('INGRESADO','ASIGNADO','EN_REVISION','OBSERVADO','DICTAMINADO','EN_FACULTAD','REVISION_FINAL','REMITIDO','ARCHIVADO'));

ALTER TABLE "tramite_historial"
  ALTER COLUMN "estado" TYPE character varying(20),
  ALTER COLUMN "estado_anterior" TYPE character varying(20);

ALTER TABLE "notificacion_envio" ALTER COLUMN "estado_tramite" TYPE character varying(20);

ALTER TABLE "solicitud_carpeta"
  ADD COLUMN "tramite_id" uuid NULL,
  ADD CONSTRAINT "fk_solicitud_carpeta_tramite" FOREIGN KEY ("tramite_id") REFERENCES "tramite" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
CREATE INDEX "idx_solicitud_tramite" ON "solicitud_carpeta" ("tramite_id");
