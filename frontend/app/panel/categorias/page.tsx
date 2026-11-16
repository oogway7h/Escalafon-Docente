import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import { FormularioCarga, LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Insignia from "@/components/common/Insignia";
import Tabla from "@/components/common/Tabla";
import { IconoEditar, IconoEliminar } from "@/components/layout/IconosPanel";
import { exigirPermiso } from "@/lib/api";
import { listarCategorias } from "@/lib/api-tramites";
import { mensajeError } from "@/lib/formularios";
import { cambiarEstadoCategoriaAction } from "./actions";

export default async function CategoriasPage({ searchParams }: { searchParams: Promise<{ error?: string }> }) {
  await exigirPermiso("CATEGORIA_GESTIONAR");
  const sp = await searchParams;
  const categorias = await listarCategorias();

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Categorías"
          descripcion={<>Catálogo normativo: {categorias.length} categorías con sus documentos probatorios.</>}
          accion={
            <LinkCarga href="/panel/categorias/nueva" className={claseBoton("primario", "sm")}>
              + Nueva categoría
            </LinkCarga>
          }
        />

        {sp.error && <Alerta tipo="error">{mensajeError(sp.error, "No se pudo completar la acción.")}</Alerta>}
          <ZonaCarga>
        <Tabla
          filas={categorias}
          clave={(c) => c.id}
          pie="Dar de baja oculta la categoría al radicar pero la conserva. Eliminar la borra de la base de datos y solo es posible si ningún docente ni trámite la usa."
          columnas={[
            { titulo: "Categoría", celda: (c) => <span className="font-semibold">{c.nombre}</span> },
            { titulo: "Puntaje", celda: (c) => (c.puntajeDesde === null ? "—" : c.puntajeHasta === null ? `${c.puntajeDesde} o más` : `${c.puntajeDesde} a ${c.puntajeHasta}`) },
            {
              titulo: "Documentos exigidos",
              celda: (c) => <span className="text-sm text-ink-soft">{c.requisitos.map((r) => r.nombre).join(" · ") || "—"}</span>,
            },
            {
              titulo: "En uso",
              celda: (c) =>
                c.docentes + c.tramites === 0 ? (
                  <span className="text-ink-faint">—</span>
                ) : (
                  <span className="text-sm text-ink-soft">
                    {c.docentes} {c.docentes === 1 ? "docente" : "docentes"} · {c.tramites} {c.tramites === 1 ? "trámite" : "trámites"}
                  </span>
                ),
            },
            { titulo: "Estado", celda: (c) => <Insignia tono={c.activa ? "exito" : "neutro"}>{c.activa ? "Activa" : "De baja"}</Insignia> },
            {
              className: "w-60 text-right",
              celda: (c) => (
                <div className="flex items-center justify-end gap-2">
                  <LinkCarga
                    href={`/panel/categorias/${c.id}/editar`}
                    className="flex h-8 w-8 items-center justify-center rounded-lg border border-line text-ink-soft hover:border-accent hover:text-accent-ink"
                    aria-label={`Editar ${c.nombre}`}
                  >
                    <IconoEditar className="h-3.5 w-3.5" />
                  </LinkCarga>
                  <FormularioCarga action={cambiarEstadoCategoriaAction}>
                    <input type="hidden" name="id" value={c.id} />
                    <input type="hidden" name="activo" value={String(!c.activa)} />
                    <button type="submit" className="rounded-lg border border-line px-3 py-1.5 text-xs font-semibold text-ink-soft hover:border-accent hover:text-accent-ink">
                      {c.activa ? "Dar de baja" : "Reactivar"}
                    </button>
                  </FormularioCarga>
                  <LinkCarga
                    href={`/panel/categorias/${c.id}/eliminar`}
                    className="flex h-8 w-8 items-center justify-center rounded-lg border border-accent-soft text-accent-ink hover:bg-accent-soft"
                    aria-label={`Eliminar ${c.nombre}`}
                  >
                    <IconoEliminar className="h-3.5 w-3.5" />
                  </LinkCarga>
                </div>
              ),
            },
          ]}
        />
          </ZonaCarga>
      </div>
    </ListadoConCarga>
  );
}
