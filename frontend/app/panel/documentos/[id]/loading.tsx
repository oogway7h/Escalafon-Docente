// Frontera de carga propia del segmento. Antes la ponia app/panel/loading.tsx, pero al ser
// ancestro de los listados hacia que en estos se viera primero este Loader y despues el
// esqueleto: dos fallbacks anidados para una sola navegacion.
export { default } from "@/components/common/Loader";
