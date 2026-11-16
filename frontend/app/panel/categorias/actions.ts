"use server";

import { redirect } from "next/navigation";
import { revalidatePath, revalidateTag } from "next/cache";
import { enviarJson } from "@/lib/api";
import { ETIQUETA_CATEGORIAS } from "@/lib/api-tramites";
import { categoriaEsquema } from "@/lib/esquemas";
import { camposDeZod, errorApi, leerErrorApi, type EstadoFormulario } from "@/lib/formularios";

export async function guardarCategoriaAction(
  id: string | null,
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = categoriaEsquema.safeParse({
    nombre: formData.get("nombre"),
    aniosAntiguedad: formData.get("aniosAntiguedad"),
    puntajeDesde: formData.get("puntajeDesde"),
    puntajeHasta: formData.get("puntajeHasta"),
    // Un requisito por línea.
    requisitos: String(formData.get("requisitos") ?? "")
      .split("\n")
      .map((r) => r.trim())
      .filter(Boolean),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(id ? `/api/categorias/${id}` : "/api/categorias", id ? "PUT" : "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo guardar la categoría.");

  revalidatePath("/panel/categorias");
  revalidateTag(ETIQUETA_CATEGORIAS, { expire: 0 }); // el portal público ve el cambio de inmediato
  redirect("/panel/categorias");
}

export async function eliminarCategoriaAction(
  id: number,
  _estadoPrevio: EstadoFormulario,
  _formData: FormData,
): Promise<EstadoFormulario> {
  const respuesta = await enviarJson(`/api/categorias/${id}`, "DELETE");
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo eliminar la categoría.");

  revalidatePath("/panel/categorias");
  revalidateTag(ETIQUETA_CATEGORIAS, { expire: 0 }); // el portal público ve el cambio de inmediato
  redirect("/panel/categorias");
}

export async function cambiarEstadoCategoriaAction(formData: FormData) {
  const id = formData.get("id") as string;
  const activo = formData.get("activo") === "true";

  const respuesta = await enviarJson(`/api/categorias/${id}/estado`, "PATCH", { activo });
  if (!respuesta.ok) {
    const { codigo } = await leerErrorApi(respuesta);
    redirect(`/panel/categorias?error=${encodeURIComponent(codigo ?? "ERROR")}`);
  }
  revalidatePath("/panel/categorias");
  revalidateTag(ETIQUETA_CATEGORIAS, { expire: 0 }); // el portal público ve el cambio de inmediato
}
