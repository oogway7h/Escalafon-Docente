import type { Metadata } from "next";
import Link from "next/link";
import { listarCategoriasPublicas } from "@/lib/api-tramites";
import SolicitudForm from "./SolicitudForm";

export const metadata: Metadata = { title: "Nueva solicitud — Escalafón Docente" };

export default async function SolicitudPage() {
  const categorias = await listarCategoriasPublicas();

  return (
    <div className="flex flex-col items-center gap-10">
      <div className="flex max-w-2xl flex-col items-center gap-4 text-center">
        <h1 className="font-display text-[clamp(2rem,5vw,3.25rem)] font-semibold leading-[1.05] tracking-[-0.025em] text-ink text-balance">
          Radica tu solicitud desde casa
        </h1>
        <p className="max-w-md text-base leading-relaxed text-ink-soft">
          Adjunta los documentos que exige tu categoría. Recibirás un código para seguir tu trámite; un funcionario validará tus documentos en ventanilla.
        </p>
      </div>

      <SolicitudForm categorias={categorias} />

      <p className="text-sm text-ink-soft">
        ¿Ya radicaste?{" "}
        <Link href="/seguimiento" className="font-semibold text-accent-ink hover:underline">
          Consulta tu trámite
        </Link>
      </p>
    </div>
  );
}
