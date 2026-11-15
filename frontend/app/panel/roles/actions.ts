"use server";

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { enviarJson } from "@/lib/api";
import { rolEsquema } from "@/lib/esquemas";
import { camposDeZod, errorApi, type EstadoFormulario } from "@/lib/formularios";

export async function guardarRolAction(
  id: number | null,
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = rolEsquema.safeParse({
    nombre: formData.get("nombre"),
    descripcion: formData.get("descripcion") || undefined,
    permisoIds: formData.getAll("permisoIds"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(id ? `/api/roles/${id}` : "/api/roles", id ? "PUT" : "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo guardar el rol.");

  revalidatePath("/panel/roles");
  redirect("/panel/roles");
}

export async function eliminarRolAction(
  id: number,
  _estadoPrevio: EstadoFormulario,
  _formData: FormData,
): Promise<EstadoFormulario> {
  const respuesta = await enviarJson(`/api/roles/${id}`, "DELETE");
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo eliminar el rol.");

  revalidatePath("/panel/roles");
  redirect("/panel/roles");
}
