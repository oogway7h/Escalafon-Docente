"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { PDFDocumentProxy } from "pdfjs-dist";
import Loader from "@/components/common/Loader";

// Visor del sistema: renderiza el PDF con pdf.js (canvas) en lugar del visor del navegador, para que se vea
// y se comporte como parte de la aplicación (barra de herramientas, miniaturas, búsqueda y resaltado).
// Las páginas van una debajo de otra y se recorren con el scroll; solo se dibujan las que están cerca de
// la vista (un PDF de cientos de páginas no llena la memoria de lienzos).
// public/pdf.worker.min.mjs es el worker de pdf.js: debe ser de la misma versión que pdfjs-dist (cópialo de
// node_modules/pdfjs-dist/build al actualizar la dependencia).

type Pdfjs = typeof import("pdfjs-dist");
let promesaPdfjs: Promise<Pdfjs> | null = null;
function cargarPdfjs(): Promise<Pdfjs> {
  promesaPdfjs ??= import("pdfjs-dist").then((m) => {
    m.GlobalWorkerOptions.workerSrc = "/pdf.worker.min.mjs";
    return m;
  });
  return promesaPdfjs;
}

const sinTildes = (s: string) => s.normalize("NFD").replace(/[̀-ͯ]/g, "").toLowerCase();

type Marca = { left: number; top: number; width: number; height: number };
type Medida = { w: number; h: number }; // tamaño de la página a escala 1
const SEPARACION = 16; // px entre páginas

function Icono({ d }: { d: string[] }) {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      {d.map((p) => (
        <path key={p} d={p} />
      ))}
    </svg>
  );
}

const botonHerramienta =
  "flex h-9 w-9 items-center justify-center rounded-lg text-white transition-colors hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-40";

function Miniatura({ pdf, n, activa, coincide, alElegir }: { pdf: PDFDocumentProxy; n: number; activa: boolean; coincide: boolean; alElegir: () => void }) {
  const lienzo = useRef<HTMLCanvasElement>(null);
  useEffect(() => {
    let cancelado = false;
    let tarea: { cancel: () => void } | null = null;
    // Con retardo: dibujar una miniatura de un escaneo cuesta casi lo mismo que la página entera (hay que
    // decodificar la imagen completa) y pdf.js tiene un solo trabajador. Primero va la página que se está leyendo.
    const espera = setTimeout(async () => {
      const page = await pdf.getPage(n);
      const base = page.getViewport({ scale: 1 });
      const viewport = page.getViewport({ scale: 64 / base.width });
      const c = lienzo.current;
      const ctx = c?.getContext("2d");
      if (cancelado || !c || !ctx) return;
      c.width = Math.floor(viewport.width);
      c.height = Math.floor(viewport.height);
      const t = page.render({ canvasContext: ctx, viewport });
      tarea = t;
      await t.promise.catch(() => {});
    }, 600);
    return () => {
      cancelado = true;
      clearTimeout(espera);
      tarea?.cancel();
    };
  }, [pdf, n]);

  return (
    <button type="button" onClick={alElegir} aria-label={`Ir a la página ${n}`} aria-current={activa ? "page" : undefined} className="flex flex-col items-center gap-1.5">
      <span className={`relative block rounded-[4px] bg-white ${activa ? "ring-2 ring-accent" : "ring-1 ring-line"}`}>
        <canvas ref={lienzo} className="block h-auto w-16 rounded-[4px]" />
        {coincide && <span className="absolute -right-1 -top-1 h-2.5 w-2.5 rounded-full border-2 border-surface bg-accent" />}
      </span>
      <span className={`text-xs ${activa ? "font-semibold text-ink" : "text-ink-faint"}`}>{n}</span>
    </button>
  );
}

