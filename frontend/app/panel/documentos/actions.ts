"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import { enviarJson } from "@/lib/api";
import { errorApi, leerErrorApi, mensajeError, type EstadoFormulario } from "@/lib/formularios";

// Expediente digital del docente: Archivo sube el escaneo directo a S3 y el servidor lo procesa con OCR.

const expediente = (docenteId: string) => `/panel/documentos/${docenteId}`;

type SubidaPreparada = { documentoId?: string; url?: string; error?: string };

export async function prepararSubidaExpedienteAction(
  docenteId: string,
  datos: { nombre: string; tipoMime: string; tamano: number; tipoDocumento: string },
): Promise<SubidaPreparada> {
  const respuesta = await enviarJson(`/api/docentes/${docenteId}/documentos/subida`, "POST", datos);
  if (!respuesta.ok) return { error: (await errorApi(respuesta, "No se pudo preparar la subida.")).error };
  const { documento, url } = await respuesta.json();
  return { documentoId: documento.id, url };
}

// Verifica el objeto subido y lanza el OCR (también reintenta un documento en ERROR).
export async function confirmarSubidaExpedienteAction(docenteId: string, docId: string): Promise<{ error?: string }> {
  const respuesta = await enviarJson(`/api/docentes/${docenteId}/documentos/${docId}/confirmar`, "POST");
  if (!respuesta.ok) {
    const { codigo } = await leerErrorApi(respuesta);
    return { error: mensajeError(codigo, "No se pudo procesar el documento.") };
  }
  revalidatePath(expediente(docenteId));
  return {};
}

export async function darDeBajaExpedienteAction(docenteId: string, docId: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const motivo = String(formData.get("motivo") ?? "").trim();
  if (motivo.length < 5) return { error: "Describe el motivo de la baja (mínimo 5 caracteres)." };
  const respuesta = await enviarJson(`/api/docentes/${docenteId}/documentos/${docId}/estado`, "PATCH", { motivo });
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo dar de baja el documento.");
  revalidatePath(expediente(docenteId));
  return { mensaje: "Documento dado de baja." };
}

// CU-17: el Director restituye un documento dado de baja; vuelve a las búsquedas sin reprocesar el OCR.
export async function restituirDocumentoAction(docId: string, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const motivo = String(formData.get("motivo") ?? "").trim();
  if (motivo.length < 5) return { error: "Describe el motivo de la restitución (mínimo 5 caracteres)." };
  const respuesta = await enviarJson(`/api/documentos/${docId}/restituir`, "POST", { motivo });
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo restituir el documento.");
  // El documento sale de la bandeja: el aviso va en la página, no en el formulario que desaparece.
  revalidatePath("/panel/documentos/bajas");
  redirect("/panel/documentos/bajas?restituido=1");
}
