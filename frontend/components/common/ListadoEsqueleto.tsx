import EncabezadoPagina from "./EncabezadoPagina";
import TablaEsqueleto from "./TablaEsqueleto";

// Pantalla de carga de un listado del panel. Replica la estructura de las páginas de listado
// (encabezado, barra de filtros, tabla) con la misma separación, para que cuando lleguen los
// datos solo cambie el contenido y no la maquetación.
//
// El título va escrito, no esqueletizado: es constante en cada ruta, así que no hay motivo para
// ocultárselo a alguien que acaba de pulsar ese mismo enlace del menú. Lo que no se sabe todavía
// —el recuento, los filtros, las filas— es lo único que aparece como silueta.
export default function ListadoEsqueleto({
  titulo,
  columnas = 4,
  filas = 8,
  conFiltros = true,
  conAccion = false,
}: {
  titulo: string;
  columnas?: number;
  filas?: number;
  conFiltros?: boolean;
  conAccion?: boolean;
}) {
  return (
    <div className="flex flex-col gap-6" role="status" aria-busy="true" aria-label={`Cargando ${titulo}`}>
      <EncabezadoPagina
        titulo={titulo}
        descripcion={<span className="esqueleto mt-1 block h-3 w-52 rounded-full" />}
        accion={conAccion ? <span className="esqueleto block h-9 w-36 rounded-lg" /> : undefined}
      />

      {conFiltros && (
        <div className="flex gap-2.5" aria-hidden>
          <span className="esqueleto block h-[42px] w-full max-w-[320px] rounded-lg" />
          <span className="esqueleto block h-[42px] w-[104px] rounded-lg" />
        </div>
      )}

      <TablaEsqueleto columnas={columnas} filas={filas} />
    </div>
  );
}
