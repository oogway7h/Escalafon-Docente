-- Baja lógica de los adjuntos del trámite (CU-16): activo=false retira el documento del expediente
-- y de las búsquedas sin borrar el objeto del almacén; el motivo y la fecha quedan para la
-- restitución del Director (CU-17) y para la bitácora de auditoría.
ALTER TABLE "tramite_documento"
  ADD COLUMN "activo" boolean NOT NULL DEFAULT true,
  ADD COLUMN "motivo_baja" character varying(255) NULL,
  ADD COLUMN "baja_en" timestamptz NULL,
  ADD COLUMN "baja_por" uuid NULL,
  ADD CONSTRAINT "fk_tramite_documento_baja" FOREIGN KEY ("baja_por") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT;
