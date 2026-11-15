"use server";

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { enviarJson } from "@/lib/api";
import { cambiarRolEsquema, usuarioEsquema } from "@/lib/esquemas";
import { camposDeZod, errorApi, leerErrorApi, type EstadoFormulario } from "@/lib/formularios";

export async function crearUsuarioAction(
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = usuarioEsquema.safeParse({
    nombres: formData.get("nombres"),
    apellidos: formData.get("apellidos"),
    correo: formData.get("correo"),
    telefono: formData.get("telefono") || undefined,
    rolId: formData.get("rolId"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson("/api/usuarios", "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo crear el usuario.");

  revalidatePath("/panel/usuarios");
  redirect("/panel/usuarios");
}

export async function editarUsuarioAction(
  id: string,
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = usuarioEsquema.omit({ rolId: true }).safeParse({
    nombres: formData.get("nombres"),
    apellidos: formData.get("apellidos"),
    correo: formData.get("correo"),
    telefono: formData.get("telefono") || undefined,
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/usuarios/${id}`, "PUT", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo editar el usuario.");

  revalidatePath("/panel/usuarios");
  redirect("/panel/usuarios");
}

export async function cambiarRolAction(
  id: string,
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = cambiarRolEsquema.safeParse({ rolId: formData.get("rolId") });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson(`/api/usuarios/${id}/rol`, "PATCH", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo cambiar el rol.");

  revalidatePath("/panel/usuarios");
  redirect("/panel/usuarios");
}

// Acciones del menú de la tabla: sin useActionState, así que el error viaja por la URL.
async function redirigirConError(respuesta: Response): Promise<never> {
  const { codigo } = await leerErrorApi(respuesta);
  redirect(`/panel/usuarios?error=${encodeURIComponent(codigo ?? "ERROR")}`);
}

export async function cambiarEstadoAction(formData: FormData) {
  const id = formData.get("id") as string;
  const activo = formData.get("activo") === "true";

  const respuesta = await enviarJson(`/api/usuarios/${id}/estado`, "PATCH", { activo });
  if (!respuesta.ok) await redirigirConError(respuesta);

  revalidatePath("/panel/usuarios");
}

export async function reenviarActivacionAction(formData: FormData) {
  const id = formData.get("id") as string;

  const respuesta = await enviarJson(`/api/usuarios/${id}/reenviar-activacion`, "POST");
  if (!respuesta.ok) await redirigirConError(respuesta);

  redirect("/panel/usuarios?exito=REENVIADO");
}
