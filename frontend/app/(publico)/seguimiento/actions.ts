"use server";

import { redirect } from "next/navigation";
import { consultarPorDocente, type ResumenPublico } from "@/lib/api-tramites";
import { MAX_CODIGO_DOCENTE } from "@/lib/tramites";

export type EstadoBusqueda = { error?: string; lista?: ResumenPublico[] };

// Un código con más de 5 caracteres es un código de seguimiento (los de docente tienen hasta 4 dígitos): abre el trámite
// directamente (MAX_CODIGO_DOCENTE en lib/tramites.ts). Si no, es un código de docente y exige su cédula; la cédula viaja por POST, nunca por la URL.

export async function consultarAction(_previo: EstadoBusqueda, formData: FormData): Promise<EstadoBusqueda> {
  const codigo = String(formData.get("codigo") ?? "").trim();
  if (!codigo) return { error: "Ingresa tu código de seguimiento o tu código de docente." };
  if (codigo.length > MAX_CODIGO_DOCENTE) redirect(`/seguimiento?codigo=${encodeURIComponent(codigo.toUpperCase())}`);

  const cedula = String(formData.get("cedula") ?? "").trim();
  if (!cedula) return { error: "Ingresa tu cédula de identidad junto con el código de docente." };

  let lista: ResumenPublico[];
  try {
    lista = await consultarPorDocente(codigo, cedula);
  } catch {
    return { error: "No se pudo consultar en este momento. Intenta de nuevo en unos minutos." };
  }
  // El mismo mensaje para código inexistente, cédula equivocada o sin trámites: no se revela cuál falló.
  if (lista.length === 0) return { error: "No encontramos trámites con esos datos. Revisa tu código de docente y tu cédula." };
  if (lista.length === 1) redirect(`/seguimiento?codigo=${encodeURIComponent(lista[0].codigo)}`);
  return { lista };
}
