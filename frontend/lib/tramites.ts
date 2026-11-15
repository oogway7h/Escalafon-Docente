// Sin "server-only": lo usan Server y Client Components.
export type EstadoTramite =
  | "INGRESADO"
  | "ASIGNADO"
  | "EN_REVISION"
  | "OBSERVADO"
  | "DICTAMINADO"
  | "EN_FACULTAD"
  | "REVISION_FINAL"
  | "REMITIDO"
  | "ARCHIVADO";

export const FASES: EstadoTramite[] = [
  "INGRESADO",
  "ASIGNADO",
  "EN_REVISION",
  "OBSERVADO",
  "DICTAMINADO",
  "EN_FACULTAD",
  "REVISION_FINAL",
  "REMITIDO",
  "ARCHIVADO",
];

export const ETIQUETA_ESTADO: Record<EstadoTramite, string> = {
  INGRESADO: "Ingresado",
  ASIGNADO: "Asignado",
  EN_REVISION: "En revisión",
  OBSERVADO: "Observado",
  DICTAMINADO: "Dictaminado",
  EN_FACULTAD: "Fuera de Escalafón",
  REVISION_FINAL: "Revisión final",
  REMITIDO: "Remitido a Rectorado",
  ARCHIVADO: "Archivado",
};

export const DESCRIPCION_ESTADO: Record<EstadoTramite, string> = {
  INGRESADO: "Tu solicitud fue recibida y espera ser asignada a un analista.",
  ASIGNADO: "Un analista tiene asignado tu expediente y lo revisará pronto.",
  EN_REVISION: "El analista está revisando tu documentación.",
  OBSERVADO: "Se encontraron observaciones en tu documentación. Revísalas abajo.",
  DICTAMINADO: "El analista terminó la revisión de tu expediente y lo remitirá a la facultad.",
  EN_FACULTAD:
    "Tu expediente está en la facultad (Tribunal, Decano y Consejo Facultativo). En esta etapa no es posible seguir su avance; te avisaremos cuando regrese a Escalafón con la resolución.",
  REVISION_FINAL: "Tu expediente regresó con la resolución y el analista la está revisando.",
  REMITIDO: "Tu expediente fue remitido al Rectorado. Falta archivar la carpeta.",
  ARCHIVADO: "Tu trámite concluyó y su legajo fue archivado.",
};

// Tipos de escaneo del expediente digital del docente (los sube el Personal de Archivo).
export const TIPOS_DOCUMENTO = [
  { valor: "CURRICULUM", etiqueta: "Currículum" },
  { valor: "REGULARIZACION", etiqueta: "Regularización" },
  { valor: "ADJUNTO", etiqueta: "Adjunto" },
  { valor: "ASOCIADO", etiqueta: "Asociado" },
  { valor: "TITULAR_A", etiqueta: "Titular A" },
  { valor: "TITULAR_B", etiqueta: "Titular B" },
  { valor: "TITULAR_C", etiqueta: "Titular C" },
  { valor: "EMERITO", etiqueta: "Emérito" },
] as const;

export const ETIQUETA_TIPO_DOCUMENTO: Record<string, string> = Object.fromEntries(TIPOS_DOCUMENTO.map((t) => [t.valor, t.etiqueta]));

// Los códigos de docente tienen hasta 4 dígitos (se admiten 5); lo que supere este largo es un código de seguimiento.
export const MAX_CODIGO_DOCENTE = 5;

// Título corto de cada estado para el portal público.
export const TITULO_ESTADO = {
  INGRESADO: "Solicitud recibida",
  ASIGNADO: "Asignado a un analista",
  EN_REVISION: "En revisión",
  OBSERVADO: "Requiere tu atención",
  DICTAMINADO: "Revisión terminada",
  EN_FACULTAD: "En la facultad",
  REVISION_FINAL: "Revisando la resolución",
  REMITIDO: "Remitido al Rectorado",
  ARCHIVADO: "Trámite concluido",
} as const;

export const ETIQUETA_TIPO = { CATEGORIZACION: "Categorización", REGULARIZACION: "Regularización" } as const;

// "Categorización · Titular A"; una regularización no tiene categoría y queda solo con el tipo.
export function tipoYCategoria(t: { tipo: keyof typeof ETIQUETA_TIPO; categoria?: string }): string {
  return t.categoria ? `${ETIQUETA_TIPO[t.tipo]} · ${t.categoria}` : ETIQUETA_TIPO[t.tipo];
}

// Tamaños máximos por archivo (los mismos que valida la API): un solo lugar para textos y validaciones.
export const LIMITE_MB = { portal: 8, pdf: 400, imagen: 20 } as const;
// Adjuntos libres de una regularización radicada desde el portal.
export const MAX_ADJUNTOS_LIBRES = 5;

export const TONO_ESTADO: Record<EstadoTramite, "acento" | "exito" | "neutro"> = {
  INGRESADO: "neutro",
  ASIGNADO: "neutro",
  EN_REVISION: "acento",
  OBSERVADO: "acento",
  DICTAMINADO: "acento",
  EN_FACULTAD: "neutro",
  REVISION_FINAL: "acento",
  REMITIDO: "exito",
  ARCHIVADO: "exito",
};

export function formatearTamano(bytes: number): string {
  if (bytes >= 1024 ** 3) return `${(bytes / 1024 ** 3).toFixed(1)} GB`;
  if (bytes >= 1024 ** 2) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${Math.max(1, Math.round(bytes / 1024))} KB`;
}

// Días completos entre dos fechas ISO; sin "hasta" cuenta hasta hoy (p. ej. días fuera de Escalafón).
export function diasEntre(desde: string, hasta?: string): number {
  const fin = hasta ? new Date(hasta).getTime() : Date.now();
  return Math.max(0, Math.floor((fin - new Date(desde).getTime()) / 86_400_000));
}

export function formatearFecha(iso: string, conHora = false): string {
  return new Date(iso).toLocaleString("es-BO", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    ...(conHora ? { hour: "2-digit", minute: "2-digit" } : {}),
  });
}

// Mañana en AAAA-MM-DD: plazo mínimo para subsanar una observación.
export function fechaManana(): string {
  return new Date(Date.now() + 86_400_000).toISOString().slice(0, 10);
}
