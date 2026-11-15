// Sin "server-only": etiquetas y formatos del tablero (CU-25) y de los reportes (CU-26).
import { mensajeError } from "./formularios";

export const TIPOS_REPORTE = [
  {
    valor: "TRAMITES",
    etiqueta: "Trámites",
    descripcion: "Trámites radicados con su estado actual, tipo, categoría, analista y resultado.",
  },
  {
    valor: "INVENTARIO",
    etiqueta: "Inventario de legajos",
    descripcion: "Carpetas físicas con su ubicación, el préstamo vigente y los días de mora.",
  },
  {
    valor: "TIEMPOS",
    etiqueta: "Tiempos de atención",
    descripcion: "Días entre la radicación y la remisión de cada trámite resuelto, con promedios por categoría.",
  },
] as const;

export type TipoReporte = (typeof TIPOS_REPORTE)[number]["valor"];

// Error de un reporte que no se pudo generar (REPORTE_VACIO es el caso 3a de CU-26).
export function mensajeReporte(codigo: string | undefined): string {
  if (codigo === "VALIDACION") return "Revisa los filtros: alguna fecha no tiene un formato válido.";
  return mensajeError(codigo, "No se pudo generar el reporte. Intenta de nuevo.");
}

export const ETIQUETA_MOTIVO_MORA = {
  FECHA_LIMITE_VENCIDA: "Pasó la fecha límite",
  EXCEDE_MAXIMO: "Supera el máximo de días",
} as const;

// 12.46 -> "12,5"
export function formatearDias(dias: number): string {
  return dias.toLocaleString("es-BO", { minimumFractionDigits: 1, maximumFractionDigits: 1 });
}

export function plural(n: number, uno: string, varios: string): string {
  return `${n.toLocaleString("es-BO")} ${n === 1 ? uno : varios}`;
}

// Hora de Bolivia de la última actualización del tablero.
export function formatearHora(iso: string): string {
  return new Date(iso).toLocaleTimeString("es-BO", { timeZone: "America/La_Paz", hour: "2-digit", minute: "2-digit" });
}

// Atajos de periodo: fechas AAAA-MM-DD en hora de Bolivia.
export function periodosRapidos(ahora = new Date()): { etiqueta: string; desde: string; hasta: string }[] {
  const hoy = new Date(ahora.toLocaleString("en-US", { timeZone: "America/La_Paz" }));
  const iso = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
  const hace = (meses: number) => new Date(hoy.getFullYear(), hoy.getMonth() - meses, hoy.getDate());
  return [
    { etiqueta: "Este mes", desde: iso(new Date(hoy.getFullYear(), hoy.getMonth(), 1)), hasta: iso(hoy) },
    { etiqueta: "Últimos 3 meses", desde: iso(hace(3)), hasta: iso(hoy) },
    { etiqueta: "Este año", desde: iso(new Date(hoy.getFullYear(), 0, 1)), hasta: iso(hoy) },
  ];
}
