import "server-only";
import { connection } from "next/server";
import { api, ApiError, cabecerasCliente } from "./api";
import type { EstadoTramite } from "./tramites";

// Docentes, categorías, trámites (CU-07 al CU-11) y portal público (CU-12, CU-13).

export type Requisito = { id: number; nombre: string; presente?: boolean };

export type Categoria = {
  id: number;
  nombre: string;
  aniosAntiguedad: number;
  puntajeDesde: number | null; // rango de puntos de categoría (art. 18 del Reglamento); null = sin puntaje
  puntajeHasta: number | null; // null con puntajeDesde = «o más»
  activa: boolean;
  requisitos: Requisito[];
  docentes: number; // docentes con esta categoría
  tramites: number; // trámites que la postulan
};

export type Legajo = { gaveta: string; orden: number; anexos: string; observacion: string };

export type Facultad = { id: number; nombre: string };

export type Docente = {
  id: string;
  codigo: string;
  cedula: string;
  nombres: string;
  apellidos: string;
  correo: string;
  telefono: string;
  categoriaId: number | null;
  categoria: string;
  facultadId: number | null;
  facultad: string;
  legajos: Legajo[];
  activo: boolean;
};

export type TramiteResumen = {
  id: string;
  codigo: string;
  docenteNombre: string;
  docenteCodigo: string;
  categoria: string;
  tipo: "CATEGORIZACION" | "REGULARIZACION";
  estado: EstadoTramite;
  origen: "STAFF" | "PORTAL";
  pendienteValidacion: boolean;
  analistaNombre: string;
  creadoEn: string;
};

export type Observacion = {
  id: string;
  documento: string;
  motivo: string;
  plazo: string;
  resuelta: boolean;
  resueltaEn: string | null;
};

// Adjunto del expediente. activo=false es una baja lógica (CU-16): conserva el archivo y el motivo.
export type DocumentoTramite = {
  id: string;
  tipoDocumento: string; // solo en los escaneos del expediente del docente ("" en adjuntos de trámite)
  nombreOriginal: string;
  tipoMime: string;
  tamano: number;
  requisitoId: number | null;
  activo: boolean;
  motivoBaja: string;
  bajaEn: string | null;
  // CU-15: NO_APLICA (portal) | SUBIENDO | PENDIENTE | PROCESANDO | PROCESADO | ERROR
  estadoOcr: "NO_APLICA" | "SUBIENDO" | "PENDIENTE" | "PROCESANDO" | "PROCESADO" | "ERROR";
  errorOcr: string;
  textoValidado: boolean;
  hashSha256: string;
  tamanoComprimido: number; // 0 = sin versión comprimida
};

export type TramiteDetalle = TramiteResumen & {
  docenteId: string;
  categoriaId: number;
  analistaId: string | null;
  resultado: "" | "FAVORABLE" | "DESFAVORABLE";
  dictamen: string;
  dictaminadoEn: string | null;
  archivadoEn: string | null;
  ubicacionLegajo: string;
  // Etapas posteriores al dictamen
  anioDesde: number; // 0 = aún sin registrar
  anioHasta: number;
  resolucionNumero: string;
  resolucionFecha: string | null;
  remisionNota: string;
  carpetasLigadas: number; // solicitudes de carpeta no rechazadas
  carpetasAbiertas: number; // las que aún no volvieron al archivo
  historial: { estado: EstadoTramite; estadoAnterior: EstadoTramite | ""; usuario: string; detalle: string; fecha: string }[];
  observaciones: Observacion[];
  documentos: DocumentoTramite[];
  requisitos: Requisito[];
};

export type Analista = { id: string; nombre: string; carga: number };

export type FiltroTramites = { estado?: string; buscar?: string; mios?: boolean; sinAsignar?: boolean };

export type VistaPublica = {
  codigo: string;
  tipo: "CATEGORIZACION" | "REGULARIZACION";
  categoria: string;
  estado: EstadoTramite;
  radicadoEn: string;
  pendienteValidacion: boolean;
  resultado: "" | "FAVORABLE" | "DESFAVORABLE";
  etapas: { estado: EstadoTramite; fecha: string }[];
  observaciones: { documento: string; motivo: string; plazo: string; resuelta: boolean }[];
};

