import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import { BuscadorCarga, FormularioCarga, LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Paginacion from "@/components/common/Paginacion";
import Insignia from "@/components/common/Insignia";
import Tabla from "@/components/common/Tabla";
import { IconoEditar } from "@/components/layout/IconosPanel";
import { exigirPermiso } from "@/lib/api";
import { listarDocentesPagina } from "@/lib/api-tramites";
import { mensajeError } from "@/lib/formularios";
import { cambiarEstadoDocenteAction } from "./actions";

export default async function DocentesPage({
  searchParams,
}: {
  searchParams: Promise<{ buscar?: string; pagina?: string; error?: string }>;
}) {
  const sesion = await exigirPermiso("DOCENTE_GESTIONAR");
  // El analista consulta el padrón; solo el Director lo modifica (altas, ediciones y bajas).
  const puedeModificar = sesion.rol === "Director";
  const sp = await searchParams;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const lista = await listarDocentesPagina(sp.buscar, pagina);
  const totalPaginas = Math.max(1, Math.ceil(lista.total / lista.tamano));

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Docentes"
          descripcion={<>{lista.total} docentes{sp.buscar ? " coinciden con la búsqueda" : " en el padrón"}.</>}
          accion={
            puedeModificar ? (
              <LinkCarga href="/panel/docentes/nuevo" className={claseBoton("primario", "sm")}>
                + Nuevo docente
              </LinkCarga>
            ) : undefined
          }
        />

        {sp.error && <Alerta tipo="error">{mensajeError(sp.error, "No se pudo completar la acción.")}</Alerta>}

        {/* Al buscar o cambiar de página, la tabla actual queda difuminada con el Loader encima. */}
          <BuscadorCarga action="/panel/docentes" defaultValue={sp.buscar} placeholder="Buscar por código, cédula o nombre" />

          <ZonaCarga>
          <Tabla
            filas={lista.datos}
            clave={(d) => d.id}
            vacio="No se encontraron docentes."
            pie={puedeModificar ? "La baja es lógica: el docente se conserva en el padrón y en sus trámites, pero no se le pueden radicar nuevos." : undefined}
            columnas={[
              {
                titulo: "Código",
                className: "w-24",
                celda: (d) =>
                  d.codigo ? <span className="font-mono text-sm font-semibold">{d.codigo}</span> : <span className="text-xs text-accent-ink">Sin código</span>,
              },
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
              {
                titulo: "Categoría y facultad",
                celda: (d) => (
                  <>
                    <div>{d.categoria || <span className="text-ink-faint">Sin categoría</span>}</div>
                    <div className="max-w-[260px] truncate text-xs text-ink-faint" title={d.facultad}>
                      {d.facultad || "Sin facultad"}
                    </div>
                  </>
                ),
              },
              {
                titulo: "Legajo",
                celda: (d) =>
                  d.legajos.length === 0 ? (
                    <span className="text-ink-faint">—</span>
                  ) : (
                    <span className="text-sm">{d.legajos.map((l) => l.gaveta || "s/g").join(" · ")}</span>
                  ),
              },
              {
                titulo: "Estado",
                celda: (d) => <Insignia tono={d.activo ? "exito" : "neutro"}>{d.activo ? "Activo" : "De baja"}</Insignia>,
              },
              ...(!puedeModificar ? [] : [{
                className: "w-48 text-right",
                celda: (d: (typeof lista.datos)[number]) => (
                  <div className="flex items-center justify-end gap-2">
                    <LinkCarga
                      href={`/panel/docentes/${d.id}/editar`}
                      className="flex h-8 w-8 items-center justify-center rounded-lg border border-line text-ink-soft hover:border-accent hover:text-accent-ink"
                      aria-label={`Editar ${d.nombres} ${d.apellidos}`}
                    >
                      <IconoEditar className="h-3.5 w-3.5" />
                    </LinkCarga>
                    <FormularioCarga action={cambiarEstadoDocenteAction}>
                      <input type="hidden" name="id" value={d.id} />
                      <input type="hidden" name="activo" value={String(!d.activo)} />
                      <button type="submit" className="rounded-lg border border-line px-3 py-1.5 text-xs font-semibold text-ink-soft hover:border-accent hover:text-accent-ink">
                        {d.activo ? "Dar de baja" : "Reactivar"}
                      </button>
                    </FormularioCarga>
                  </div>
                ),
              }]),
            ]}
          />
          </ZonaCarga>

          <Paginacion pagina={lista.pagina} totalPaginas={totalPaginas} searchParams={{ buscar: sp.buscar }} />
      </div>
    </ListadoConCarga>
  );
}
