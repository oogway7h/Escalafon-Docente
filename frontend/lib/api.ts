import "server-only";
import { cache } from "react";
import { headers } from "next/headers";
import { notFound, redirect } from "next/navigation";
import { borrarTokenSesion, leerTokenSesion } from "./sesion";

export class ApiError extends Error {
  constructor(
    public status: number,
    public codigo: string,
    message: string,
  ) {
    super(message);
  }
}

// Reenvía a la API la IP del navegador: sin esto todas las peticiones salen con la IP del servidor de Next y el
// limitador de tasa y la auditoría ven a todos los usuarios como uno solo. Se toma la última entrada de
// x-forwarded-for (la que agregó el proxy más cercano); en producción el proxy inverso debe sobrescribir la cabecera.
export async function cabecerasCliente(): Promise<Record<string, string>> {
  const ip = (await headers()).get("x-forwarded-for")?.split(",").pop()?.trim();
  return ip ? { "X-Forwarded-For": ip } : {};
}

export async function api(path: string, init?: RequestInit): Promise<Response> {
  const apiUrl = process.env.API_URL_INTERNA;
  if (!apiUrl) throw new Error("API_URL_INTERNA no está configurada");

  const token = await leerTokenSesion();
  const respuesta = await fetch(`${apiUrl}${path}`, {
    ...init,
    headers: {
      ...(init?.headers ?? {}),
      ...(await cabecerasCliente()),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    cache: "no-store",
  });

  if (respuesta.status === 401) {
    // Next solo permite borrar cookies en acciones de servidor y route handlers; durante el render
    // de una página lanzaría un error (500). Ahí basta con redirigir: el nuevo login la reemplaza.
    await borrarTokenSesion().catch(() => {});
    redirect("/login");
  }
  if (respuesta.status === 403) {
    notFound();
  }
  return respuesta;
}

export function enviarJson(path: string, method: string, datos?: unknown): Promise<Response> {
  return api(path, {
    method,
    headers: { "Content-Type": "application/json" },
    body: datos === undefined ? undefined : JSON.stringify(datos),
  });
}

export type UsuarioResumen = {
  id: string;
  nombres: string;
  apellidos: string;
  correo: string;
  rolId: number;
};

export type Sesion = {
  usuario: UsuarioResumen;
  rol: string;
  permisos: string[];
  terminosAceptados: boolean; // false: se muestra el acuerdo de usuario antes del panel
};

// cache() deduplica la llamada dentro de la misma petición: el layout del panel y cada página
// pueden pedir la sesión sin duplicar la ida y vuelta a la API.
export const obtenerSesion = cache(async (): Promise<Sesion> => {
  const respuesta = await api("/api/auth/sesion");
  if (!respuesta.ok) {
    throw new ApiError(respuesta.status, "ERROR_SESION", "No se pudo obtener la sesión");
  }
  return respuesta.json();
});

// Devuelve la sesión o responde 404 si el rol no tiene el permiso.
export async function exigirPermiso(permiso: string): Promise<Sesion> {
  const sesion = await obtenerSesion();
  if (!sesion.permisos.includes(permiso)) notFound();
  return sesion;
}

export async function cerrarSesionApi(): Promise<void> {
  await api("/api/auth/logout", { method: "POST" });
}

export type Perfil = {
  id: string;
  nombres: string;
  apellidos: string;
  correo: string;
  telefono: string;
  rolId: number;
};

export const obtenerPerfil = cache(async (): Promise<Perfil> => {
  const respuesta = await api("/api/perfil");
  if (!respuesta.ok) throw new ApiError(respuesta.status, "ERROR_PERFIL", "No se pudo obtener el perfil");
  return respuesta.json();
});

export type UsuarioListado = {
  id: string;
  nombres: string;
  apellidos: string;
  correo: string;
  telefono: string;
  rolId: number;
  activo: boolean;
};

export type ListaUsuarios = { datos: UsuarioListado[]; total: number; pagina: number; tamano: number };

export type RolConConteo = {
  id: number;
  nombre: string;
  descripcion: string;
  esSistema: boolean;
  activo: boolean;
  usuariosCount: number;
};

export type Permiso = { id: number; codigo: string; nombre: string; modulo: string; casoUso: string };

export type RolConPermisos = {
  id: number;
  nombre: string;
  descripcion: string;
  esSistema: boolean;
  activo: boolean;
  permisos: Permiso[];
};

export type FiltroUsuarios = { buscar?: string; rolId?: string; activo?: string; pagina?: string };

export async function listarUsuarios(filtro: FiltroUsuarios): Promise<ListaUsuarios> {
  const qs = new URLSearchParams();
  if (filtro.buscar) qs.set("buscar", filtro.buscar);
  if (filtro.rolId) qs.set("rolId", filtro.rolId);
  if (filtro.activo) qs.set("activo", filtro.activo);
  if (filtro.pagina) qs.set("pagina", filtro.pagina);
  const respuesta = await api(`/api/usuarios?${qs.toString()}`);
  if (!respuesta.ok) throw new ApiError(respuesta.status, "ERROR_USUARIOS", "No se pudo cargar la lista de usuarios");
  return respuesta.json();
}

export async function obtenerUsuario(id: string): Promise<UsuarioListado> {
  const respuesta = await api(`/api/usuarios/${id}`);
  if (!respuesta.ok) throw new ApiError(respuesta.status, "ERROR_USUARIO", "No se pudo cargar el usuario");
  return respuesta.json();
}

export async function listarRoles(): Promise<RolConConteo[]> {
  const respuesta = await api("/api/roles");
  if (!respuesta.ok) throw new ApiError(respuesta.status, "ERROR_ROLES", "No se pudo cargar la lista de roles");
  return respuesta.json();
}

export async function obtenerRol(id: string): Promise<RolConPermisos> {
  const respuesta = await api(`/api/roles/${id}`);
  if (!respuesta.ok) throw new ApiError(respuesta.status, "ERROR_ROL", "No se pudo cargar el rol");
  return respuesta.json();
}

export async function listarPermisos(): Promise<Permiso[]> {
  const respuesta = await api("/api/permisos");
  if (!respuesta.ok) throw new ApiError(respuesta.status, "ERROR_PERMISOS", "No se pudo cargar el catálogo de permisos");
  return respuesta.json();
}

// Versión de la foto de perfil del usuario (cambia con cada foto nueva), o null si no tiene o no se pudo saber.
// Nunca lanza ni redirige: el panel debe abrir igual aunque la API o la migración de la foto fallen.
export const obtenerVersionFoto = cache(async (): Promise<number | null> => {
  try {
    const apiUrl = process.env.API_URL_INTERNA;
    const token = await leerTokenSesion();
    if (!apiUrl || !token) return null;
    const r = await fetch(`${apiUrl}/api/perfil/foto`, {
      headers: { Authorization: `Bearer ${token}`, ...(await cabecerasCliente()) },
      cache: "no-store",
      signal: AbortSignal.timeout(3000),
    });
    return r.ok ? ((await r.json()) as { version: number }).version : null;
  } catch {
    return null;
  }
});
