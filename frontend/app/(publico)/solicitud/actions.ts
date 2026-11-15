"use server";

import { redirect } from "next/navigation";
import { apiPublica } from "@/lib/api-tramites";
import { mensajeError, type EstadoFormulario } from "@/lib/formularios";
import { LIMITE_MB, MAX_ADJUNTOS_LIBRES } from "@/lib/tramites";

const MAX_BYTES = LIMITE_MB.portal * 1024 * 1024;

export async function radicarPublicoAction(_previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const codigoDocente = String(formData.get("codigo") ?? "").trim();
  const cedula = String(formData.get("cedula") ?? "").trim();
  const tipo = formData.get("tipo") === "REGULARIZACION" ? "REGULARIZACION" : "CATEGORIZACION";
  const categoriaId = String(formData.get("categoriaId") ?? "");
  const campos: Record<string, string> = {};
  if (!codigoDocente) campos.codigo = "Ingresa tu código de docente";
  if (cedula.length < 3) campos.cedula = "Ingresa tu cédula de identidad";
  if (tipo === "CATEGORIZACION" && !categoriaId) campos.categoriaId = "Selecciona la categoría a la que postulas";

  // Reenvía solo los campos esperados. Categorización: un adjunto por requisito (doc_<id>); regularización, que no
  // lleva categoría: adjuntos libres en el campo adjunto.
  const salida = new FormData();
  salida.set("codigo", codigoDocente);
  salida.set("cedula", cedula);
  salida.set("tipo", tipo);
  if (tipo === "CATEGORIZACION") salida.set("categoriaId", categoriaId);
  let libres = 0;
  for (const [clave, valor] of formData.entries()) {
    if (!(valor instanceof File) || valor.size === 0) continue;
    const esRequisito = tipo === "CATEGORIZACION" && clave.startsWith("doc_");
    const esLibre = tipo === "REGULARIZACION" && clave === "adjunto";
    if (!esRequisito && !esLibre) continue;
    if (valor.size > MAX_BYTES) campos[clave] = `El archivo supera los ${LIMITE_MB.portal} MB`;
    if (esLibre) {
      libres++;
      salida.append(clave, valor, valor.name);
    } else {
      salida.set(clave, valor, valor.name);
    }
  }
  if (tipo === "REGULARIZACION" && (libres < 1 || libres > MAX_ADJUNTOS_LIBRES)) {
    campos.adjunto ??= `Adjunta entre 1 y ${MAX_ADJUNTOS_LIBRES} documentos`;
  }
  if (Object.keys(campos).length > 0) return { camposInvalidos: campos };

  const respuesta = await apiPublica("/api/publico/solicitudes", { method: "POST", body: salida });
  if (!respuesta.ok) {
    const cuerpo = await respuesta.json().catch(() => null);
    if (cuerpo?.codigo === "DOCUMENTOS_FALTANTES") {
      return { error: `Faltan documentos: ${(cuerpo.detalles?.faltan ?? []).join(", ")}.` };
    }
    return { error: mensajeError(cuerpo?.codigo, "No se pudo radicar la solicitud. Inténtalo más tarde.") };
  }

  const { codigo } = await respuesta.json();
  redirect(`/seguimiento?codigo=${encodeURIComponent(codigo)}&nueva=1`);
}
