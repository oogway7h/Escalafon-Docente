-- Se retira la foto de evidencia de la devolución de carpetas (CU-23): queda el estado material y la observación.
ALTER TABLE "movimiento_custodia" DROP COLUMN "evidencia_clave";
