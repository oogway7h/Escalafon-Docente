import Form from "next/form";
import Alerta from "@/components/common/Alerta";
import { LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import { exigirPermiso } from "@/lib/api";
import { listarSolicitudes, type SolicitudCarpeta } from "@/lib/api-legajos";
import { mensajeError } from "@/lib/formularios";
import { enDias, formatearFecha } from "@/lib/legajos";
import { DevolucionForm, EntregaForm } from "./AccionesCarpeta";

const MAXIMO = 50; // ponytail: sin paginación; la cola de archivo rara vez pasa de unas decenas

type Pestana = "atender" | "prestadas" | "vencidas";

const EXITOS: Record<string, string> = {
  ENTREGADA: "Carpeta entregada: la salida quedó en la bitácora.",
  APROBADA: "Solicitud aprobada: la carpeta queda lista para entregar.",
  RECHAZADA: "Solicitud rechazada: el solicitante verá el motivo.",
  DEVUELTA: "Devolución registrada: la carpeta volvió al archivo.",
};

const DIA = 86_400_000;
const diasEntre = (desde: string, hasta: string) => Math.round((Date.parse(hasta) - Date.parse(desde)) / DIA);

function hace(iso: string): string {
  const min = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 60_000));
  if (min < 60) return "hace un momento";
  if (min < 24 * 60) return `hace ${Math.floor(min / 60)} h`;
  const d = Math.floor(min / (24 * 60));
  return d === 1 ? "ayer" : `hace ${d} días`;
}

