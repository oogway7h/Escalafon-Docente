-- Texto OCR por página (CU-15): permite buscar dentro de un documento y saltar a la página donde aparece.
-- texto_ocr se conserva como la concatenación de las páginas. El índice de texto completo por página es
-- manual (Atlas no lo conoce) y servirá también a la búsqueda documental (CU-18).
CREATE TABLE "documento_pagina" (
  "documento_id" uuid NOT NULL,
  "pagina" integer NOT NULL,
  "texto" text NOT NULL DEFAULT '',
  PRIMARY KEY ("documento_id", "pagina"),
  CONSTRAINT "fk_documento_pagina_documento" FOREIGN KEY ("documento_id") REFERENCES "tramite_documento" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX ix_documento_pagina_fts ON "documento_pagina" USING gin (to_tsvector('spanish', texto));

-- Lo ya procesado se parte por el salto de página (\f) que dejó ocrmypdf; sin salto, es una sola página.
INSERT INTO "documento_pagina" ("documento_id", "pagina", "texto")
SELECT d.id, p.n, trim(p.t)
FROM "tramite_documento" d
CROSS JOIN LATERAL unnest(string_to_array(d.texto_ocr, E'\f')) WITH ORDINALITY AS p(t, n)
WHERE d.estado_ocr = 'PROCESADO' AND d.texto_ocr IS NOT NULL
ON CONFLICT DO NOTHING;
