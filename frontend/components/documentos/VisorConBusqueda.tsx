"use client";

import { useState } from "react";
import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import { claseControl } from "@/components/common/Campo";
import VisorPdf from "./VisorPdf";

export type Coincidencia = { pagina: number; fragmento: string };

// El servicio marca las palabras halladas entre ⟦ y ⟧.
export function Fragmento({ texto }: { texto: string }) {
  return texto.split(/⟦(.*?)⟧/).map((t, i) =>
    i % 2 ? (
      <mark key={i} className="rounded-[3px] bg-accent-soft px-0.5 font-semibold text-accent-dark">
        {t}
      </mark>
    ) : (
      t
    ),
  );
}

// Refleja la búsqueda y la página en la URL (para compartir o recargar) sin navegar: el documento ya está
// abierto en el navegador y una navegación volvería a pedir toda la página al servidor.
function reflejarEnUrl(cambios: Record<string, string>) {
  const url = new URL(window.location.href);
  for (const [k, v] of Object.entries(cambios)) {
    if (v) url.searchParams.set(k, v);
    else url.searchParams.delete(k);
  }
  window.history.replaceState(null, "", url);
}

// Visor + panel lateral (buscador interno, coincidencias y, como children, los detalles del documento).
// Buscar y saltar entre coincidencias ocurre en el navegador: la búsqueda pide solo las páginas halladas a
// `rutaBuscar` y el salto mueve el visor, sin volver a renderizar la página en el servidor.
export default function VisorConBusqueda({
  src,
  esImagen,
  nombre,
  paginaInicial,
  consultaInicial,
  coincidenciasIniciales,
  rutaBuscar,
  children,
}: {
  src: string;
  esImagen: boolean;
  nombre: string;
  paginaInicial: number;
  consultaInicial: string;
  coincidenciasIniciales: Coincidencia[];
  rutaBuscar?: string; // sin ella no hay buscador interno (la búsqueda la hace la página que lo contiene)
  children?: React.ReactNode;
}) {
  const [pagina, setPagina] = useState(paginaInicial);
  const [salto, setSalto] = useState(0);
  const [consulta, setConsulta] = useState(consultaInicial);
  const [halladas, setHalladas] = useState(coincidenciasIniciales);
  const [buscando, setBuscando] = useState(false);
  const [error, setError] = useState("");

  // Si la página que lo contiene busca de nuevo (otra consulta u otra página pedida), se adopta lo que envía.
  const firma = `${consultaInicial}|${paginaInicial}|${coincidenciasIniciales.length}`;
  const [firmaVista, setFirmaVista] = useState(firma);
  if (firmaVista !== firma) {
    setFirmaVista(firma);
    setConsulta(consultaInicial);
    setHalladas(coincidenciasIniciales);
    setPagina(paginaInicial);
    setSalto((n) => n + 1);
  }

  function irA(n: number) {
    setPagina(n);
    setSalto((s) => s + 1);
    reflejarEnUrl({ pagina: String(n) });
  }

  async function buscar(formData: FormData) {
    const q = String(formData.get("q") ?? "").trim();
    setError("");
    if (!q) {
      setConsulta("");
      setHalladas([]);
      reflejarEnUrl({ q: "" });
      return;
    }
    setBuscando(true);
    try {
      const r = await fetch(`${rutaBuscar}?q=${encodeURIComponent(q)}`);
      const cuerpo = await r.json().catch(() => null);
      if (!r.ok) {
        setError(cuerpo?.error ?? "No se pudo buscar en el documento.");
        return;
      }
      const lista: Coincidencia[] = cuerpo?.coincidencias ?? [];
      setConsulta(q);
      setHalladas(lista);
      reflejarEnUrl({ q });
      if (lista.length > 0) irA(lista[0].pagina);
    } catch {
      setError("No se pudo buscar en el documento. Revisa tu conexión.");
    } finally {
      setBuscando(false);
    }
  }

  return (
    <div className="flex flex-wrap items-start gap-5">
      <div className="min-w-0 flex-[999_1_520px]">
        <VisorPdf
          src={src}
          esImagen={esImagen}
          nombre={nombre}
          paginaInicial={pagina}
          salto={salto}
          consulta={consulta}
          paginasConCoincidencia={halladas.map((c) => c.pagina)}
        />
      </div>

      <aside className="flex min-w-[260px] flex-1 flex-col gap-4 lg:max-w-[340px]">
        {rutaBuscar && (
          <form action={buscar} className="flex gap-2">
            <input name="q" defaultValue={consulta} placeholder="Buscar en este documento" aria-label="Buscar en el documento" className={`${claseControl()} min-w-0 flex-1`} />
            <button type="submit" disabled={buscando} className={claseBoton("secundario")}>
              {buscando ? "Buscando…" : "Buscar"}
            </button>
          </form>
        )}
        {error && <Alerta tipo="error">{error}</Alerta>}

        {consulta && !error && (
          <div className="rounded-xl border border-line p-4">
            <h3 className="mb-2 text-sm font-semibold text-ink">{halladas.length === 0 ? "Sin coincidencias" : `Coincidencias · ${halladas.length}`}</h3>
            <ul className="flex max-h-[52vh] flex-col gap-1 overflow-y-auto">
              {halladas.map((c) => (
                <li key={c.pagina}>
                  <button
                    type="button"
                    onClick={() => irA(c.pagina)}
                    aria-current={c.pagina === pagina ? "true" : undefined}
                    className={`flex w-full flex-col gap-0.5 rounded-lg px-3 py-2 text-left hover:bg-bg ${c.pagina === pagina ? "bg-accent-soft" : ""}`}
                  >
                    <span className={`text-xs font-semibold ${c.pagina === pagina ? "text-accent-dark" : "text-ink-soft"}`}>Pág. {c.pagina}</span>
                    <span className="text-sm leading-snug text-ink-soft">
                      <Fragmento texto={c.fragmento} />
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}

        {children}
      </aside>
    </div>
  );
}
