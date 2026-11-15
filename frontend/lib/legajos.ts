// Sin "server-only": etiquetas de legajos y solicitudes de carpeta (CU-19 al CU-24).

export const ETIQUETA_SOLICITUD: Record<string, string> = {
  PENDIENTE: "Pendiente",
  ATENDIDA: "Aprobada",
  RECHAZADA: "Rechazada",
  DESPACHADA: "Prestada",
  DEVUELTA: "Devuelta",
};

export const TONO_SOLICITUD: Record<string, "acento" | "exito" | "neutro"> = {
  PENDIENTE: "acento",
  ATENDIDA: "acento",
  RECHAZADA: "neutro",
  DESPACHADA: "acento",
  DEVUELTA: "exito",
};

export const ESTADOS_MATERIAL = [
  { valor: "BUENO", etiqueta: "Buen estado" },
  { valor: "REGULAR", etiqueta: "Hojas sueltas" },
  { valor: "DETERIORADO", etiqueta: "Deteriorado" },
];

export function formatearFecha(iso: string | null): string {
  if (!iso) return "—";
  return new Date(iso.length === 10 ? `${iso}T12:00:00` : iso).toLocaleDateString("es-BO", { day: "2-digit", month: "short", year: "numeric" });
}

// Fecha AAAA-MM-DD a N días de hoy, para el valor inicial de <input type="date">.
export function enDias(n: number): string {
  const d = new Date();
  d.setDate(d.getDate() + n);
  return d.toISOString().slice(0, 10);
}