async function leerJson<T>(respuesta: Response, codigo: string, mensaje: string): Promise<T> {
  if (!respuesta.ok) throw new ApiError(respuesta.status, codigo, mensaje);
  return respuesta.json();
}

export async function listarDocentes(buscar?: string): Promise<Docente[]> {
  const qs = buscar ? `?buscar=${encodeURIComponent(buscar)}` : "";
  return leerJson(await api(`/api/docentes${qs}`), "ERROR_DOCENTES", "No se pudo cargar el padrón de docentes");
}

export type ListaDocentes = { datos: Docente[]; total: number; pagina: number; tamano: number };

// Una página del padrón (10 por defecto), con búsqueda por código, cédula o nombre.
export async function listarDocentesPagina(buscar: string | undefined, pagina: number, tamano = 10): Promise<ListaDocentes> {
  const qs = new URLSearchParams({ pagina: String(pagina), tamano: String(tamano) });
  if (buscar) qs.set("buscar", buscar);
  return leerJson(await api(`/api/docentes?${qs}`), "ERROR_DOCENTES", "No se pudo cargar el padrón de docentes");
}

export async function obtenerDocente(id: string): Promise<Docente> {
  return leerJson(await api(`/api/docentes/${id}`), "ERROR_DOCENTE", "No se pudo cargar el docente");
}

export async function listarFacultades(): Promise<Facultad[]> {
  return leerJson(await api("/api/facultades"), "ERROR_FACULTADES", "No se pudo cargar las facultades");
}

export async function listarCategorias(): Promise<Categoria[]> {
  return leerJson(await api("/api/categorias"), "ERROR_CATEGORIAS", "No se pudo cargar el catálogo de categorías");
}

export async function obtenerCategoria(id: string): Promise<Categoria> {
  return leerJson(await api(`/api/categorias/${id}`), "ERROR_CATEGORIA", "No se pudo cargar la categoría");
}

export async function listarTramites(filtro: FiltroTramites = {}): Promise<TramiteResumen[]> {
  const qs = new URLSearchParams();
  if (filtro.estado) qs.set("estado", filtro.estado);
  if (filtro.buscar) qs.set("buscar", filtro.buscar);
  if (filtro.mios) qs.set("mios", "true");
  if (filtro.sinAsignar) qs.set("sinAsignar", "true");
  return leerJson(await api(`/api/tramites?${qs}`), "ERROR_TRAMITES", "No se pudo cargar la lista de trámites");
}

export type ListaTramites = { datos: TramiteResumen[]; total: number; pagina: number; tamano: number };

// Una página de trámites (10 por defecto) para el listado general.
export async function listarTramitesPagina(filtro: FiltroTramites, pagina: number, tamano = 10): Promise<ListaTramites> {
  const qs = new URLSearchParams({ pagina: String(pagina), tamano: String(tamano) });
  if (filtro.estado) qs.set("estado", filtro.estado);
  if (filtro.buscar) qs.set("buscar", filtro.buscar);
  return leerJson(await api(`/api/tramites?${qs}`), "ERROR_TRAMITES", "No se pudo cargar la lista de trámites");
}

export async function obtenerTramite(id: string): Promise<TramiteDetalle> {
  return leerJson(await api(`/api/tramites/${id}`), "ERROR_TRAMITE", "No se pudo cargar el trámite");
}

// --- Expediente digital del docente ---

export type CoincidenciaExpediente = {
  documentoId: string;
  tipoDocumento: string;
  nombreOriginal: string;
  pagina: number;
  fragmento: string; // las palabras halladas van entre ⟦ y ⟧
  rango: number;
};

export async function listarExpediente(docenteId: string): Promise<DocumentoTramite[]> {
  return leerJson(await api(`/api/docentes/${docenteId}/documentos`), "ERROR_EXPEDIENTE", "No se pudo cargar el expediente del docente");
}