// Una página del documento. Fuera de la vista es solo un recuadro blanco del tamaño correcto; cerca de ella
// se dibuja en un lienzo y se marcan las palabras buscadas sobre la capa de texto que dejó el OCR.
function Pagina({
  pdf,
  pdfjs,
  n,
  medida,
  escala,
  dibujar,
  terminos,
  nombre,
  alMedir,
  refPagina,
}: {
  pdf: PDFDocumentProxy;
  pdfjs: Pdfjs;
  n: number;
  medida: Medida;
  escala: number;
  dibujar: boolean;
  terminos: string[];
  nombre: string;
  alMedir: (n: number, m: Medida) => void;
  refPagina: (n: number, el: HTMLDivElement | null) => void;
}) {
  const lienzo = useRef<HTMLCanvasElement>(null);
  // Las marcas se guardan con la búsqueda y la escala para las que se calcularon: si alguna cambió, no se muestran.
  const [calculadas, setCalculadas] = useState<{ terminos: string[]; escala: number; lista: Marca[] } | null>(null);
  const marcas = dibujar && calculadas && calculadas.terminos === terminos && calculadas.escala === escala ? calculadas.lista : [];
  const ancho = medida.w * escala;
  const alto = medida.h * escala;

  useEffect(() => {
    if (!dibujar) return;
    let cancelado = false;
    let tarea: { cancel: () => void; promise: Promise<unknown> } | null = null;
    (async () => {
      const page = await pdf.getPage(n);
      const base = page.getViewport({ scale: 1 });
      if (base.width !== medida.w || base.height !== medida.h) alMedir(n, { w: base.width, h: base.height });
      const viewport = page.getViewport({ scale: escala * (medida.w / base.width) });
      const c = lienzo.current;
      const ctx = c?.getContext("2d");
      if (cancelado || !c || !ctx) return;
      const dpr = window.devicePixelRatio || 1;
      c.width = Math.floor(viewport.width * dpr);
      c.height = Math.floor(viewport.height * dpr);
      const t = page.render({ canvasContext: ctx, viewport, transform: dpr === 1 ? undefined : [dpr, 0, 0, dpr, 0, 0] });
      tarea = t;
      await t.promise.catch(() => {});
    })();
    return () => {
      cancelado = true;
      tarea?.cancel();
    };
  }, [pdf, n, dibujar, escala, medida.w, medida.h, alMedir]);

  // Resaltado de lo buscado: va aparte del dibujo. Cambiar la búsqueda solo recalcula las marcas (leer la capa
  // de texto es barato); antes volvía a rasterizar todas las páginas cercanas.
  useEffect(() => {
    if (!dibujar || terminos.length === 0) return;
    let cancelado = false;
    (async () => {
      const page = await pdf.getPage(n);
      const base = page.getViewport({ scale: 1 });
      const viewport = page.getViewport({ scale: escala * (medida.w / base.width) });
      const nuevas: Marca[] = [];
      {
        const contenido = await page.getTextContent();
        const k = viewport.scale;
        for (const item of contenido.items) {
          if (!("str" in item) || !item.str) continue;
          const texto = sinTildes(item.str);
          const [, , , , e, f] = pdfjs.Util.transform(viewport.transform, item.transform);
          const altoTexto = Math.hypot(item.transform[2], item.transform[3]) * k;
          const anchoTexto = item.width * k;
          for (const termino of terminos) {
            let desde = texto.indexOf(termino);
            while (desde !== -1) {
              nuevas.push({
                left: e + (desde / texto.length) * anchoTexto,
                top: f - altoTexto,
                width: Math.max(4, (termino.length / texto.length) * anchoTexto),
                height: altoTexto * 1.15,
              });
              desde = texto.indexOf(termino, desde + termino.length);
            }
          }
        }
      }
      if (!cancelado) setCalculadas({ terminos, escala, lista: nuevas });
    })();
    return () => {
      cancelado = true;
    };
  }, [pdf, pdfjs, n, dibujar, escala, medida.w, medida.h, terminos]);

  return (
    <div
      ref={(el) => refPagina(n, el)}
      data-pagina={n}
      className="relative mx-auto shrink-0 bg-white shadow-[0_1px_3px_rgba(15,27,46,0.12),0_18px_40px_rgba(15,27,46,0.12)]"
      style={{ width: ancho, height: alto }}
    >
      {dibujar && (
        <>
          <canvas ref={lienzo} className="block" style={{ width: ancho, height: alto }} aria-label={`Página ${n} de ${nombre}`} />
          {marcas.map((m, i) => (
            <span key={i} className="pointer-events-none absolute rounded-[3px] bg-accent/25 outline outline-2 outline-accent/70" style={{ left: m.left, top: m.top, width: m.width, height: m.height }} />
          ))}
        </>
      )}
    </div>
  );
}

