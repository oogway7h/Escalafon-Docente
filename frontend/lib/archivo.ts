import "server-only";
import { createHash } from "node:crypto";
import { createReadStream, createWriteStream } from "node:fs";
import { mkdir, readdir, rename, rm, stat } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import { Readable } from "node:stream";
import { pipeline } from "node:stream/promises";
import { api } from "./api";
import { leerTokenSesion } from "./sesion";

// Entrega un documento a través del propio servidor: el visor lo lee con la sesión del usuario y con
// soporte de rangos (pdf.js pide el PDF por tramos), sin que el bucket necesite CORS de lectura.
//
// Rendimiento: antes cada tramo pasaba por la API (permisos + URL firmada) y luego por S3. Ahora
//  1. la URL firmada se recuerda por sesión y documento unos minutos (los tramos siguientes no tocan la API);
//  2. la primera vez el objeto se descarga entero a disco en segundo plano y desde ahí se sirven los
//     tramos (y las aperturas siguientes, de cualquier usuario autorizado) sin ir a S3;
//  3. el navegador lo guarda en su caché (la ruta lleva ?v= con la versión del archivo).

const DIRECTORIO = path.join(tmpdir(), "escalafon-documentos");
const LIMITE_DISCO = 2 * 1024 ** 3; // ponytail: limpieza por antigüedad al pasar 2 GB; un LRU real si el disco del servidor aprieta
const VIGENCIA_URL = 10 * 60 * 1000; // la API firma por 30 min; los permisos se vuelven a comprobar a los 10

// Caché del navegador corta: un documento dado de baja deja de poder abrirse desde la caché en 10 minutos.
const CACHE_NAVEGADOR = "private, max-age=600";

const urls = new Map<string, { url: string; vence: number }>();
const descargas = new Map<string, Promise<void>>();

const sha = (s: string) => createHash("sha256").update(s).digest("hex").slice(0, 32);
const EXTENSION: Record<string, string> = { "application/pdf": "pdf", "image/png": "png", "image/jpeg": "jpg" };
const TIPO: Record<string, string> = { pdf: "application/pdf", png: "image/png", jpg: "image/jpeg" };

async function urlFirmada(rutaApi: string): Promise<string | Response> {
  const clave = `${sha((await leerTokenSesion()) ?? "")}|${rutaApi}`;
  const memo = urls.get(clave);
  if (memo && memo.vence > Date.now()) return memo.url;
  const respuesta = await api(rutaApi);
  if (!respuesta.ok) return new Response("No se pudo abrir el documento.", { status: respuesta.status });
  const { url } = (await respuesta.json()) as { url: string };
  urls.set(clave, { url, vence: Date.now() + VIGENCIA_URL });
  return url;
}

// Copia local del objeto: la clave sale de la ruta del objeto en el bucket (sin la firma), así el original
// y su versión comprimida son entradas distintas.
async function enDisco(url: string): Promise<{ archivo: string; tamano: number; tipo: string } | null> {
  const base = path.join(DIRECTORIO, sha(new URL(url).pathname));
  for (const [ext, tipo] of Object.entries(TIPO)) {
    const s = await stat(`${base}.${ext}`).catch(() => null);
    if (s) return { archivo: `${base}.${ext}`, tamano: s.size, tipo };
  }
  return null;
}

