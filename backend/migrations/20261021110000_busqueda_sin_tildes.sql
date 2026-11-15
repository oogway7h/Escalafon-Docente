-- La búsqueda dentro de un documento debe ignorar tildes ("resolucion" encuentra "RESOLUCIÓN"):
-- configuración de texto español + unaccent, y el índice por página pasa a usarla.
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE TEXT SEARCH CONFIGURATION es_unaccent (COPY = spanish);
ALTER TEXT SEARCH CONFIGURATION es_unaccent
  ALTER MAPPING FOR hword, hword_part, word WITH unaccent, spanish_stem;
DROP INDEX IF EXISTS ix_documento_pagina_fts;
CREATE INDEX ix_documento_pagina_fts ON "documento_pagina" USING gin (to_tsvector('es_unaccent', texto));
