"use server";

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { enviarJson } from "@/lib/api";
import { listarCarpetas } from "@/lib/api-legajos";
import { listarTramites, type Docente } from "@/lib/api-tramites";
import { errorApi, leerErrorApi, type EstadoFormulario } from "@/lib/formularios";
import { buscarDocentePorCodigoAction } from "../tramites/actions";

// CU-20: el analista solicita una carpeta. El error viaja por la URL de la pantalla de origen.
export async function solicitarCarpetaAction(formData: FormData) {
  const tramiteId = String(formData.get("tramiteId") ?? "");
  const respuesta = await enviarJson(`/api/legajos/${formData.get("legajoId")}/solicitudes`, "POST", {
    tramiteId,
    motivo: String(formData.get("motivo") ?? "").trim(),
  });
  if (!respuesta.ok) {
    const { codigo } = await leerErrorApi(respuesta);
    // 422 de validación (motivo vacío) llega con código VALIDACION.
    redirect(`/panel/carpetas/nueva?tramite=${encodeURIComponent(tramiteId)}&error=${encodeURIComponent(codigo === "VALIDACION" ? "MOTIVO_REQUERIDO" : (codigo ?? "ERROR"))}`);
  }
  revalidatePath("/panel/carpetas");
  revalidatePath(`/panel/tramites/${tramiteId}`);
  redirect("/panel/carpetas?exito=SOLICITADA");
}

// --- Solicitud desde el modal: se busca al docente por código y se muestran sus datos ---

export type DatosSolicitud = {
  docente: Docente | null;
  tramites: { id: string; codigo: string; descripcion: string }[]; // trámites abiertos a cargo del analista para ese docente (ligar uno es opcional)
  carpetas: { id: number; codigo: string; gaveta: string }[]; // carpetas del docente disponibles en el archivo
};

export async function datosParaSolicitudAction(codigo: string): Promise<DatosSolicitud> {
  const docente = await buscarDocentePorCodigoAction(codigo);
  if (!docente) return { docente: null, tramites: [], carpetas: [] };
  const [mios, carpetas] = await Promise.all([
    listarTramites({ mios: true }),
    listarCarpetas({ buscar: docente.codigo, estado: "EN_ARCHIVO" }, 1, 20),
  ]);
  return {
    docente,
    tramites: mios
      .filter((t) => t.docenteCodigo === docente.codigo && t.estado !== "ARCHIVADO" && t.estado !== "REMITIDO")
      .map((t) => ({ id: t.id, codigo: t.codigo, descripcion: t.categoria ? `${t.codigo} · ${t.categoria}` : t.codigo })),
    carpetas: carpetas.datos.filter((c) => c.docenteId === docente.id).map((c) => ({ id: c.id, codigo: c.codigo, gaveta: c.gaveta })),
  };
}

export async function solicitarCarpetaModalAction(_previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const motivo = String(formData.get("motivo") ?? "").trim();
  if (!motivo) return { error: "Indica el motivo de la solicitud." };
  // El trámite es opcional: se puede pedir cualquier carpeta aunque no haya un trámite abierto de ese docente.
  const tramiteId = String(formData.get("tramiteId") ?? "");
  const respuesta = await enviarJson(`/api/legajos/${formData.get("legajoId")}/solicitudes`, "POST", { ...(tramiteId ? { tramiteId } : {}), motivo });
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo enviar la solicitud.");
  revalidatePath("/panel/carpetas");
  if (tramiteId) revalidatePath(`/panel/tramites/${tramiteId}`);
  redirect("/panel/carpetas?exito=SOLICITADA");
}
