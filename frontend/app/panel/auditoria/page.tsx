import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import { BotonFiltrar, FiltrosCarga, LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import Campo, { claseControl } from "@/components/common/Campo";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Insignia from "@/components/common/Insignia";
import Paginacion from "@/components/common/Paginacion";
import Tabla from "@/components/common/Tabla";
import { exigirPermiso } from "@/lib/api";
import { listarAuditoria, obtenerFiltrosAuditoria, type ResultadoAuditoria } from "@/lib/api-auditoria";
import { etiquetaAccion, formatearFechaHora, resumirDetalle, tonoAccion } from "@/lib/auditoria";
import { mensajeError } from "@/lib/formularios";

const ETIQUETA_PARAMETRO: Record<string, string> = {
  usuarioId: "usuario",
  desde: "fecha inicial",
  hasta: "fecha final",
  pagina: "página",
  tamano: "tamaño",
};

function mensajeConsulta(resultado: Extract<ResultadoAuditoria, { ok: false }>): string {
  if (resultado.codigo === "VALIDACION" && resultado.detalles) {
    const campos = Object.entries(resultado.detalles).map(([campo, msj]) => `${ETIQUETA_PARAMETRO[campo] ?? campo}: ${msj}`);
    return `Revisa los filtros. ${campos.join(". ")}.`;
  }
  return mensajeError(resultado.codigo, "No se pudo consultar la auditoría. Intenta de nuevo en unos momentos.");
}

export default async function AuditoriaPage({
  searchParams,
}: {
  searchParams: Promise<{ usuarioId?: string; accion?: string; desde?: string; hasta?: string; pagina?: string }>;
}) {
  await exigirPermiso("AUDITORIA_CONSULTAR");
  const sp = await searchParams;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const filtro = { usuarioId: sp.usuarioId, accion: sp.accion, desde: sp.desde, hasta: sp.hasta };
  const [resultado, opciones] = await Promise.all([listarAuditoria(filtro, pagina), obtenerFiltrosAuditoria()]);
  const hayFiltros = Object.values(filtro).some(Boolean);

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Auditoría"
          descripcion={
            resultado.ok ? (
              <>
                {resultado.lista.total} {resultado.lista.total === 1 ? "registro" : "registros"}
                {hayFiltros ? " coinciden con los filtros" : " en la bitácora de acciones de los usuarios del sistema"}.
              </>
            ) : (
              "Bitácora de las acciones realizadas por los usuarios del sistema."
            )
          }
        />

        {/* key: al limpiar o cambiar los filtros por URL, los campos se reinician con los valores vigentes. */}
        <FiltrosCarga
          key={JSON.stringify(filtro)}
          action="/panel/auditoria"
          className="grid items-end gap-3 sm:grid-cols-2 lg:grid-cols-[minmax(0,1.4fr)_minmax(0,1.2fr)_minmax(0,1fr)_minmax(0,1fr)_auto]"
        >
          <Campo id="usuarioId" etiqueta="Usuario">
            <select
              id="usuarioId"
              name="usuarioId"
              defaultValue={sp.usuarioId ?? ""}
              disabled={!opciones}
              className={claseControl({ bloqueado: !opciones })}
            >
              <option value="">Todos los usuarios</option>
              {opciones?.usuarios.map((u) => (
                <option key={u.id} value={u.id}>
                  {u.apellidos}, {u.nombres} — {u.correo}
                </option>
              ))}
            </select>
          </Campo>
          <Campo id="accion" etiqueta="Tipo de acción">
            <select
              id="accion"
              name="accion"
              defaultValue={sp.accion ?? ""}
              disabled={!opciones}
              className={claseControl({ bloqueado: !opciones })}
            >
              <option value="">Todas las acciones</option>
              {opciones?.acciones.map((a) => (
                <option key={a} value={a}>
                  {etiquetaAccion(a)}
                </option>
              ))}
            </select>
          </Campo>
          <Campo id="desde" etiqueta="Desde">
            <input id="desde" type="date" name="desde" defaultValue={sp.desde} className={claseControl()} />
          </Campo>
          <Campo id="hasta" etiqueta="Hasta">
            <input id="hasta" type="date" name="hasta" defaultValue={sp.hasta} className={claseControl()} />
          </Campo>
          <div className="flex gap-2.5">
            <BotonFiltrar variante="primario">Aplicar filtros</BotonFiltrar>
            {hayFiltros ? (
              <LinkCarga href="/panel/auditoria" className={claseBoton("secundario", "sm")}>
                Limpiar
              </LinkCarga>
            ) : (
              <button type="reset" className={claseBoton("secundario", "sm")}>
                Limpiar
              </button>
            )}
          </div>
        </FiltrosCarga>

        {!opciones && resultado.ok && (
          <Alerta tipo="aviso">No se pudo cargar la lista de usuarios y acciones para los filtros.</Alerta>
        )}

        {!resultado.ok ? (
          <Alerta tipo="error">{mensajeConsulta(resultado)}</Alerta>
        ) : (
          <>
            <ZonaCarga>
              {/* En pantallas angostas la tabla se desplaza horizontalmente en vez de recortarse. */}
              <div className="overflow-x-auto rounded-xl">
                <Tabla
                  className="min-w-[760px]"
                  filas={resultado.lista.datos}
                  clave={(r) => r.id}
                  vacio={
                    resultado.lista.total > 0
                      ? "No hay registros en esta página."
                      : hayFiltros
                        ? "No hay registros que coincidan con los filtros."
                        : "Todavía no hay registros de auditoría."
                  }
                  pie="Los registros de auditoría son de solo lectura: no se pueden modificar ni eliminar. Las fechas están en hora de Bolivia."
                  columnas={[
                    {
                      titulo: "Fecha y hora",
                      className: "w-44 whitespace-nowrap",
                      celda: (r) => <span className="text-sm tabular-nums">{formatearFechaHora(r.fecha)}</span>,
                    },
                    {
                      titulo: "Usuario",
                      celda: (r) =>
                        r.usuario ? (
                          <>
                            <div className="font-semibold">
                              {r.usuario.nombres} {r.usuario.apellidos}
                            </div>
                            <div className="text-xs text-ink-faint">{r.usuario.correo}</div>
                          </>
                        ) : (
                          <span className="text-sm text-ink-faint">
                            {r.accion === "RADICAR_TRAMITE_PORTAL" ? "Portal público" : "No identificado"}
                          </span>
                        ),
                    },
                    {
                      titulo: "Acción",
                      celda: (r) => <Insignia tono={tonoAccion(r.accion)}>{etiquetaAccion(r.accion)}</Insignia>,
                    },
                    {
                      titulo: "Detalle",
                      celda: (r) => {
                        const lineas = resumirDetalle(r.detalle);
                        return lineas.length === 0 ? (
                          <span className="text-ink-faint">—</span>
                        ) : (
                          <ul className="max-w-[360px] text-xs text-ink-soft">
                            {lineas.map((l, i) => (
                              <li key={i} className="break-words">
                                {l}
                              </li>
                            ))}
                          </ul>
                        );
                      },
                    },
                    {
                      titulo: "Dirección IP",
                      className: "w-36 whitespace-nowrap",
                      celda: (r) =>
                        r.direccionIp ? (
                          <span className="font-mono text-xs">{r.direccionIp}</span>
                        ) : (
                          <span className="text-ink-faint">—</span>
                        ),
                    },
                  ]}
                />
              </div>
            </ZonaCarga>

            <Paginacion
              pagina={resultado.lista.pagina}
              totalPaginas={Math.max(1, Math.ceil(resultado.lista.total / resultado.lista.tamano))}
              searchParams={filtro}
            />
          </>
        )}
      </div>
    </ListadoConCarga>
  );
}
