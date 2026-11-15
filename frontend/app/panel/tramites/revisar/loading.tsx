import ListadoEsqueleto from "@/components/common/ListadoEsqueleto";

// Next exige este archivo por convención para la carga del segmento. En vez del Loader genérico,
// muestra la silueta de la tabla que está a punto de aparecer.
export default function Cargando() {
  return <ListadoEsqueleto titulo="Mi bandeja" columnas={6} conFiltros={false} />;
}
