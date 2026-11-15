-- Expediente digital del docente: los escaneos que sube el Personal de Archivo (Currículum, Regularización,
-- Adjunto, Asociado, Titular A/B/C, Emérito) cuelgan del docente y no de un trámite, así que tramite_id pasa a
-- ser opcional. Cada documento tiene exactamente un dueño: un trámite (adjunto) o un docente (escaneo).
ALTER TABLE "tramite_documento" ALTER COLUMN "tramite_id" DROP NOT NULL;
ALTER TABLE "tramite_documento"
  ADD COLUMN "docente_id" uuid NULL,
  ADD COLUMN "tipo_documento" character varying(20) NULL;
ALTER TABLE "tramite_documento"
  ADD CONSTRAINT "fk_tramite_documento_docente" FOREIGN KEY ("docente_id") REFERENCES "docente" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
CREATE INDEX "idx_tramite_documento_docente" ON "tramite_documento" ("docente_id");
ALTER TABLE "tramite_documento"
  ADD CONSTRAINT "ck_documento_duenio" CHECK ((tramite_id IS NOT NULL) <> (docente_id IS NOT NULL)),
  ADD CONSTRAINT "ck_documento_tipo" CHECK (tipo_documento IS NULL OR tipo_documento IN
    ('CURRICULUM', 'REGULARIZACION', 'ADJUNTO', 'ASOCIADO', 'TITULAR_A', 'TITULAR_B', 'TITULAR_C', 'EMERITO'));

-- La búsqueda usa documento_pagina (es_unaccent); este índice sobre texto_ocr ya no lo consulta nadie.
DROP INDEX IF EXISTS ix_tramite_documento_texto_fts;
