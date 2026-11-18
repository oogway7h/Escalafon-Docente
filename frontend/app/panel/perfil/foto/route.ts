import { api } from "@/lib/api";

// Foto de perfil del usuario con sesión. La API da un enlace firmado de corta vida; aquí se lee y se entrega
// al navegador desde el mismo origen (la CSP solo permite imágenes propias). No se guarda en disco: pesa
// poco y, al reemplazarse la foto, la clave del objeto no cambia (la versión viaja en ?v=).
export async function GET() {
  const r = await api("/api/perfil/foto");
  if (!r.ok) return new Response(null, { status: r.status === 404 ? 404 : 502 });
  const { url, tipoMime } = (await r.json()) as { url: string; tipoMime: string };
  const imagen = await fetch(url, { cache: "no-store" });
  if (!imagen.ok || !imagen.body) return new Response(null, { status: 502 });
  return new Response(imagen.body, {
    headers: {
      "Content-Type": tipoMime,
      "Cache-Control": "private, max-age=86400", // la URL lleva ?v=<versión>: una foto nueva es otra URL
      "X-Content-Type-Options": "nosniff",
    },
  });
}