function descargarEnSegundoPlano(url: string) {
  const base = path.join(DIRECTORIO, sha(new URL(url).pathname));
  if (descargas.has(base)) return;
  const tarea = (async () => {
    await mkdir(DIRECTORIO, { recursive: true });
    const r = await fetch(url, { cache: "no-store" });
    // El tipo sale del almacén; si viene genérico (p. ej. la versión comprimida), de la extensión del objeto.
    const porRuta = path.extname(new URL(url).pathname).slice(1).toLowerCase().replace("jpeg", "jpg");
    const ext = EXTENSION[r.headers.get("content-type")?.split(";")[0] ?? ""] ?? (porRuta in TIPO ? porRuta : undefined);
    if (!r.ok || !r.body || !ext) return;
    const temporal = `${base}.${process.pid}.parcial`;
    try {
      await pipeline(Readable.fromWeb(r.body as never), createWriteStream(temporal));
      await rename(temporal, `${base}.${ext}`);
    } catch {
      await rm(temporal, { force: true });
      return;
    }
    await limpiar();
  })()
    .catch(() => {})
    .finally(() => descargas.delete(base));
  descargas.set(base, tarea);
}

// Borra los archivos más antiguos cuando la caché supera el límite.
async function limpiar() {
  const nombres = (await readdir(DIRECTORIO)).filter((n) => !n.endsWith(".parcial"));
  const archivos = await Promise.all(
    nombres.map(async (n) => {
      const s = await stat(path.join(DIRECTORIO, n));
      return { ruta: path.join(DIRECTORIO, n), tamano: s.size, uso: s.mtimeMs };
    }),
  );
  let total = archivos.reduce((t, a) => t + a.tamano, 0);
  for (const a of archivos.sort((x, y) => x.uso - y.uso)) {
    if (total <= LIMITE_DISCO) break;
    await rm(a.ruta, { force: true });
    total -= a.tamano;
  }
}

function servirDesdeDisco(request: Request, local: { archivo: string; tamano: number; tipo: string }): Response {
  const cabeceras = new Headers({
    "cache-control": CACHE_NAVEGADOR,
    "accept-ranges": "bytes",
    "content-type": local.tipo,
  });
  const m = /^bytes=(\d*)-(\d*)$/.exec(request.headers.get("range") ?? "");
  if (!m || (!m[1] && !m[2])) {
    cabeceras.set("content-length", String(local.tamano));
    return new Response(Readable.toWeb(createReadStream(local.archivo)) as ReadableStream, { status: 200, headers: cabeceras });
  }
  // bytes=a-b, bytes=a- o bytes=-n (los últimos n bytes, donde pdf.js busca el índice del PDF).
  let inicio = m[1] ? Number(m[1]) : Math.max(0, local.tamano - Number(m[2]));
  let fin = m[1] && m[2] ? Number(m[2]) : local.tamano - 1;
  fin = Math.min(fin, local.tamano - 1);
  if (inicio > fin) {
    cabeceras.set("content-range", `bytes */${local.tamano}`);
    return new Response(null, { status: 416, headers: cabeceras });
  }
  inicio = Math.max(0, inicio);
  cabeceras.set("content-range", `bytes ${inicio}-${fin}/${local.tamano}`);
  cabeceras.set("content-length", String(fin - inicio + 1));
  return new Response(Readable.toWeb(createReadStream(local.archivo, { start: inicio, end: fin })) as ReadableStream, { status: 206, headers: cabeceras });
}

export async function entregarArchivo(request: Request, rutaApi: string): Promise<Response> {
  const url = await urlFirmada(rutaApi);
  if (url instanceof Response) return url;

  const local = await enDisco(url);
  if (local) return servirDesdeDisco(request, local);
  descargarEnSegundoPlano(url);

  // Mientras se descarga la copia local, este tramo se pide directo a S3.
  const rango = request.headers.get("range");
  const origen = await fetch(url, { headers: rango ? { range: rango } : {}, cache: "no-store" });
  if (!origen.ok && origen.status !== 206) return new Response("El almacén no entregó el documento.", { status: 502 });

  const cabeceras = new Headers({ "cache-control": CACHE_NAVEGADOR, "accept-ranges": "bytes" });
  for (const nombre of ["content-type", "content-length", "content-range"]) {
    const valor = origen.headers.get(nombre);
    if (valor) cabeceras.set(nombre, valor);
  }
  return new Response(origen.body, { status: origen.status, headers: cabeceras });
}
