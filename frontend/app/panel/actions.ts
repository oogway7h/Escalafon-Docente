"use server";

import { redirect } from "next/navigation";
import { cerrarSesionApi } from "@/lib/api";
import { borrarTokenSesion } from "@/lib/sesion";

export async function cerrarSesion() {
  try {
    await cerrarSesionApi();
  } catch {
    // Si la API no responde, igual cerramos la sesión localmente.
  } finally {
    await borrarTokenSesion();
  }
  redirect("/login");
}
