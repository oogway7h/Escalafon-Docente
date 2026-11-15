import type { EstadoTramite } from "@/lib/tramites";
import { diasEntre, ETIQUETA_ESTADO, formatearFecha } from "@/lib/tramites";

// Recorrido real del expediente. "Observado" es un desvío de "En revisión", no una fase más.
// EN_FACULTAD no es un punto: es el tramo en que la carpeta está fuera del Departamento (facultad,
// tribunal, decano y consejo) y el sistema no puede rastrearla; se dibuja punteado entre medias.
const ORDEN: EstadoTramite[] = ["INGRESADO", "ASIGNADO", "EN_REVISION", "DICTAMINADO", "EN_FACULTAD", "REVISION_FINAL", "REMITIDO", "ARCHIVADO"];

export default function FasesTramite({
  estado,
  etapas,
}: {
  estado: EstadoTramite;
  etapas: { estado: EstadoTramite; fecha: string }[];
}) {
  const observado = estado === "OBSERVADO";
  const actual = ORDEN.indexOf(observado ? "EN_REVISION" : estado);
  const fechaDe = (e: EstadoTramite) => etapas.findLast((x) => x.estado === e)?.fecha;
  const llena = (i: number) => actual >= i; // el tramo que llega a la posición i ya se recorrió

  return (
    <div className="overflow-x-auto pb-2">
      <ol className="flex min-w-[720px] items-start" aria-label="Fases del trámite">
        {ORDEN.map((fase, i) => {
          if (fase === "EN_FACULTAD") {
            const salida = fechaDe("EN_FACULTAD");
            const regreso = fechaDe("REVISION_FINAL");
            const dentro = actual === i;
            return (
              <li key={fase} className="relative flex min-w-36 flex-[2.4] flex-col items-center px-2 pt-9 text-center" aria-current={dentro ? "step" : undefined}>
                <span aria-hidden className={`absolute left-0 right-0 top-[13px] border-t-2 border-dashed ${llena(i) ? "border-accent" : "border-line"}`} />
                <span
                  className={`absolute left-1/2 top-[13px] -translate-x-1/2 -translate-y-1/2 whitespace-nowrap rounded-full border bg-surface px-2.5 py-0.5 text-xs font-semibold uppercase tracking-wide ${
                    dentro ? "animate-pulse border-accent text-accent-ink" : llena(i) ? "border-accent/50 text-ink-soft" : "border-line text-ink-faint"
                  }`}
                >
                  Sin seguimiento
                </span>
                <span className={`text-xs font-semibold leading-tight ${dentro || llena(i) ? "text-ink" : "text-ink-faint"}`}>{ETIQUETA_ESTADO[fase]}</span>
                <span className="text-xs leading-tight text-ink-faint">Facultad · Tribunal · Decano · Consejo</span>
                {salida && (
                  <span className="mt-1 text-xs text-ink-faint">
                    Salió {formatearFecha(salida)}
                    {regreso ? ` · regresó ${formatearFecha(regreso)}` : ""}
                    <span className="block font-semibold text-ink-soft">{diasEntre(salida, regreso)} días fuera</span>
                  </span>
                )}
              </li>
            );
          }
          const hecha = i < actual || (i === actual && estado === "ARCHIVADO");
          const enCurso = i === actual && !hecha;
          const fecha = fechaDe(fase);
          return (
            <li key={fase} className="contents">
              {i > 0 && ORDEN[i - 1] !== "EN_FACULTAD" && (
                <span aria-hidden className={`mt-[13px] h-0.5 min-w-3 flex-1 rounded-full ${llena(i) ? "bg-accent" : "bg-line"}`} />
              )}
              {i > 0 && ORDEN[i - 1] === "EN_FACULTAD" && <span aria-hidden className="hidden" />}
              <div className="relative z-10 flex w-[72px] shrink-0 flex-col items-center gap-2 text-center" aria-current={enCurso ? "step" : undefined}>
                <span
                  className={`flex h-7 w-7 items-center justify-center rounded-full border-2 text-xs font-semibold transition-colors ${
                    hecha
                      ? "border-accent bg-accent text-white"
                      : enCurso
                        ? "border-accent bg-surface text-accent-ink shadow-[0_0_0_4px_var(--color-accent-soft)]"
                        : "border-line bg-surface text-ink-faint"
                  }`}
                >
                  {hecha ? "✓" : i + 1 - (i > ORDEN.indexOf("EN_FACULTAD") ? 1 : 0)}
                </span>
                <span className={`text-xs leading-tight ${enCurso || hecha ? "font-semibold text-ink" : "text-ink-faint"}`}>
                  {ETIQUETA_ESTADO[fase]}
                  {enCurso && observado && <span className="block text-accent-ink">con observaciones</span>}
                </span>
                {fecha && (hecha || enCurso) && <span className="text-xs text-ink-faint">{formatearFecha(fecha)}</span>}
              </div>
            </li>
          );
        })}
      </ol>
    </div>
  );
}
