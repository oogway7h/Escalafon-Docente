import Link from "next/link";
import Insignia from "@/components/common/Insignia";
import { IndicadorFila } from "@/components/common/EnlaceCarga";
import Tabla from "@/components/common/Tabla";
import type { TramiteResumen } from "@/lib/api-tramites";
import { formatearFecha, tipoYCategoria } from "@/lib/tramites";
import EstadoInsignia from "./EstadoInsignia";

// Listado de trámites reutilizado por "Todos", "Asignar" y "Mi bandeja".
export default function TablaTramites({
  tramites,
  vacio,
  accion,
}: {
  tramites: TramiteResumen[];
  vacio?: string;
  accion?: (t: TramiteResumen) => React.ReactNode;
}) {
  return (
    <Tabla
      filas={tramites}
      clave={(t) => t.id}
      href={(t) => `/panel/tramites/${t.id}`}
      vacio={vacio ?? "No hay trámites."}
      columnas={[
        {
          titulo: "Trámite",
          celda: (t) => (
            <>
              <Link href={`/panel/tramites/${t.id}`} className="font-mono text-sm font-semibold text-ink hover:text-accent-ink hover:underline">
                {t.codigo}
              </Link>
              <IndicadorFila />
              <div className="text-xs text-ink-faint">
                {tipoYCategoria(t)}
              </div>
            </>
          ),
        },
        {
          titulo: "Docente",
          celda: (t) => (
            <>
              <div className="font-semibold">{t.docenteNombre}</div>
              <div className="text-xs text-ink-faint">Cód. {t.docenteCodigo}</div>
            </>
          ),
        },
        {
          titulo: "Estado",
          celda: (t) => (
            <div className="flex flex-wrap items-center gap-1.5">
              <EstadoInsignia estado={t.estado} />
              {t.pendienteValidacion && <Insignia tono="acento">Por validar</Insignia>}
              {t.origen === "PORTAL" && <Insignia tono="neutro">Portal</Insignia>}
            </div>
          ),
        },
        { titulo: "Analista", celda: (t) => t.analistaNombre || <span className="text-ink-faint">Sin asignar</span> },
        { titulo: "Radicado", celda: (t) => <span className="text-sm text-ink-soft">{formatearFecha(t.creadoEn)}</span> },
        ...(accion ? [{ className: "w-56", celda: accion }] : []),
      ]}
    />
  );
}