export default function VisorPdf({
  src,
  esImagen,
  nombre,
  paginaInicial = 1,
  salto = 0,
  consulta = "",
  paginasConCoincidencia = [],
}: {
  src: string; // /panel/.../archivo/<docId>: el servidor lo entrega con rangos
  esImagen: boolean;
  nombre: string;
  paginaInicial?: number;
  salto?: number; // cambia cada vez que quien lo usa pide ir a paginaInicial (aunque sea la misma página otra vez)
  consulta?: string;
  paginasConCoincidencia?: number[];
}) {
  const escenario = useRef<HTMLDivElement>(null);
  const paginasDom = useRef(new Map<number, HTMLDivElement>());
  const [pdfjs, setPdfjs] = useState<Pdfjs | null>(null);
  const [pdf, setPdf] = useState<PDFDocumentProxy | null>(null);
  const [medidas, setMedidas] = useState<Medida[]>([]);
  const [pagina, setPagina] = useState(paginaInicial);
  const [zoom, setZoom] = useState(1);
  const [ancho, setAncho] = useState(0);
  const [cerca, setCerca] = useState<Set<number>>(() => new Set());
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(!esImagen);
  // Bytes del documento ya recibidos por el navegador. pdf.js lo descarga entero en segundo plano: mientras
  // llega, las páginas lejanas tardan; después, todo sale de la memoria.
  const [descarga, setDescarga] = useState<{ cargado: number; total: number } | null>(null);
  const [campo, setCampo] = useState(String(paginaInicial));
  const total = medidas.length;
  const listo = Boolean(pdf && pdfjs && ancho > 0);

  // Ancho disponible: cada página se ajusta al ancho del escenario.
  useEffect(() => {
    const el = escenario.current;
    if (!el) return;
    const medir = () => setAncho(Math.max(240, el.clientWidth - 48));
    medir();
    const o = new ResizeObserver(medir);
    o.observe(el);
    return () => o.disconnect();
  }, []);

  // Carga del documento. Tramos de 1 MB: menos idas y vueltas al servidor que los 64 KB por defecto.
  useEffect(() => {
    if (esImagen) return;
    let cancelado = false;
    let tarea: { destroy: () => Promise<void> } | null = null;
    (async () => {
      try {
        const lib = await cargarPdfjs();
        const t = lib.getDocument({ url: src, rangeChunkSize: 1 << 20 });
        tarea = t;
        t.onProgress = ({ loaded, total }: { loaded: number; total: number }) => {
          if (!cancelado && total > 0) setDescarga({ cargado: Math.min(loaded, total), total });
        };
        const doc = await t.promise;
        const primera = (await doc.getPage(1)).getViewport({ scale: 1 });
        if (cancelado) return;
        setPdfjs(lib);
        setPdf(doc);
        // Hasta dibujar cada página se supone del tamaño de la primera (en un escaneo suelen ser iguales).
        setMedidas(Array.from({ length: doc.numPages }, () => ({ w: primera.width, h: primera.height })));
      } catch {
        if (!cancelado) setError("No se pudo abrir el documento.");
      } finally {
        if (!cancelado) setCargando(false);
      }
    })();
    return () => {
      cancelado = true;
      void tarea?.destroy();
    };
  }, [src, esImagen]);

  const alMedir = useCallback((n: number, m: Medida) => {
    setMedidas((prev) => (prev[n - 1] && prev[n - 1].w === m.w && prev[n - 1].h === m.h ? prev : prev.map((x, i) => (i === n - 1 ? m : x))));
  }, []);
  const refPagina = useCallback((n: number, el: HTMLDivElement | null) => {
    if (el) paginasDom.current.set(n, el);
    else paginasDom.current.delete(n);
  }, []);

  // Qué páginas dibujar: las que están a menos de una pantalla de la vista.
  useEffect(() => {
    const raiz = escenario.current;
    if (!listo || !raiz) return;
    const o = new IntersectionObserver(
      (entradas) =>
        setCerca((prev) => {
          const s = new Set(prev);
          for (const e of entradas) {
            const n = Number((e.target as HTMLElement).dataset.pagina);
            if (e.isIntersecting) s.add(n);
            else s.delete(n);
          }
          return s;
        }),
      { root: raiz, rootMargin: "100% 0px" },
    );
    for (const el of paginasDom.current.values()) o.observe(el);
    return () => o.disconnect();
  }, [listo, total]);

  // Desplaza el escenario hasta la página n. Solo se anima el paso a una página vecina: un desplazamiento suave
  // a través de decenas de páginas tarda y hace dibujar (y cancelar) todas las intermedias antes de la buscada.
  const desplazarA = useCallback((n: number, suave: boolean) => {
    const raiz = escenario.current;
    const el = paginasDom.current.get(n);
    if (!raiz || !el) return;
    const top = el.offsetTop - SEPARACION;
    const cerca = Math.abs(top - raiz.scrollTop) <= raiz.clientHeight * 1.5;
    raiz.scrollTo({ top, behavior: suave && cerca ? "smooth" : "auto" });
  }, []);

  const ir = (n: number) => {
    const destino = Math.min(Math.max(1, Math.round(n) || 1), total || 1);
    setPagina(destino);
    setCampo(String(destino));
    desplazarA(destino, true);
  };

  // Al abrir el documento (o al cambiar el zoom) se coloca en la página pedida / la que se estaba leyendo.
  const paginaRef = useRef(pagina);
  useEffect(() => {
    paginaRef.current = pagina;
  }, [pagina]);
  useEffect(() => {
    if (pdf && ancho) desplazarA(paginaRef.current, false);
  }, [pdf, zoom, ancho, desplazarA]);

  // Si quien lo usa pide otra página (otra coincidencia), el visor salta a ella sin recargar el documento.
  const [paginaPedida, setPaginaPedida] = useState(paginaInicial);
  if (paginaPedida !== paginaInicial) {
    setPaginaPedida(paginaInicial);
    setPagina(paginaInicial);
    setCampo(String(paginaInicial));
  }
  const [saltoVisto, setSaltoVisto] = useState(salto);
  if (saltoVisto !== salto) {
    setSaltoVisto(salto);
    setPagina(paginaInicial);
    setCampo(String(paginaInicial));
  }
  useEffect(() => {
    if (pdf) desplazarA(paginaInicial, true);
  }, [paginaInicial, salto, pdf, desplazarA]);

  // Página actual = la que ocupa el tercio superior de la vista mientras se hace scroll.
  const alDesplazar = () => {
    const raiz = escenario.current;
    if (!raiz) return;
    const linea = raiz.scrollTop + raiz.clientHeight / 3;
    let actual = 1;
    for (let n = 1; n <= total; n++) {
      const el = paginasDom.current.get(n);
      if (el && el.offsetTop <= linea) actual = n;
      else if (el) break;
    }
    if (actual !== pagina) {
      setPagina(actual);
      setCampo(String(actual));
    }
  };

  const terminos = useMemo(
    () =>
      sinTildes(consulta)
        .split(/\s+/)
        .filter((t) => t.length >= 2),
    [consulta],
  );

  const hits = useMemo(() => [...new Set(paginasConCoincidencia)].sort((a, b) => a - b), [paginasConCoincidencia]);
  const posicionHit = hits.indexOf(pagina);
  const siguienteHit = hits.find((p) => p > pagina) ?? hits[0];
  const anteriorHit = [...hits].reverse().find((p) => p < pagina) ?? hits[hits.length - 1];

  // Miniaturas: la página actual y sus vecinas. Las páginas con coincidencia se recorren con las flechas de la
  // barra y con la lista lateral: dibujar una miniatura por cada una dejaba al trabajador de pdf.js ocupado con
  // decenas de páginas antes de mostrar la que se buscó.
  const enRiel = useMemo(() => {
    const s: number[] = [];
    for (let n = pagina - 1; n <= pagina + 2; n++) if (n >= 1 && n <= total) s.push(n);
    return s;
  }, [pagina, total]);

  const porcentaje = descarga ? Math.floor((descarga.cargado / descarga.total) * 100) : 0;
  const descargando = !esImagen && !error && (cargando || (descarga !== null && descarga.cargado < descarga.total));
  const mb = (b: number) => (b / (1024 * 1024)).toFixed(1).replace(".", ",");

  return (
    <div
      className="flex flex-col gap-3 outline-none"
      tabIndex={0}
      onKeyDown={(e) => {
        if ((e.target as HTMLElement).tagName === "INPUT") return;
        if (e.key === "ArrowRight" || e.key === "PageDown") {
          e.preventDefault();
          ir(pagina + 1);
        }
        if (e.key === "ArrowLeft" || e.key === "PageUp") {
          e.preventDefault();
          ir(pagina - 1);
        }
      }}
    >
      {!esImagen && (
        <div className="flex flex-wrap items-center justify-center gap-x-3 gap-y-1.5 self-center rounded-xl bg-sidebar px-3 py-2 text-white shadow-[0_8px_24px_rgba(15,27,46,0.22)]">
          <div className="flex items-center gap-1">
            <button type="button" className={botonHerramienta} onClick={() => ir(pagina - 1)} disabled={pagina <= 1} aria-label="Página anterior">
              <Icono d={["m15 18-6-6 6-6"]} />
            </button>
            <form
              className="flex items-center gap-1.5 text-sm font-semibold"
              onSubmit={(e) => {
                e.preventDefault();
                ir(Number(campo));
              }}
            >
              <input
                value={campo}
                onChange={(e) => setCampo(e.target.value.replace(/\D/g, ""))}
                onBlur={() => ir(Number(campo))}
                inputMode="numeric"
                aria-label="Número de página"
                className="h-8 w-12 rounded-lg bg-white/10 text-center text-white outline-none focus:ring-2 focus:ring-white/40"
              />
              <span className="text-white/60">/ {total || "—"}</span>
            </form>
            <button type="button" className={botonHerramienta} onClick={() => ir(pagina + 1)} disabled={pagina >= total} aria-label="Página siguiente">
              <Icono d={["m9 18 6-6-6-6"]} />
            </button>
          </div>
          <span className="hidden h-5 w-px bg-white/20 sm:block" aria-hidden="true" />
          <div className="flex items-center gap-1">
            <button type="button" className={botonHerramienta} onClick={() => setZoom((z) => Math.max(0.5, +(z - 0.25).toFixed(2)))} disabled={zoom <= 0.5} aria-label="Alejar">
              <Icono d={["M5 12h14"]} />
            </button>
            <button type="button" onClick={() => setZoom(1)} className="min-w-[52px] rounded-lg px-1 py-1.5 text-center text-sm font-semibold hover:bg-white/10" aria-label="Ajustar al ancho">
              {Math.round(zoom * 100)}%
            </button>
            <button type="button" className={botonHerramienta} onClick={() => setZoom((z) => Math.min(3, +(z + 0.25).toFixed(2)))} disabled={zoom >= 3} aria-label="Acercar">
              <Icono d={["M12 5v14", "M5 12h14"]} />
            </button>
          </div>
          {hits.length > 0 && (
            <>
              <span className="hidden h-5 w-px bg-white/20 sm:block" aria-hidden="true" />
              <div className="flex items-center gap-1.5 text-sm">
                <span className="max-w-[160px] truncate text-white/90">«{consulta}»</span>
                <span className="text-white/60">{posicionHit >= 0 ? `${posicionHit + 1} de ${hits.length}` : `${hits.length} páginas`}</span>
                <button type="button" className={botonHerramienta} onClick={() => ir(anteriorHit)} aria-label="Coincidencia anterior">
                  <Icono d={["m6 15 6-6 6 6"]} />
                </button>
                <button type="button" className={botonHerramienta} onClick={() => ir(siguienteHit)} aria-label="Coincidencia siguiente">
                  <Icono d={["m6 9 6 6 6-6"]} />
                </button>
              </div>
            </>
          )}
        </div>
      )}

      {descargando && (
        <div className="flex flex-col gap-1.5" role="status" aria-live="polite">
          <div className="flex items-baseline justify-between gap-3 text-xs text-ink-soft">
            <span>{descarga ? "Descargando el documento…" : "Abriendo el documento…"}</span>
            {descarga && (
              <span className="text-ink-faint">
                {porcentaje}% · {mb(descarga.cargado)} de {mb(descarga.total)} MB
              </span>
            )}
          </div>
          <progress value={descarga ? descarga.cargado : undefined} max={descarga ? descarga.total : undefined} className="h-1.5 w-full" aria-label="Descarga del documento" />
          {descarga && !cargando && <span className="text-xs text-ink-faint">Ya puedes leerlo; las búsquedas y los saltos de página serán inmediatos cuando termine.</span>}
        </div>
      )}

      <div className="flex min-w-0 gap-3">
        {!esImagen && pdf && total > 1 && (
          <nav aria-label="Páginas" className="hidden w-[104px] shrink-0 flex-col items-center gap-3 self-start rounded-xl border border-line bg-surface px-2 py-3 md:flex">
            {enRiel.map((n) => (
              <Miniatura key={n} pdf={pdf} n={n} activa={n === pagina} coincide={hits.includes(n)} alElegir={() => ir(n)} />
            ))}
          </nav>
        )}

        <div
          ref={escenario}
          onScroll={alDesplazar}
          className="relative min-h-[560px] min-w-0 flex-1 overflow-auto rounded-xl border border-line bg-bg p-6 dark:bg-black/20"
          style={{ maxHeight: "78vh" }}
        >
          {cargando && (
            <div className="py-16" role="status" aria-label="Abriendo documento">
              <Loader compacto />
            </div>
          )}
          {error && <p className="py-20 text-center text-sm text-accent-ink">{error}</p>}
          {esImagen ? (
            // eslint-disable-next-line @next/next/no-img-element -- imagen del propio expediente servida con sesión
            <img src={src} alt={nombre} className="mx-auto h-auto max-w-full rounded-[4px] bg-white shadow-[0_1px_3px_rgba(15,27,46,0.12),0_18px_40px_rgba(15,27,46,0.12)]" />
          ) : (
            pdf &&
            pdfjs &&
            ancho > 0 && (
              <div className="flex flex-col items-center" style={{ gap: SEPARACION }}>
                {medidas.map((m, i) => (
                  <Pagina
                    key={i + 1}
                    pdf={pdf}
                    pdfjs={pdfjs}
                    n={i + 1}
                    medida={m}
                    escala={(ancho / m.w) * zoom}
                    dibujar={cerca.has(i + 1)}
                    terminos={terminos}
                    nombre={nombre}
                    alMedir={alMedir}
                    refPagina={refPagina}
                  />
                ))}
              </div>
            )
          )}
        </div>
      </div>
    </div>
  );
}
