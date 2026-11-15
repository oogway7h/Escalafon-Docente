import "server-only";
import { unstable_rethrow } from "next/navigation";
import { api, enviarJson } from "./api";
import type { EstadoTramite } from "./tramites";

// Tablero de indicadores (CU-25, GET /api/tablero) y reportes exportables (CU-26, POST /api/reportes).

export type MotivoMora = "FECHA_LIMITE_VENCIDA" | "EXCEDE_MAXIMO";

export type Tablero = {
  generadoEn: string; // ISO 8601
  periodo: { desde: string | null; hasta: string | null };
  diasMaximosPrestamo: number;
  tramites: { radicados: number; enCurso: number; resueltos: number; promedioDias: number };
  porEstado: { estado: EstadoTramite; cantidad: number }[];
  tiempos: {
    tipo: "CATEGORIZACION" | "REGULARIZACION";
    categoria: string;
    cantidad: number;
    promedioDias: number;
    minimoDias: number;
    maximoDias: number;
  }[];
  analistas: { id: string; nombre: string; activos: number; enFacultad: number; resueltos: number }[];
  prestamosVigentes: number;
  alertas: {
    legajoCodigo: string;
    docenteNombre: string;
    custodio: string;
    tramiteCodigo: string;
    salida: string;
    fechaLimite: string | null; // AAAA-MM-DD
    diasPrestado: number;
    diasMora: number;
    motivo: MotivoMora;
  }[];
};

export type ResultadoTablero =
  | { ok: true; tablero: Tablero }
  | { ok: false; codigo?: string; detalles?: Record<string, string> };

// Igual que la auditoría: 401 y 403 los resuelve api(); un 422 o un fallo de red vuelven como error
// para que la pantalla lo explique sin perder el periodo elegido.
export async function obtenerTablero(periodo: { desde?: string; hasta?: string }): Promise<ResultadoTablero> {
  const qs = new URLSearchParams();
  for (const clave of ["desde", "hasta"] as const) {
    const valor = periodo[clave]?.trim();
    if (valor) qs.set(clave, valor);
  }
  try {
    const respuesta = await api(`/api/tablero?${qs}`);
    if (!respuesta.ok) {
      const cuerpo = await respuesta.json().catch(() => null);
      return { ok: false, codigo: cuerpo?.codigo, detalles: cuerpo?.detalles ?? undefined };
    }
    return { ok: true, tablero: await respuesta.json() };
  } catch (error) {
    unstable_rethrow(error);
    return { ok: false };
  }
}

export type SolicitudReporte = {
  tipo: string;
  formato: string;
  desde?: string;
  hasta?: string;
  estado?: string;
  tipoTramite?: string;
  categoriaId?: number;
  estadoLegajo?: string;
};

// Devuelve la respuesta de la API tal cual: el route handler la reenvía como descarga.
export function generarReporte(solicitud: SolicitudReporte): Promise<Response> {
  return enviarJson("/api/reportes", "POST", solicitud);
}

// Reportes en lenguaje natural: la API devuelve los filtros que entendió de la frase (sin generar nada).
export function interpretarReporte(frase: string, categorias: { id: number; nombre: string }[]): Promise<Response> {
  return enviarJson("/api/reportes/interpretar", "POST", { frase, categorias });
}
