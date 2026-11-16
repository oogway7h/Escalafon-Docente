"use server";

import { redirect } from "next/navigation";
import { enviarJson } from "@/lib/api";

// El usuario acepta la versión vigente del acuerdo de usuario (POST /api/auth/terminos) y entra al panel.
export async function aceptarTerminos(formData: FormData) {
  if (formData.get("acepto") !== "on") redirect("/panel");
  const respuesta = await enviarJson("/api/auth/terminos", "POST");
  if (!respuesta.ok) throw new Error("No se pudo registrar la aceptación del acuerdo");
  redirect("/panel");
}