// Pantalla única de Archivo: atender las solicitudes, entregar las carpetas y recibirlas de vuelta (CU-21 al CU-23).
// «Aprobar y entregar» hace en un solo paso lo que antes eran dos pantallas.
export default async function CustodiaPage({
  searchParams,
}: {
  searchParams: Promise<{ tab?: string; sel?: string; buscar?: string; exito?: string; error?: string }>;
}) {
  await exigirPermiso("CUSTODIA_REGISTRAR");
  const sp = await searchParams;
  const [pendientes, aprobadas, prestadas] = await Promise.all([
    listarSolicitudes({ estado: "PENDIENTE" }, 1, MAXIMO),
    listarSolicitudes({ estado: "ATENDIDA" }, 1, MAXIMO),
    listarSolicitudes({ estado: "DESPACHADA" }, 1, MAXIMO),
  ]);

  const hoy = enDias(0);
  const q = (sp.buscar ?? "").trim().toLowerCase();
  const coincide = (s: SolicitudCarpeta) => !q || [s.legajoCodigo, s.docenteNombre, s.solicitante, s.tramiteCodigo].some((v) => v.toLowerCase().includes(q));
  const esVencida = (s: SolicitudCarpeta) => s.fechaLimite !== null && s.fechaLimite < hoy;

  const porAtender = [...pendientes.datos, ...aprobadas.datos];
  const prestadasOrden = [...prestadas.datos].sort((a, b) => (a.fechaLimite ?? "9").localeCompare(b.fechaLimite ?? "9"));
  const conteos = { atender: pendientes.total + aprobadas.total, prestadas: prestadas.total, vencidas: prestadas.datos.filter(esVencida).length };

  const tab: Pestana = sp.tab === "prestadas" || sp.tab === "vencidas" ? sp.tab : "atender";
  const lista = (tab === "atender" ? porAtender : tab === "prestadas" ? prestadasOrden : prestadasOrden.filter(esVencida)).filter(coincide);
  const actual = lista.find((s) => s.id === sp.sel) ?? lista[0] ?? null;

  const url = (p: { tab?: Pestana; sel?: string }) => {
    const qs = new URLSearchParams();
    const t = p.tab ?? tab;
    if (t !== "atender") qs.set("tab", t);
    if (p.sel) qs.set("sel", p.sel);
    if (sp.buscar) qs.set("buscar", sp.buscar);
    const s = qs.toString();
    return s ? `/panel/custodia?${s}` : "/panel/custodia";
  };

  const pestana = (id: Pestana, etiqueta: string, n: number, alerta = false) => {
    const activa = tab === id;
    return (
      <LinkCarga
        key={id}
        href={url({ tab: id })}
        aria-current={activa ? "page" : undefined}
        className={`flex items-center gap-2.5 rounded-xl border px-[18px] py-3 text-base font-semibold transition-colors ${activa ? "border-ink bg-ink text-surface" : "border-line bg-surface text-ink hover:border-accent"}`}
      >
        {etiqueta}
        <span className={`min-w-6 rounded-full px-2 py-px text-center text-xs ${activa ? "bg-white/20 text-white" : alerta && n > 0 ? "bg-accent-soft text-accent-dark" : "bg-bg text-ink-soft"}`}>{n}</span>
      </LinkCarga>
    );
  };

  return (
    <ListadoConCarga>
      <div className="mx-auto flex w-full max-w-[1280px] flex-col gap-6">
        <header className="flex flex-wrap items-end justify-between gap-4">
          <div>
            <h1 className="font-display text-3xl font-semibold tracking-[-0.02em] text-ink">Solicitudes y custodia</h1>
            <p className="text-sm text-ink-soft">Atiende, entrega y recibe carpetas desde un solo lugar. Cada movimiento queda en la bitácora.</p>
          </div>
          <Form action="/panel/custodia" scroll={false} className="flex w-[340px] max-w-full items-center gap-2.5 rounded-xl border border-line bg-surface px-4 focus-within:border-accent">
            {tab !== "atender" && <input type="hidden" name="tab" value={tab} />}
            <svg className="shrink-0 text-ink-soft" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
              <circle cx="11" cy="11" r="7" />
              <path d="m20 20-3.5-3.5" />
            </svg>
            <input type="search" name="buscar" defaultValue={sp.buscar} placeholder="Carpeta, docente o solicitante" aria-label="Buscar carpeta, docente o solicitante" className="h-[46px] min-w-0 flex-1 bg-transparent text-base text-ink outline-none placeholder:text-ink-faint" />
          </Form>
        </header>

        {sp.exito && EXITOS[sp.exito] && <Alerta tipo="exito">{EXITOS[sp.exito]}</Alerta>}
        {sp.error && <Alerta tipo="error">{mensajeError(sp.error, "No se pudo completar la acción.")}</Alerta>}

        <div className="flex flex-wrap gap-2.5">
          {pestana("atender", "Por atender", conteos.atender, true)}
          {pestana("prestadas", "Prestadas", conteos.prestadas)}
          {pestana("vencidas", "Vencidas", conteos.vencidas, true)}
        </div>

        <ZonaCarga>
          <div className="flex flex-wrap items-start gap-6">
            <section className="flex min-w-0 flex-[999_1_520px] flex-col gap-3" aria-label="Solicitudes">
              {lista.length === 0 && (
                <p className="rounded-xl border border-line bg-surface p-6 text-sm text-ink-soft">
                  {q
                    ? "Ninguna carpeta coincide con la búsqueda."
                    : tab === "atender"
                      ? "No hay solicitudes por atender."
                      : tab === "prestadas"
                        ? "No hay carpetas fuera del archivo."
                        : "No hay carpetas vencidas."}
                </p>
              )}
              {lista.map((s) => {
                const elegida = actual?.id === s.id;
                const vence = s.fechaLimite ? diasEntre(hoy, s.fechaLimite) : null;
                const vencida = vence !== null && vence < 0;
                const etiqueta =
                  s.estado === "PENDIENTE"
                    ? { texto: "Pendiente", clase: "bg-accent-soft text-accent-dark" }
                    : s.estado === "ATENDIDA"
                      ? { texto: "Aprobada · falta entregar", clase: "bg-success/15 text-success" }
                      : vencida
                        ? { texto: `Vencida hace ${-vence!} ${-vence! === 1 ? "día" : "días"}`, clase: "bg-accent text-white" }
                        : { texto: vence === 0 ? "Vence hoy" : `Vence en ${vence} ${vence === 1 ? "día" : "días"}`, clase: "bg-bg text-ink-soft" };
                return (
                  <LinkCarga
                    key={s.id}
                    href={url({ sel: s.id })}
                    scroll={false}
                    aria-current={elegida ? "true" : undefined}
                    className={`flex flex-col gap-2 rounded-xl border-2 px-[18px] py-4 transition-colors ${elegida ? "border-accent bg-accent-soft/40" : "border-transparent bg-surface shadow-[0_0_0_1px_var(--color-line)] hover:shadow-[0_0_0_1px_var(--color-accent)]"}`}
                  >
                    <span className="flex flex-wrap items-center gap-x-2.5 gap-y-1">
                      <span className="font-mono text-sm font-semibold text-ink">{s.legajoCodigo}</span>
                      <span className="text-base font-semibold text-ink">{s.docenteNombre}</span>
                      <span className={`ml-auto rounded-full px-2.5 py-0.5 text-xs font-semibold ${etiqueta.clase}`}>{etiqueta.texto}</span>
                    </span>
                    <span className="text-sm text-ink-soft">{s.estado === "DESPACHADA" ? `En poder de ${s.solicitante}` : s.motivo}</span>
                    <span className="text-xs text-ink-faint">
                      {s.estado === "DESPACHADA"
                        ? `Entregada ${formatearFecha(s.atendidoEn)} · debía volver el ${formatearFecha(s.fechaLimite)}`
                        : `${s.solicitante}${s.tramiteCodigo ? ` · trámite ${s.tramiteCodigo}` : ""} · ${hace(s.creadoEn)}`}
                    </span>
                  </LinkCarga>
                );
              })}
            </section>

            {actual && (
              <aside className="flex min-w-[320px] max-w-[460px] flex-[1_1_400px] flex-col gap-[18px] rounded-xl border border-line bg-surface p-6 shadow-[0_1px_2px_rgba(15,27,46,0.04),0_12px_32px_rgba(15,27,46,0.06)] lg:sticky lg:top-6">
                <header className="flex flex-col gap-1">
                  <span className="font-mono text-sm font-semibold text-ink-soft">{actual.legajoCodigo}</span>
                  <h2 className="font-display text-xl font-semibold tracking-[-0.01em] text-ink">{actual.docenteNombre}</h2>
                </header>

                {actual.estado === "DESPACHADA" && actual.fechaLimite && actual.fechaLimite < hoy && (
                  <p className="rounded-xl bg-accent-soft px-3.5 py-3 text-sm font-semibold leading-snug text-accent-dark">
                    Vencida hace {diasEntre(actual.fechaLimite, hoy)} {diasEntre(actual.fechaLimite, hoy) === 1 ? "día" : "días"}. Quien la recibe queda como último custodio responsable.
                  </p>
                )}

                <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
                  <dt className="text-ink-faint">{actual.estado === "DESPACHADA" ? "En poder de" : "Solicita"}</dt>
                  <dd className="text-right font-semibold text-ink">{actual.solicitante}</dd>
                  {actual.tramiteCodigo && (
                    <>
                      <dt className="text-ink-faint">Trámite</dt>
                      <dd className="text-right font-mono text-sm font-semibold text-ink">{actual.tramiteCodigo}</dd>
                    </>
                  )}
                  <dt className="text-ink-faint">Motivo</dt>
                  <dd className="text-right font-semibold text-ink">{actual.motivo}</dd>
                  {actual.estado === "DESPACHADA" && (
                    <>
                      <dt className="text-ink-faint">Entregada</dt>
                      <dd className="text-right font-semibold text-ink">{formatearFecha(actual.atendidoEn)}</dd>
                      <dt className="text-ink-faint">Debe volver</dt>
                      <dd className="text-right font-semibold text-ink">{formatearFecha(actual.fechaLimite)}</dd>
                    </>
                  )}
                </dl>

                <hr className="border-line" />

                {actual.estado === "DESPACHADA" ? (
                  <DevolucionForm key={actual.id} legajoId={actual.legajoId} id={actual.id} tab={tab} />
                ) : (
                  <EntregaForm key={actual.id} id={actual.id} tab={tab} aprobada={actual.estado === "ATENDIDA"} />
                )}
              </aside>
            )}
          </div>
        </ZonaCarga>
      </div>
    </ListadoConCarga>
  );
}
