import { claseBoton } from "@/components/common/Boton";
import { BuscadorCarga, LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Insignia from "@/components/common/Insignia";
import Paginacion from "@/components/common/Paginacion";
import Tabla from "@/components/common/Tabla";
import { IconoEditar } from "@/components/layout/IconosPanel";
import { exigirPermiso } from "@/lib/api";
import { listarCarpetas } from "@/lib/api-legajos";

export default async function LegajosPage({ searchParams }: { searchParams: Promise<{ buscar?: string; pagina?: string }> }) {
  await exigirPermiso("LEGAJO_GESTIONAR");
  const sp = await searchParams;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const lista = await listarCarpetas({ buscar: sp.buscar }, pagina);
  const totalPaginas = Math.max(1, Math.ceil(lista.total / lista.tamano));

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Inventario de legajos"
          descripcion={<>{lista.total} carpetas{sp.buscar ? " coinciden con la búsqueda" : " en el inventario"}.</>}
          accion={
            <LinkCarga href="/panel/legajos/nuevo" className={claseBoton("primario", "sm")}>
              + Nuevo legajo
            </LinkCarga>
          }
        />

        <BuscadorCarga action="/panel/legajos" defaultValue={sp.buscar} placeholder="Código de legajo o de docente, cédula o nombre" />

        <ZonaCarga>
          <Tabla
            filas={lista.datos}
            clave={(l) => l.id}
            vacio="No se encontraron legajos."
            columnas={[
              { titulo: "Legajo", className: "w-28", celda: (l) => <span className="font-mono text-sm font-semibold">{l.codigo}</span> },
              {
                titulo: "Docente",
                celda: (l) => (
                  <>
                    <div className="font-semibold">{l.docenteNombre}</div>
                    <div className="text-xs text-ink-faint">Código {l.docenteCodigo || "—"}</div>
                  </>
                ),
              },
              {
                titulo: "Ubicación",
                celda: (l) => (
                  <>
                    <div>
                      Gaveta {l.gaveta || "s/g"} · N.º {l.orden}
                    </div>
                    {l.anexos && <div className="text-xs text-ink-faint">{l.anexos}</div>}
                  </>
                ),
              },
              {
                titulo: "Estado",
                celda: (l) => <Insignia tono={l.estado === "EN_ARCHIVO" ? "exito" : "acento"}>{l.estado === "EN_ARCHIVO" ? "En archivo" : "Prestado"}</Insignia>,
              },
              {
                className: "w-16 text-right",
                celda: (l) => (
                  <LinkCarga
                    href={`/panel/legajos/${l.id}/editar`}
                    className="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-line text-ink-soft hover:border-accent hover:text-accent-ink"
                    aria-label={`Editar legajo ${l.codigo}`}
                  >
                    <IconoEditar className="h-3.5 w-3.5" />
                  </LinkCarga>
                ),
              },
            ]}
          />
        </ZonaCarga>

        <Paginacion pagina={lista.pagina} totalPaginas={totalPaginas} searchParams={{ buscar: sp.buscar }} />
      </div>
    </ListadoConCarga>
  );
}
