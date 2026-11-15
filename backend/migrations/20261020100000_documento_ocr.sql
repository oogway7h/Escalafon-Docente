-- Carga de documentos con OCR y compresión (CU-15). El archivo se sube directo a S3 y el ia-service
-- (Tesseract + ocrmypdf) extrae el texto y genera una versión comprimida; el original se conserva.
-- hash_sha256 pasa a ser opcional porque se calcula después de la subida.
ALTER TABLE "tramite_documento"
  ALTER COLUMN "hash_sha256" DROP NOT NULL,
  ADD COLUMN "estado_ocr" character varying(12) NOT NULL DEFAULT 'NO_APLICA',
  ADD COLUMN "texto_ocr" text NULL,
  ADD COLUMN "texto_validado" boolean NOT NULL DEFAULT false,
  ADD COLUMN "error_ocr" character varying(300) NULL,
  ADD COLUMN "clave_comprimida" character varying(300) NULL,
  ADD COLUMN "tamano_comprimido" bigint NULL,
  ADD CONSTRAINT "ck_tramite_documento_estado_ocr" CHECK (estado_ocr IN ('NO_APLICA','SUBIENDO','PENDIENTE','PROCESANDO','PROCESADO','ERROR'));

-- Búsqueda de texto completo sobre el texto extraído (la usará CU-18). Manual: Atlas no la conoce.
CREATE INDEX ix_tramite_documento_texto_fts ON "tramite_documento" USING gin (to_tsvector('spanish', coalesce(texto_ocr, '')));
