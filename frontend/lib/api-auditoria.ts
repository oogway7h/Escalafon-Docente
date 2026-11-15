import "server-only";
import { unstable_rethrow } from "next/navigation";
import { api } from "./api";

// Consulta de la bitácora de auditoría (CU-06). Solo lectura: GET /api/auditoria y /api/auditoria/filtros.

export type AutorAuditoria = { id: string; nombres: string; apellidos: string; correo: string };

export type RegistroAuditoria = {
  id: number;
  fecha: string; // ISO 8601
  usuario: AutorAuditoria | null; // null: login con un correo inexistente o solicitud del portal
  accion: string;
  entidad: string;
  entidadId: string;
  detalle: Record<string, unknown> | null;
  direccionIp: string;
};

export type ListaAuditoria = { datos: RegistroAuditoria[]; total: number; pagina: number; tamano: number };

export type FiltrosAuditoria = { acciones: string[]; usuarios: AutorAuditoria[] };

// Mismos nombres que los parámetros de la API; desde/hasta en AAAA-MM-DD.
export type FiltroAuditoria = { usuarioId?: string; accion?: string; desde?: string; hasta?: string };

export type ResultadoAuditoria =
  | { ok: true; lista: ListaAuditoria }
  | { ok: false; codigo?: string; detalles?: Record<string, string> };

// Una página de la bitácora (10 por página, como el listado de docentes). 401 y 403 los resuelve
// api() (login / 404). Un 422 o un fallo de red se devuelven como error para que la pantalla lo
// explique sin perder los filtros.
export async function listarAuditoria(filtro: FiltroAuditoria, pagina: number, tamano = 10): Promise<ResultadoAuditoria> {
  const qs = new URLSearchParams({ pagina: String(pagina), tamano: String(tamano) });
  for (const clave of ["usuarioId", "accion", "desde", "hasta"] as const) {
    const valor = filtro[clave]?.trim();
    if (valor) qs.set(clave, valor);
  }
  try {
    const respuesta = await api(`/api/auditoria?${qs}`);
    if (!respuesta.ok) {
      const cuerpo = await respuesta.json().catch(() => null);
      return { ok: false, codigo: cuerpo?.codigo, detalles: cuerpo?.detalles ?? undefined };
    }
    return { ok: true, lista: await respuesta.json() };
  } catch (error) {
    unstable_rethrow(error);
    return { ok: false };
  }
}

export async function obtenerFiltrosAuditoria(): Promise<FiltrosAuditoria | null> {
  try {
    const respuesta = await api("/api/auditoria/filtros");
    return respuesta.ok ? respuesta.json() : null;
  } catch (error) {
    unstable_rethrow(error);
    return null;
  }
}
