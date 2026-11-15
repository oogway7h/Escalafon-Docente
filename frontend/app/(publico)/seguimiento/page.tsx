import type { Metadata } from "next";
import Link from "next/link";
import Alerta from "@/components/common/Alerta";
import FasesTramite from "@/components/tramites/FasesTramite";
import { consultarSeguimiento } from "@/lib/api-tramites";
import { DESCRIPCION_ESTADO, ETIQUETA_ESTADO, ETIQUETA_TIPO, TITULO_ESTADO, formatearFecha } from "@/lib/tramites";
import BuscadorSeguimiento from "./BuscadorSeguimiento";

export const metadata: Metadata = { title: "Seguimiento de trámite — Escalafón Docente" };

export default async function SeguimientoPage({
  searchParams,
}: {
  searchParams: Promise<{ codigo?: string; nueva?: string }>;
}) {
  const { codigo, nueva } = await searchParams;
  const vista = codigo ? await consultarSeguimiento(codigo) : null;
  const pendientes = vista?.observaciones.filter((o) => !o.resuelta) ?? [];
  // Las etapas llegan en orden cronológico; el historial muestra la más reciente primero.
  const historial = [...(vista?.etapas ?? [])].reverse();

  return (
    <div className="mx-auto flex w-full max-w-3xl flex-col gap-8">
      <section className="flex flex-col gap-4 rounded-[24px] bg-sidebar px-6 py-8 text-white sm:px-10 sm:py-10">
        <h1 className="font-display text-[clamp(1.75rem,4.5vw,2.75rem)] font-extrabold leading-[1.05] tracking-[-0.025em] text-balance">
          Consulta el estado de tu trámite
        </h1>
        <p className="max-w-xl text-[15.5px] leading-relaxed text-white/80">
          Ingresa el código que recibiste al radicar tu solicitud. Si lo perdiste, usa tu código de docente y tu cédula. No necesitas una cuenta.
        </p>
        <BuscadorSeguimiento codigoInicial={codigo} sobreOscuro />
        {!codigo && (
          <p className="text-[14px] text-white/80">
            ¿Aún no radicaste?{" "}
            <Link href="/solicitud" className="font-semibold text-white underline underline-offset-2 hover:no-underline">
              Inicia tu solicitud en línea
            </Link>
          </p>
        )}
      </section>

      {codigo && !vista && (
        <Alerta tipo="aviso" className="animate-materializar w-full text-center">
          No encontramos ningún trámite con el código <strong className="font-mono">{codigo}</strong>. Revísalo e inténtalo de nuevo.
        </Alerta>
      )}

      {vista && (
        <section className="animate-materializar material w-full rounded-[24px] p-6 shadow-[0_24px_60px_-24px_rgba(15,27,46,0.35)] sm:p-10">
          {nueva && (
            <Alerta tipo="exito" className="mb-8 text-center">
              <strong>¡Solicitud radicada!</strong> Guarda este código: <span className="font-mono font-bold">{vista.codigo}</span>. Un funcionario validará tus documentos en ventanilla.
            </Alerta>
          )}

          <div className="mb-9 flex flex-col gap-1.5 text-center">
            <p className="font-mono text-[13px] font-bold tracking-wider text-ink-faint">{vista.codigo}</p>
            <h2 className="font-display text-[clamp(1.5rem,3.5vw,2.25rem)] font-extrabold leading-[1.1] tracking-[-0.02em] text-ink">
              {vista.estado === "INGRESADO" && vista.pendienteValidacion ? "Pendiente de validación" : TITULO_ESTADO[vista.estado]}
            </h2>
            <p className="mx-auto max-w-md text-[15px] leading-relaxed text-ink-soft">
              {vista.estado === "INGRESADO" && vista.pendienteValidacion
                ? "Tu solicitud llegó por el portal. Falta que un funcionario valide tus documentos en ventanilla."
                : DESCRIPCION_ESTADO[vista.estado]}
            </p>
          </div>

          <FasesTramite estado={vista.estado} etapas={vista.etapas} />

          <dl className="mt-10 grid grid-cols-2 gap-5 border-t border-line pt-6 text-center sm:grid-cols-3">
            <Dato etiqueta="Tipo">{ETIQUETA_TIPO[vista.tipo]}</Dato>
            <Dato etiqueta="Categoría">{vista.categoria || "—"}</Dato>
            <Dato etiqueta="Radicado">{formatearFecha(vista.radicadoEn)}</Dato>
          </dl>

          {vista.resultado && (
            <p className="mt-8 rounded-[14px] bg-bg px-5 py-4 text-center text-[15px] text-ink">
              Resultado:{" "}
              <strong className={vista.resultado === "FAVORABLE" ? "text-success" : "text-accent-ink"}>
                {vista.resultado === "FAVORABLE" ? "Favorable" : "Desfavorable"}
              </strong>
            </p>
          )}

        </section>
      )}

      {vista && (
        <div className="flex flex-wrap items-start gap-6">
          {vista.observaciones.length > 0 && (
            <section className="flex min-w-0 flex-[999_1_360px] flex-col gap-3 rounded-[20px] border border-line bg-surface p-6">
              <h3 className="font-display text-[17px] font-extrabold text-ink">
                Observaciones {pendientes.length > 0 ? `(${pendientes.length} por subsanar)` : "(todas resueltas)"}
              </h3>
              {vista.observaciones.map((o, i) => (
                <div key={i} className={`rounded-[14px] border p-4 ${o.resuelta ? "border-line opacity-70" : "border-accent/40 bg-accent-soft"}`}>
                  <p className="text-[14px] font-bold text-ink">
                    {o.documento} {o.resuelta && <span className="font-semibold text-success">· resuelta</span>}
                  </p>
                  <p className="mt-1 text-[14px] text-ink-soft">{o.motivo}</p>
                  {!o.resuelta && <p className="mt-1 text-[12.5px] text-ink-faint">Plazo para subsanar: {formatearFecha(o.plazo)}</p>}
                </div>
              ))}
              {pendientes.length > 0 && <p className="text-[13px] text-ink-faint">Acércate a ventanilla con los documentos corregidos dentro del plazo.</p>}
            </section>
          )}

          {historial.length > 0 && (
            <section className="flex max-w-full flex-[1_1_260px] flex-col gap-4 rounded-[20px] border border-line bg-surface p-6">
              <h3 className="font-display text-[17px] font-extrabold text-ink">Historial</h3>
              <ol className="flex flex-col gap-4">
                {historial.map((e, i) => (
                  <li key={`${e.estado}-${e.fecha}`} className="grid grid-cols-[12px_1fr] items-start gap-3">
                    <span aria-hidden className={`mt-1.5 h-2.5 w-2.5 rounded-full ${i === 0 ? "bg-accent" : "bg-ink-faint"}`} />
                    <span className="flex flex-col">
                      <span className="text-[14.5px] font-semibold text-ink">{ETIQUETA_ESTADO[e.estado] ?? e.estado}</span>
                      <span className="text-[13px] text-ink-faint">{formatearFecha(e.fecha)}</span>
                    </span>
                  </li>
                ))}
              </ol>
            </section>
          )}
        </div>
      )}

      {vista && <p className="text-[13px] text-ink-faint">Por tu privacidad, aquí solo se muestra el estado del trámite: tus datos personales no aparecen.</p>}
    </div>
  );
}
function Dato({ etiqueta, children }: { etiqueta: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col gap-1">
      <dt className="text-[11.5px] font-bold uppercase tracking-wide text-ink-faint">{etiqueta}</dt>
      <dd className="text-[15px] font-semibold text-ink">{children}</dd>
    </div>
  );
}
