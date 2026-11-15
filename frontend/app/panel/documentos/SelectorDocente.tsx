import { BuscadorCarga, LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import { claseBoton } from "@/components/common/Boton";
import Paginacion from "@/components/common/Paginacion";
import Tabla from "@/components/common/Tabla";
import { listarDocentesPagina } from "@/lib/api-tramites";

// Elegir al docente cuyo expediente se sube o se consulta: se busca por código, cédula o nombre.
export default async function SelectorDocente({
  ruta,
  buscar,
  pagina,
  etiquetaAccion,
  destino,
  sinBuscador = false,
}: {
  ruta: string; // página que recibe ?buscar= y ?pagina=
  buscar?: string;
  pagina: number;
  etiquetaAccion: string;
  destino: (docenteId: string) => string;
  sinBuscador?: boolean; // la página ya tiene su propio buscador
}) {
  const lista = await listarDocentesPagina(buscar, pagina);
  const totalPaginas = Math.max(1, Math.ceil(lista.total / lista.tamano));
  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-5">
        {!sinBuscador && <BuscadorCarga action={ruta} defaultValue={buscar} placeholder="Buscar docente por código, cédula o nombre" />}
        <ZonaCarga>
          <Tabla
            filas={lista.datos}
            clave={(d) => d.id}
            vacio="No se encontraron docentes."
            columnas={[
              { titulo: "Código", className: "w-24", celda: (d) => <span className="font-mono text-sm font-semibold">{d.codigo || "—"}</span> },
              {
                titulo: "Docente",
                celda: (d) => (
                  <>
                    <div className="font-semibold">
                      {d.apellidos}, {d.nombres}
                    </div>
                    {d.cedula && <div className="text-xs text-ink-faint">CI {d.cedula}</div>}
                  </>
                ),
              },
              { titulo: "Categoría", celda: (d) => d.categoria || <span className="text-ink-faint">Sin categoría</span> },
              {
                className: "w-56 text-right",
                celda: (d) => (
                  <LinkCarga href={destino(d.id)} className={claseBoton("secundario", "sm")}>
                    {etiquetaAccion}
                  </LinkCarga>
                ),
              },
            ]}
          />
        </ZonaCarga>
        <Paginacion pagina={lista.pagina} totalPaginas={totalPaginas} searchParams={{ buscar }} />
      </div>
    </ListadoConCarga>
  );
}
