import { api } from "@/lib/api";
import type { CoincidenciaExpediente } from "@/lib/api-tramites";
import { ETIQUETA_TIPO_DOCUMENTO } from "@/lib/tramites";
import type { Sugerencia } from "@/components/common/BuscadorSugerencias";

// Paso 2 del buscador: páginas del expediente del docente que contienen lo escrito (las 8 más relevantes).
export async function GET(request: Request) {
  const sp = new URL(request.url).searchParams;
  const docente = sp.get("docente") ?? "";
  const q = sp.get("q")?.trim() ?? "";
  if (!docente || q.length < 2) return Response.json([]);
  const r = await api(`/api/docentes/${encodeURIComponent(docente)}/documentos/buscar?q=${encodeURIComponent(q)}`);
  if (!r.ok) return Response.json([], { status: r.status });
  const { coincidencias }: { coincidencias: CoincidenciaExpediente[] } = await r.json();
  const lista: Sugerencia[] = coincidencias.slice(0, 8).map((c) => ({
    clave: `${c.documentoId}-${c.pagina}`,
    texto: c.nombreOriginal,
    detalle: `${ETIQUETA_TIPO_DOCUMENTO[c.tipoDocumento] ?? c.tipoDocumento} · página ${c.pagina}`,
    fragmento: c.fragmento,
    href: `/panel/documentos/buscar?${new URLSearchParams({ docente, q, doc: c.documentoId, pagina: String(c.pagina) })}`,
  }));
  return Response.json(lista);
}