// Busca por contenido en todos los documentos procesados del docente; devuelve el código de error de la API si falla.
export async function buscarEnExpediente(docenteId: string, q: string): Promise<{ coincidencias: CoincidenciaExpediente[]; error?: string }> {
  const respuesta = await api(`/api/docentes/${docenteId}/documentos/buscar?q=${encodeURIComponent(q)}`);
  if (!respuesta.ok) {
    const cuerpo = await respuesta.json().catch(() => null);
    return { coincidencias: [], error: cuerpo?.codigo ?? "ERROR" };
  }
  return { coincidencias: (await respuesta.json()).coincidencias };
}

// --- CU-17: documentos dados de baja ---

export type Pagina<T> = { datos: T[]; total: number; pagina: number; tamano: number };

export type DocumentoDeBaja = DocumentoTramite & {
  tramiteId: string;
  tramiteCodigo: string;
  tramiteEstado: EstadoTramite | "";
  docenteId: string;
  docenteNombre: string;
  docenteCodigo: string;
  bajaPorNombre: string;
};

export async function listarDocumentosDeBaja(buscar: string, pagina: number): Promise<Pagina<DocumentoDeBaja>> {
  const qs = new URLSearchParams({ pagina: String(pagina), tamano: "10" });
  if (buscar) qs.set("buscar", buscar);
  return leerJson(await api(`/api/documentos/bajas?${qs}`), "ERROR_BAJAS", "No se pudieron cargar los documentos dados de baja");
}

export async function listarAnalistas(): Promise<Analista[]> {
  return leerJson(await api("/api/analistas"), "ERROR_ANALISTAS", "No se pudo cargar la lista de analistas");
}

// --- Portal público: sin sesión ni redirecciones ---

export async function apiPublica(path: string, init?: RequestInit): Promise<Response> {
  const apiUrl = process.env.API_URL_INTERNA;
  if (!apiUrl) throw new Error("API_URL_INTERNA no está configurada");
  return fetch(`${apiUrl}${path}`, { ...init, headers: { ...(init?.headers ?? {}), ...(await cabecerasCliente()) }, cache: "no-store" });
}

// null = el código no existe.
export async function consultarSeguimiento(codigo: string): Promise<VistaPublica | null> {
  const respuesta = await apiPublica(`/api/publico/tramites/${encodeURIComponent(codigo.trim())}`);
  if (respuesta.status === 404) return null;
  return leerJson(respuesta, "ERROR_SEGUIMIENTO", "No se pudo consultar el trámite");
}

export type ResumenPublico = {
  codigo: string;
  tipo: "CATEGORIZACION" | "REGULARIZACION";
  categoria: string;
  estado: EstadoTramite;
  radicadoEn: string;
  pendienteValidacion: boolean;
};

// Los trámites de un docente a partir de su código de docente Y su cédula (para quien perdió el código de
// seguimiento). Va por POST: la cédula no debe viajar en una URL. Una combinación que no coincide devuelve una
// lista vacía, igual que un docente sin trámites.
export async function consultarPorDocente(codigo: string, cedula: string): Promise<ResumenPublico[]> {
  const respuesta = await apiPublica("/api/publico/consulta-docente", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ codigo: codigo.trim(), cedula: cedula.trim() }),
  });
  return leerJson(respuesta, "ERROR_SEGUIMIENTO", "No se pudo consultar los trámites");
}

// Etiqueta de caché del catálogo público: las acciones de /panel/categorias la invalidan al guardar.
export const ETIQUETA_CATEGORIAS = "categorias-publicas";

// El catálogo cambia muy poco y lo piden la portada, /solicitud y los reportes: se guarda 5 minutos en la caché
// de datos de Next. Va sin las cabeceras del cliente (apiPublica) porque la respuesta es la misma para todos.
export async function listarCategoriasPublicas(): Promise<Categoria[]> {
  const apiUrl = process.env.API_URL_INTERNA;
  if (!apiUrl) throw new Error("API_URL_INTERNA no está configurada");
  await connection(); // la página se sigue generando por petición (la API no existe durante el build); solo se cachea el dato
  const respuesta = await fetch(`${apiUrl}/api/publico/categorias`, { next: { revalidate: 300, tags: [ETIQUETA_CATEGORIAS] } });
  return leerJson(respuesta, "ERROR_CATEGORIAS", "No se pudo cargar las categorías");
}
