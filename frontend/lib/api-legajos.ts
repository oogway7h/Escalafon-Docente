import "server-only";
import { api, ApiError } from "./api";

// Legajos, solicitudes de carpeta y custodia (CU-19 al CU-24).

export type Carpeta = {
  id: number;
  codigo: string;
  docenteId: string;
  docenteCodigo: string;
  docenteNombre: string;
  gaveta: string;
  orden: number;
  anexos: string;
  observacion: string;
  estado: "EN_ARCHIVO" | "PRESTADO";
};

export type EstadoSolicitud = "PENDIENTE" | "ATENDIDA" | "RECHAZADA" | "DESPACHADA" | "DEVUELTA";

export type SolicitudCarpeta = {
  id: string;
  legajoId: number;
  legajoCodigo: string;
  docenteNombre: string;
  tramiteId: string;
  tramiteCodigo: string;
  solicitante: string;
  motivo: string;
  estado: EstadoSolicitud;
  respuesta: string;
  creadoEn: string;
  atendidoEn: string | null;
  fechaLimite: string | null;
};

export type Movimiento = {
  id: number;
  legajoCodigo: string;
  docenteNombre: string;
  tipo: "SALIDA" | "ENTRADA";
  custodio: string;
  fechaLimite: string | null;
  estadoMaterial: "BUENO" | "REGULAR" | "DETERIORADO";
  observacion: string;
  registradoPor: string;
  fecha: string;
  vencido: boolean;
};

export type Pagina<T> = { datos: T[]; total: number; pagina: number; tamano: number };

async function pagina<T>(ruta: string, filtro: Record<string, string | undefined>, nPagina: number, tamano: number, mensaje: string): Promise<Pagina<T>> {
  const qs = new URLSearchParams({ pagina: String(nPagina), tamano: String(tamano) });
  for (const [k, v] of Object.entries(filtro)) if (v) qs.set(k, v);
  const respuesta = await api(`${ruta}?${qs}`);
  if (!respuesta.ok) throw new ApiError(respuesta.status, "ERROR_LEGAJOS", mensaje);
  return respuesta.json();
}

export const listarCarpetas = (filtro: { buscar?: string; estado?: string }, nPagina: number, tamano = 10) =>
  pagina<Carpeta>("/api/legajos", filtro, nPagina, tamano, "No se pudo cargar el inventario de legajos");

export async function obtenerCarpeta(id: string): Promise<Carpeta> {
  const respuesta = await api(`/api/legajos/${id}`);
  if (!respuesta.ok) throw new ApiError(respuesta.status, "ERROR_LEGAJO", "No se pudo cargar el legajo");
  return respuesta.json();
}

export const listarSolicitudes = (filtro: { estado?: string }, nPagina: number, tamano = 10) =>
  pagina<SolicitudCarpeta>("/api/solicitudes-carpeta", filtro, nPagina, tamano, "No se pudo cargar las solicitudes de carpeta");

export const listarBitacora = (filtro: { buscar?: string; tipo?: string; vencidos?: string }, nPagina: number, tamano = 15) =>
  pagina<Movimiento>("/api/custodia/bitacora", filtro, nPagina, tamano, "No se pudo cargar la bitácora de custodia");

export type Notificacion = {
  id: string;
  tramiteCodigo: string;
  estadoTramite: "OBSERVADO" | "EN_FACULTAD" | "REMITIDO" | "ARCHIVADO";
  destino: string;
  estado: "PENDIENTE" | "ENVIADO" | "FALLIDO";
  intentos: number;
  error: string;
  creadoEn: string;
  enviadoEn: string | null;
};

// CU-14: avisos de cambio de estado enviados a los docentes.
export const listarNotificaciones = (filtro: { estado?: string }, nPagina: number, tamano = 15) =>
  pagina<Notificacion>("/api/notificaciones", filtro, nPagina, tamano, "No se pudo cargar las notificaciones");
