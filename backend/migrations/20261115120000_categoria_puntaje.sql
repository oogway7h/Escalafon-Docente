-- Modify "categoria" table: rango de puntos de categoría (Reglamento del Escalafón, art. 18)
ALTER TABLE "categoria" ADD COLUMN "puntaje_desde" smallint NULL, ADD COLUMN "puntaje_hasta" smallint NULL;
