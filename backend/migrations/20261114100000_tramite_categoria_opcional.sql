-- Las regularizaciones no postulan a ninguna categoría: categoria_id pasa a ser opcional, pero sigue siendo
-- obligatoria en las categorizaciones.
ALTER TABLE "tramite" ALTER COLUMN "categoria_id" DROP NOT NULL, ADD CONSTRAINT "chk_tramite_categoria" CHECK (((tipo)::text = 'REGULARIZACION'::text) OR (categoria_id IS NOT NULL));
