"use server";

import { redirect } from "next/navigation";
import { cabecerasCliente } from "@/lib/api";
import { loginEsquema } from "@/lib/esquemas";
import { camposDeZod, errorApi, type EstadoFormulario } from "@/lib/formularios";
import { escribirTokenSesion } from "@/lib/sesion";

export async function iniciarSesion(
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = loginEsquema.safeParse({
    correo: formData.get("correo"),
    contrasena: formData.get("contrasena"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const apiUrl = process.env.API_URL_INTERNA;
  if (!apiUrl) {
    return { error: "El servidor de autenticación no está configurado todavía." };
  }

  // fetch directo y no api(): api() trata el 401 como sesión vencida y redirige a /login,
  // lo que taparía el mensaje de credenciales inválidas.
  let respuesta: Response;
  try {
    respuesta = await fetch(`${apiUrl}/api/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json", ...(await cabecerasCliente()) },
      body: JSON.stringify(analisis.data),
      cache: "no-store",
    });
  } catch {
    return { error: "No se pudo conectar con el servidor. Intenta nuevamente." };
  }

  if (!respuesta.ok) return errorApi(respuesta, "No se pudo iniciar sesión.");

  const { token, expiraEn } = await respuesta.json();
  await escribirTokenSesion(token, expiraEn);

  redirect("/panel");
}
