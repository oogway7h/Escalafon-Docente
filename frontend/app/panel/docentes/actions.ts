"use server";

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { enviarJson } from "@/lib/api";
import { docenteEsquema } from "@/lib/esquemas";
import { camposDeZod, errorApi, leerErrorApi, type EstadoFormulario } from "@/lib/formularios";

export async function guardarDocenteAction(
  id: string | null,
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = docenteEsquema.safeParse({
    codigo: formData.get("codigo"),
    cedula: formData.get("cedula") || undefined,
    nombres: formData.get("nombres"),
    apellidos: formData.get("apellidos"),
    correo: formData.get("correo") || "",
    telefono: formData.get("telefono") || undefined,
    categoriaId: formData.get("categoriaId") || undefined,
    facultadId: formData.get("facultadId") || undefined,
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(id ? `/api/docentes/${id}` : "/api/docentes", id ? "PUT" : "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo guardar el docente.");

  revalidatePath("/panel/docentes");
  redirect("/panel/docentes");
}

// Baja lógica / reactivación desde el listado: el error viaja por la URL.
export async function cambiarEstadoDocenteAction(formData: FormData) {
  const id = formData.get("id") as string;
  const activo = formData.get("activo") === "true";

  const respuesta = await enviarJson(`/api/docentes/${id}/estado`, "PATCH", { activo });
  if (!respuesta.ok) {
    const { codigo } = await leerErrorApi(respuesta);
    redirect(`/panel/docentes?error=${encodeURIComponent(codigo ?? "ERROR")}`);
  }
  revalidatePath("/panel/docentes");
}
