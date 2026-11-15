"use server";

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { enviarJson } from "@/lib/api";
import { listarDocentesPagina, type Docente } from "@/lib/api-tramites";
import {
  archivarEsquema,
  asignarEsquema,
  bajaDocumentoEsquema,
  dictamenEsquema,
  observacionEsquema,
  periodoEsquema,
  radicarEsquema,
  regresoEsquema,
  remisionEsquema,
} from "@/lib/esquemas";
import { camposDeZod, errorApi, leerErrorApi, mensajeError, type EstadoFormulario } from "@/lib/formularios";

const detalle = (id: string) => `/panel/tramites/${id}`;

function terminar(id: string, exito: string): never {
  revalidatePath("/panel/tramites", "layout");
  redirect(`${detalle(id)}?exito=${exito}`);
}

// --- CU-08 ---

// Busca por código exacto (en vez de cargar los ~1600 docentes en un <select>).
export async function buscarDocentePorCodigoAction(codigo: string): Promise<Docente | null> {
  const buscado = codigo.trim();
  if (!buscado) return null;
  const { datos } = await listarDocentesPagina(buscado, 1, 20);
  return datos.find((d) => d.codigo.toLowerCase() === buscado.toLowerCase()) ?? null;
}

export async function radicarAction(_previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = radicarEsquema.safeParse({
    docenteId: formData.get("docenteId"),
    tipo: formData.get("tipo"),
    categoriaId: formData.get("categoriaId"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson("/api/tramites", "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo radicar el trámite.");

  const creado = await respuesta.json();
  terminar(creado.id, "RADICADO");
}

// --- Acciones de un solo botón: el error viaja por la URL ---

async function accionSimple(id: string, ruta: string, exito: string, cuerpo?: unknown): Promise<never> {
  const respuesta = await enviarJson(`/api/tramites/${id}${ruta}`, "POST", cuerpo);
  if (!respuesta.ok) {
    const { codigo } = await leerErrorApi(respuesta);
    redirect(`${detalle(id)}?error=${encodeURIComponent(codigo ?? "ERROR")}`);
  }
  terminar(id, exito);
}

export async function validarAction(formData: FormData) {
  return accionSimple(formData.get("id") as string, "/validar", "VALIDADO");
}

export async function iniciarRevisionAction(formData: FormData) {
  return accionSimple(formData.get("id") as string, "/revision", "REVISION");
}

export async function resolverObservacionAction(formData: FormData) {
  const oid = formData.get("observacionId") as string;
  return accionSimple(formData.get("id") as string, `/observaciones/${oid}/resolver`, "RESUELTA");
}

// --- CU-09 ---

export async function asignarAction(id: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = asignarEsquema.safeParse({
    analistaId: formData.get("analistaId"),
    motivo: formData.get("motivo") || undefined,
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/tramites/${id}/asignar`, "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo asignar el trámite.");
  terminar(id, "ASIGNADO");
}

// --- CU-10 ---

export async function observarAction(id: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = observacionEsquema.safeParse({
    documento: formData.get("documento"),
    motivo: formData.get("motivo"),
    plazo: formData.get("plazo"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/tramites/${id}/observaciones`, "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo registrar la observación.");
  terminar(id, "OBSERVADO");
}

// --- CU-11 ---

export async function dictaminarAction(id: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = dictamenEsquema.safeParse({ resultado: formData.get("resultado"), texto: formData.get("texto") });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/tramites/${id}/dictamen`, "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo emitir el dictamen.");
  terminar(id, "DICTAMINADO");
}

export async function archivarAction(id: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = archivarEsquema.safeParse({ ubicacion: formData.get("ubicacion") });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/tramites/${id}/archivar`, "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo archivar el trámite.");
  terminar(id, "ARCHIVADO");
}

// --- Etapas posteriores al dictamen ---

export async function periodoAction(id: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = periodoEsquema.safeParse({ anioDesde: formData.get("anioDesde"), anioHasta: formData.get("anioHasta") });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/tramites/${id}/periodo`, "PUT", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudieron registrar los años a evaluar.");
  terminar(id, "PERIODO");
}

// El analista remite el expediente a la facultad: desde ahí el trámite no se puede rastrear.
export async function remitirFacultadAction(formData: FormData) {
  return accionSimple(formData.get("id") as string, "/remitir-facultad", "FACULTAD");
}

export async function regresoAction(id: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = regresoEsquema.safeParse({
    resolucionNumero: formData.get("resolucionNumero"),
    resolucionFecha: formData.get("resolucionFecha"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/tramites/${id}/regreso`, "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo registrar el regreso.");
  terminar(id, "REGRESO");
}

export async function remisionAction(id: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = remisionEsquema.safeParse({ nota: formData.get("nota") });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/tramites/${id}/remitir-rectorado`, "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo remitir al Rectorado.");
  terminar(id, "RECTORADO");
}

// --- CU-16 ---

// Baja lógica de un adjunto cargado por error: el archivo se conserva y el motivo queda en la bitácora.
export async function darDeBajaDocumentoAction(id: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = bajaDocumentoEsquema.safeParse({
    documentoId: formData.get("documentoId"),
    motivo: formData.get("motivo"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const { documentoId, motivo } = analisis.data;
  const respuesta = await enviarJson(`/api/tramites/${id}/documentos/${documentoId}/estado`, "PATCH", { motivo });
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo dar de baja el documento.");
  terminar(id, "DOCUMENTO_BAJA");
}

// --- CU-15: carga de documentos con OCR ---

type SubidaPreparada = { documentoId?: string; url?: string; error?: string };

// Registra el documento y devuelve la URL firmada: el navegador sube el archivo directo a S3.
export async function prepararSubidaAction(
  id: string,
  datos: { nombre: string; tipoMime: string; tamano: number; requisitoId: number | null },
): Promise<SubidaPreparada> {
  const respuesta = await enviarJson(`/api/tramites/${id}/documentos/subida`, "POST", datos);
  if (!respuesta.ok) return { error: (await errorApi(respuesta, "No se pudo preparar la subida.")).error };
  const { documento, url } = await respuesta.json();
  return { documentoId: documento.id, url };
}

// Verifica el objeto subido y lanza el OCR (también reintenta un documento en ERROR).
export async function confirmarSubidaAction(id: string, docId: string): Promise<{ error?: string }> {
  const respuesta = await enviarJson(`/api/tramites/${id}/documentos/${docId}/confirmar`, "POST");
  if (!respuesta.ok) {
    const { codigo } = await leerErrorApi(respuesta);
    return { error: mensajeError(codigo, "No se pudo procesar el documento.") };
  }
  revalidatePath(detalle(id));
  return {};
}

export async function corregirTextoAction(id: string, docId: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const respuesta = await enviarJson(`/api/tramites/${id}/documentos/${docId}/texto`, "PATCH", { paginas: formData.getAll("pagina").map(String) });
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo guardar el texto.");
  revalidatePath(detalle(id));
  return { mensaje: "Texto guardado y marcado como validado." };
}
