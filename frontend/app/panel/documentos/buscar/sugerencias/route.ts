import { api } from "@/lib/api";
import type { Docente } from "@/lib/api-tramites";
import type { Sugerencia } from "@/components/common/BuscadorSugerencias";

// Paso 1 del buscador: docentes que coinciden con lo escrito (código por prefijo, nombre sin tildes).
export async function GET(request: Request) {
  const q = new URL(request.url).searchParams.get("q")?.trim() ?? "";
  if (q.length < 2) return Response.json([]);
  const r = await api(`/api/docentes/sugerencias?q=${encodeURIComponent(q)}`);
  if (!r.ok) return Response.json([], { status: r.status });
  const docentes: Docente[] = await r.json();
  const lista: Sugerencia[] = docentes.map((d) => ({
    clave: d.id,
    texto: `${d.apellidos}, ${d.nombres}`,
    detalle: [d.codigo && `Código ${d.codigo}`, d.cedula && `CI ${d.cedula}`, d.categoria].filter(Boolean).join(" · "),
    href: `/panel/documentos/buscar?docente=${d.id}`,
  }));
  return Response.json(lista);
}
