import Alerta from "@/components/common/Alerta";
import { BotonFiltrar, FiltrosCarga, LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import Campo, { claseControl } from "@/components/common/Campo";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Tarjeta from "@/components/common/Tarjeta";
import { exigirPermiso } from "@/lib/api";
import { obtenerTablero, type ResultadoTablero } from "@/lib/api-indicadores";
import { mensajeError } from "@/lib/formularios";
import { ETIQUETA_MOTIVO_MORA, formatearDias, formatearHora, periodosRapidos, plural } from "@/lib/indicadores";
import { formatearFecha } from "@/lib/legajos";
import { ETIQUETA_ESTADO, ETIQUETA_TIPO } from "@/lib/tramites";
import ActualizacionAutomatica from "./ActualizacionAutomatica";

const MINUTOS_REFRESCO = 5;

function mensajeConsulta(resultado: Extract<ResultadoTablero, { ok: false }>): string {
  if (resultado.codigo === "VALIDACION") return "Revisa el periodo: las fechas deben tener el formato AAAA-MM-DD.";
  return mensajeError(resultado.codigo, "No se pudo calcular el tablero. Intenta de nuevo en unos momentos.");
}

// Barra horizontal de una sola serie: crece desde la izquierda, extremo redondeado de 4px.
// `color` es una clase de fondo; por defecto el azul del gráfico. La pista gris da la escala completa.
function Barra({ valor, maximo, titulo, color = "bg-grafico" }: { valor: number; maximo: number; titulo: string; color?: string }) {
  const ancho = maximo > 0 ? (valor / maximo) * 100 : 0;
  return (
    <div className="flex h-6 min-w-16 flex-1 items-center" title={titulo}>
      <div className="h-3 w-full rounded-md bg-line/60">{valor > 0 && <div className={`h-3 min-w-[4px] rounded-md ${color}`} style={{ width: `${ancho}%` }} />}</div>
    </div>
  );
}

// Un color por estado: el azul se oscurece a medida que el trámite avanza; Observado (naranja) es lo único que pide acción.
const COLOR_ESTADO: Record<string, string> = {
  INGRESADO: "bg-[#a9c4e6] dark:bg-[#4c6a96]",
  ASIGNADO: "bg-[#8db0dc] dark:bg-[#5578a8]",
  EN_REVISION: "bg-[#5f93cf] dark:bg-[#6a96d8]",
  OBSERVADO: "bg-[#e08a1e]",
  DICTAMINADO: "bg-[#2f6fb5] dark:bg-[#5b8fd0]",
  EN_FACULTAD: "bg-[#255b99] dark:bg-[#4f82c4]",
  REVISION_FINAL: "bg-[#1d4f86] dark:bg-[#4577b6]",
  REMITIDO: "bg-[#173f6c] dark:bg-[#3d6aa8]",
  ARCHIVADO: "bg-[#0f2d52] dark:bg-[#3560a0]",
};

const TONO_INDICADOR = {
  azul: "bg-[#2f6fb5] dark:bg-[#5b8fd0]",
  turquesa: "bg-[#1f6f8b] dark:bg-[#4fb0cf]",
  naranja: "bg-[#e08a1e]",
  rojo: "bg-accent",
};

function Indicador({ etiqueta, valor, detalle, tono, alerta = false }: { etiqueta: string; valor: string; detalle: string; tono: keyof typeof TONO_INDICADOR; alerta?: boolean }) {
  return (
    <div className="rounded-xl border border-line bg-surface px-5 py-4">
      <span aria-hidden="true" className={`mb-2 block h-1 w-10 rounded-full ${TONO_INDICADOR[tono]}`} />
      <p className="text-sm font-semibold text-ink-soft">{etiqueta}</p>
      <p className={`mt-1 font-display text-4xl font-semibold leading-tight tracking-[-0.02em] ${alerta ? "text-accent-ink" : "text-ink"}`}>{valor}</p>
      <p className="text-sm text-ink-soft">{detalle}</p>
    </div>
  );
}

// CU-25: indicadores de gestión del periodo y alertas de carpetas en mora (solo lectura).
export default async function TableroPage({ searchParams }: { searchParams: Promise<{ desde?: string; hasta?: string }> }) {
  const sesion = await exigirPermiso("TABLERO_VER");
  const sp = await searchParams;
  const periodo = { desde: sp.desde, hasta: sp.hasta };
  const resultado = await obtenerTablero(periodo);
  const hayPeriodo = Boolean(sp.desde || sp.hasta);
  const puedeCustodia = sesion.permisos.includes("CUSTODIA_REGISTRAR");

  const textoPeriodo =
    sp.desde && sp.hasta
      ? `del ${formatearFecha(sp.desde)} al ${formatearFecha(sp.hasta)}`
      : sp.desde
        ? `desde el ${formatearFecha(sp.desde)}`
        : sp.hasta
          ? `hasta el ${formatearFecha(sp.hasta)}`
          : "en todo el historial";

  return (
    <ListadoConCarga>
      <ActualizacionAutomatica minutos={MINUTOS_REFRESCO} />
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Tablero de indicadores"
          descripcion={
            resultado.ok ? (
              <>
                Actualizado a las {formatearHora(resultado.tablero.generadoEn)} (hora de Bolivia). Se actualiza solo cada {MINUTOS_REFRESCO} minutos.
              </>
            ) : (
              "Indicadores de gestión de trámites y custodia de carpetas."
            )
          }
        />

        {/* Un solo filtro para todo el tablero: el periodo acota los trámites radicados y resueltos. */}
        <div className="flex flex-col gap-3">
          <FiltrosCarga key={JSON.stringify(periodo)} action="/panel/tablero" className="flex flex-wrap items-end gap-3">
            <Campo id="desde" etiqueta="Desde">
              <input id="desde" type="date" name="desde" defaultValue={sp.desde} className={claseControl()} />
            </Campo>
            <Campo id="hasta" etiqueta="Hasta">
              <input id="hasta" type="date" name="hasta" defaultValue={sp.hasta} className={claseControl()} />
            </Campo>
            <BotonFiltrar variante="primario">Aplicar periodo</BotonFiltrar>
          </FiltrosCarga>
          <div className="flex flex-wrap items-center gap-2">
            {periodosRapidos().map((p) => {
              const activo = sp.desde === p.desde && sp.hasta === p.hasta;
              return (
                <LinkCarga
                  key={p.etiqueta}
                  href={`/panel/tablero?desde=${p.desde}&hasta=${p.hasta}`}
                  aria-current={activo ? "true" : undefined}
                  className={`rounded-full border px-3 py-1 text-sm font-semibold transition-colors ${activo ? "border-ink bg-ink text-surface" : "border-line bg-surface text-ink-soft hover:border-accent hover:text-accent-ink"}`}
                >
                  {p.etiqueta}
                </LinkCarga>
              );
            })}
            <LinkCarga
              href="/panel/tablero"
              aria-current={!hayPeriodo ? "true" : undefined}
              className={`rounded-full border px-3 py-1 text-sm font-semibold transition-colors ${!hayPeriodo ? "border-ink bg-ink text-surface" : "border-line bg-surface text-ink-soft hover:border-accent hover:text-accent-ink"}`}
            >
              Todo el historial
            </LinkCarga>
          </div>
        </div>

        {!resultado.ok ? (
          <Alerta tipo="error">{mensajeConsulta(resultado)}</Alerta>
        ) : (
          <ZonaCarga>
            <Contenido tablero={resultado.tablero} textoPeriodo={textoPeriodo} puedeCustodia={puedeCustodia} />
          </ZonaCarga>
        )}
      </div>
    </ListadoConCarga>
  );
}

function Contenido({
  tablero: t,
  textoPeriodo,
  puedeCustodia,
}: {
  tablero: Extract<ResultadoTablero, { ok: true }>["tablero"];
  textoPeriodo: string;
  puedeCustodia: boolean;
}) {
  const maxEstado = Math.max(0, ...t.porEstado.map((e) => e.cantidad));
  const maxTiempo = Math.max(0, ...t.tiempos.map((g) => g.promedioDias));
  const maxCarga = Math.max(0, ...t.analistas.map((a) => a.activos));

  return (
    <div className="flex flex-col gap-6">
      {t.alertas.length > 0 && (
        <section aria-label="Aviso de mora" className="flex flex-wrap items-center gap-x-4 gap-y-2 rounded-xl border border-[#e8b95a] bg-[#fff4e0] px-5 py-4 dark:border-[#8a6a2a] dark:bg-[#33270f]">
          <svg className="shrink-0 text-[#8a4b00] dark:text-[#f0b860]" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
            <path d="M12 9v4M12 17h.01" />
            <path d="M10.3 3.9 1.8 18a2 2 0 0 0 1.7 3h17a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z" />
          </svg>
          <p className="min-w-0 flex-[1_1_320px] text-sm text-[#5c3300] dark:text-[#f5d9a8]">
            <strong>{plural(t.alertas.length, "carpeta retenida", "carpetas retenidas")} fuera de plazo.</strong> La alerta se cierra sola al registrar la devolución.
          </p>
          <a href="#alertas" className="text-sm font-bold text-[#5c3300] underline dark:text-[#f5d9a8]">
            Ver carpetas
          </a>
        </section>
      )}

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Indicador tono="azul" etiqueta="Trámites radicados" valor={t.tramites.radicados.toLocaleString("es-BO")} detalle={textoPeriodo} />
        <Indicador tono="naranja" etiqueta="En curso" valor={t.tramites.enCurso.toLocaleString("es-BO")} detalle="De los radicados, aún sin remitir al Rectorado" />
        <Indicador
          tono="turquesa"
          etiqueta="Tiempo promedio de resolución"
          valor={t.tramites.resueltos > 0 ? `${formatearDias(t.tramites.promedioDias)} días` : "—"}
          detalle={t.tramites.resueltos > 0 ? `${plural(t.tramites.resueltos, "trámite resuelto", "trámites resueltos")} ${textoPeriodo}` : `Ningún trámite resuelto ${textoPeriodo}`}
        />
        <Indicador
          tono="rojo"
          etiqueta="Carpetas en mora"
          valor={t.alertas.length.toLocaleString("es-BO")}
          detalle={`De ${plural(t.prestamosVigentes, "carpeta prestada", "carpetas prestadas")} hoy`}
          alerta={t.alertas.length > 0}
        />
      </div>

      <div className="grid gap-6 xl:grid-cols-2">
        <Tarjeta titulo="Trámites por estado" descripcion={`Radicados ${textoPeriodo}, según su estado actual.`}>
          <ul className="flex flex-col gap-1">
            {t.porEstado.map((e) => (
              <li key={e.estado} className="flex items-center gap-3 text-sm">
                <span className="w-40 shrink-0 text-ink">{ETIQUETA_ESTADO[e.estado] ?? e.estado}</span>
                <Barra valor={e.cantidad} maximo={maxEstado} color={COLOR_ESTADO[e.estado]} titulo={`${ETIQUETA_ESTADO[e.estado] ?? e.estado}: ${e.cantidad}`} />
                <span className="w-10 shrink-0 text-right font-semibold tabular-nums text-ink">{e.cantidad}</span>
              </li>
            ))}
          </ul>
        </Tarjeta>

        <Tarjeta titulo="Carga por analista" descripcion={`Trámites a cargo hoy y resueltos ${textoPeriodo}. Fuera de Escalafón no cuenta como carga.`}>
          {t.analistas.length === 0 ? (
            <p className="text-sm text-ink-faint">No hay analistas activos.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[440px] border-collapse text-sm">
                <thead>
                  <tr className="text-left text-xs font-semibold uppercase tracking-wide text-ink-faint">
                    <th className="pb-2 pr-3">Analista</th>
                    <th className="pb-2 pr-3">A cargo</th>
                    <th className="w-24 pb-2 pr-3 text-right">Fuera de Esc.</th>
                    <th className="w-24 pb-2 text-right">Resueltos</th>
                  </tr>
                </thead>
                <tbody>
                  {t.analistas.map((a) => (
                    <tr key={a.id} className="border-t border-line">
                      <td className="py-2 pr-3 font-semibold text-ink">{a.nombre}</td>
                      <td className="py-2 pr-3">
                        <div className="flex items-center gap-2">
                          <Barra valor={a.activos} maximo={maxCarga} titulo={`${a.nombre}: ${a.activos} a cargo`} />
                          <span className="w-8 text-right tabular-nums">{a.activos}</span>
                        </div>
                      </td>
                      <td className="py-2 pr-3 text-right tabular-nums text-ink-soft">{a.enFacultad}</td>
                      <td className="py-2 text-right tabular-nums text-ink-soft">{a.resueltos}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </Tarjeta>
      </div>

      <section id="alertas" className="scroll-mt-6">
      <Tarjeta
        titulo="Carpetas retenidas fuera de plazo"
        descripcion={`Carpetas fuera del archivo con la fecha límite vencida o con más de ${t.diasMaximosPrestamo} días de préstamo. Cada alerta se cierra sola al registrar la devolución.`}
      >
        {t.alertas.length === 0 ? (
          <Alerta tipo="exito">Ninguna carpeta está en mora.</Alerta>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[760px] border-collapse text-sm">
              <thead>
                <tr className="text-left text-xs font-semibold uppercase tracking-wide text-ink-faint">
                  <th className="pb-2 pr-3">Legajo</th>
                  <th className="pb-2 pr-3">Docente</th>
                  <th className="pb-2 pr-3">Custodio</th>
                  <th className="pb-2 pr-3">Salida / límite</th>
                  <th className="pb-2 pr-3">Atraso</th>
                  {puedeCustodia && <th className="pb-2" />}
                </tr>
              </thead>
              <tbody>
                {t.alertas.map((a) => (
                  <tr key={a.legajoCodigo} className="border-t border-line align-top">
                    <td className="py-3 pr-3 font-mono text-sm font-semibold text-ink">{a.legajoCodigo}</td>
                    <td className="py-3 pr-3">
                      <span className="font-semibold text-ink">{a.docenteNombre}</span>
                      {a.tramiteCodigo && <span className="block font-mono text-xs text-ink-faint">{a.tramiteCodigo}</span>}
                    </td>
                    <td className="py-3 pr-3 text-ink">{a.custodio}</td>
                    <td className="py-3 pr-3 text-sm text-ink-soft">
                      {formatearFecha(a.salida)}
                      {a.fechaLimite && <span className="block">límite {formatearFecha(a.fechaLimite)}</span>}
                    </td>
                    <td className="py-3 pr-3">
                      <span className="inline-block rounded-full bg-accent-soft px-2.5 py-0.5 text-xs font-semibold text-accent-ink">{plural(a.diasMora, "día", "días")}</span>
                      <span className="block text-xs text-ink-faint">{ETIQUETA_MOTIVO_MORA[a.motivo]}</span>
                    </td>
                    {puedeCustodia && (
                      <td className="py-3 text-right">
                        <LinkCarga href="/panel/custodia?tab=vencidas" className="whitespace-nowrap text-sm font-semibold text-accent-ink hover:underline">
                          Registrar devolución
                        </LinkCarga>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Tarjeta>
      </section>

      <Tarjeta
        titulo="Tiempo de resolución por tipo y categoría"
        descripcion={`Días entre la radicación y la remisión al Rectorado de los trámites resueltos ${textoPeriodo}.`}
      >
        {t.tiempos.length === 0 ? (
          <p className="text-sm text-ink-faint">Ningún trámite se resolvió en este periodo.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[640px] border-collapse text-sm">
              <thead>
                <tr className="text-left text-xs font-semibold uppercase tracking-wide text-ink-faint">
                  <th className="pb-2 pr-3">Tipo</th>
                  <th className="pb-2 pr-3">Categoría</th>
                  <th className="w-20 pb-2 pr-3 text-right">Trámites</th>
                  <th className="pb-2 pr-3">Promedio</th>
                  <th className="w-28 pb-2 pr-3 text-right">Más rápido</th>
                  <th className="w-28 pb-2 text-right">Más lento</th>
                </tr>
              </thead>
              <tbody>
                {t.tiempos.map((g) => (
                  <tr key={`${g.tipo}-${g.categoria}`} className="border-t border-line">
                    <td className="py-2 pr-3 text-ink-soft">{ETIQUETA_TIPO[g.tipo] ?? g.tipo}</td>
                    <td className="py-2 pr-3 font-semibold text-ink">{g.categoria}</td>
                    <td className="py-2 pr-3 text-right tabular-nums">{g.cantidad}</td>
                    <td className="py-2 pr-3">
                      <div className="flex items-center gap-2">
                        <Barra valor={g.promedioDias} maximo={maxTiempo} color="bg-[#1f6f8b] dark:bg-[#4fb0cf]" titulo={`${g.categoria}: ${formatearDias(g.promedioDias)} días en promedio`} />
                        <span className="w-20 text-right tabular-nums">{formatearDias(g.promedioDias)} días</span>
                      </div>
                    </td>
                    <td className="py-2 pr-3 text-right tabular-nums text-ink-soft">{formatearDias(g.minimoDias)} días</td>
                    <td className="py-2 text-right tabular-nums text-ink-soft">{formatearDias(g.maximoDias)} días</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Tarjeta>
    </div>
  );
}
