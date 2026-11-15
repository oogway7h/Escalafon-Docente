"use client";

import { useRouter } from "next/navigation";
import { useEffect, useId, useRef, useState, useTransition } from "react";
import BotonCarga from "./BotonCarga";
import { useCargaListado } from "./CargaListado";

// Barra de búsqueda con sugerencias debajo mientras se escribe, como la de un navegador: ↑/↓ recorren
// las sugerencias, Enter abre la marcada (o busca el texto) y Esc las cierra. La lista sale de `fuente`
// (un route handler que responde Sugerencia[]); las respuestas se recuerdan para que borrar sea instantáneo.

export type Sugerencia = {
  clave: string;
  texto: string; // línea principal; se resalta lo escrito
  detalle?: string; // segunda línea (código, categoría, página…)
  fragmento?: string; // texto con las palabras halladas entre ⟦ ⟧
  href: string;
};

const sinTildes = (s: string) => s.normalize("NFD").replace(/[̀-ͯ]/g, "").toLowerCase();

// Resalta lo escrito dentro del texto, sin distinguir tildes ni mayúsculas.
function Resaltado({ texto, consulta }: { texto: string; consulta: string }) {
  const normal = [...texto].map((c) => sinTildes(c)[0] ?? c).join("");
  const q = sinTildes(consulta.trim());
  const i = q && normal.length === texto.length ? normal.indexOf(q) : -1;
  if (i < 0) return <>{texto}</>;
  return (
    <>
      {texto.slice(0, i)}
      <strong className="font-semibold text-ink">{texto.slice(i, i + q.length)}</strong>
      {texto.slice(i + q.length)}
    </>
  );
}

function Fragmento({ texto }: { texto: string }) {
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

export default function BuscadorSugerencias({
  accion,
  nombre = "q",
  ocultos = {},
  valorInicial = "",
  fuente,
  placeholder,
  etiqueta,
  textoBoton = "Buscar",
  autoFocus = false,
}: {
  accion: string; // página que recibe la búsqueda completa (Enter sin sugerencia marcada)
  nombre?: string; // parámetro de la búsqueda en la URL
  ocultos?: Record<string, string>; // otros parámetros que se conservan
  valorInicial?: string;
  fuente: string; // URL que devuelve Sugerencia[]; se le agrega q=
  placeholder: string;
  etiqueta: string;
  textoBoton?: string;
  autoFocus?: boolean;
}) {
  const router = useRouter();
  const listado = useCargaListado();
  const [propio, iniciar] = useTransition();
  const pendiente = listado?.pendiente ?? propio;
  const navegar = (href: string) => (listado ? listado.navegar(href) : iniciar(() => router.push(href)));

  const id = useId();
  const [texto, setTexto] = useState(valorInicial);
  const [sugerencias, setSugerencias] = useState<Sugerencia[]>([]);
  const [abierta, setAbierta] = useState(false);
  const [activa, setActiva] = useState(-1);
  const [buscando, setBuscando] = useState(false);
  const [fallo, setFallo] = useState(false); // la fuente respondió con error: no es lo mismo que «sin coincidencias»
  const memo = useRef(new Map<string, Sugerencia[]>());
  const escrito = useRef(false); // solo se sugiere tras escribir, no con el valor inicial

  useEffect(() => {
    const q = texto.trim();
    if (!escrito.current || q.length < 2) return;
    const guardadas = memo.current.get(q);
    if (guardadas) {
      setSugerencias(guardadas);
      setFallo(false);
      return;
    }
    const control = new AbortController();
    const espera = setTimeout(async () => {
      setBuscando(true);
      try {
        const r = await fetch(`${fuente}${fuente.includes("?") ? "&" : "?"}q=${encodeURIComponent(q)}`, { signal: control.signal });
        if (!r.ok) throw new Error(`HTTP ${r.status}`);
        const lista: Sugerencia[] = await r.json();
        memo.current.set(q, lista); // solo se recuerdan las respuestas buenas
        setSugerencias(lista);
        setFallo(false);
        setActiva(-1);
      } catch {
        // abortada = la reemplazó otra búsqueda más reciente; cualquier otro error se avisa
        if (!control.signal.aborted) {
          setSugerencias([]);
          setFallo(true);
        }
      } finally {
        if (!control.signal.aborted) setBuscando(false);
      }
    }, 120);
    return () => {
      clearTimeout(espera);
      control.abort();
    };
  }, [texto, fuente]);

  const q = texto.trim();
  const visibles = q.length >= 2 ? sugerencias : [];
  const mostrar = abierta && q.length >= 2 && (visibles.length > 0 || !buscando);

  const elegir = (s: Sugerencia) => {
    setAbierta(false);
    setTexto(s.texto);
    navegar(s.href);
  };
  const enviar = () => {
    setAbierta(false);
    const qs = new URLSearchParams(ocultos);
    if (q) qs.set(nombre, q);
    navegar(`${accion}?${qs}`);
  };

  return (
    <form
      role="search"
      action={accion}
      className="relative flex flex-col"
      onSubmit={(e) => {
        e.preventDefault();
        if (activa >= 0 && visibles[activa]) elegir(visibles[activa]);
        else enviar();
      }}
    >
      {Object.entries(ocultos).map(([k, v]) => (
        <input key={k} type="hidden" name={k} value={v} />
      ))}
      <div
        className={`flex h-16 items-center gap-3 border border-line bg-surface pl-5 pr-2.5 shadow-[0_1px_2px_rgba(15,27,46,0.04),0_8px_24px_rgba(15,27,46,0.06)] focus-within:border-accent focus-within:ring-[3px] focus-within:ring-accent/20 ${mostrar ? "rounded-t-xl" : "rounded-xl"}`}
      >
        <svg className="shrink-0 text-ink-soft" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
          <circle cx="11" cy="11" r="7" />
          <path d="m20 20-3.5-3.5" />
        </svg>
        <input
          type="search"
          name={nombre}
          value={texto}
          autoFocus={autoFocus}
          autoComplete="off"
          spellCheck={false}
          data-foco-en-contenedor
          placeholder={placeholder}
          aria-label={etiqueta}
          role="combobox"
          aria-expanded={mostrar}
          aria-controls={`${id}-lista`}
          aria-autocomplete="list"
          aria-activedescendant={activa >= 0 ? `${id}-${activa}` : undefined}
          onChange={(e) => {
            escrito.current = true;
            setTexto(e.target.value);
            setAbierta(true);
            setActiva(-1);
          }}
          onFocus={() => setAbierta(true)}
          onBlur={() => setAbierta(false)}
          onKeyDown={(e) => {
            if (e.key === "ArrowDown" && visibles.length) {
              e.preventDefault();
              setAbierta(true);
              setActiva((a) => (a + 1) % visibles.length);
            } else if (e.key === "ArrowUp" && visibles.length) {
              e.preventDefault();
              setActiva((a) => (a <= 0 ? visibles.length - 1 : a - 1));
            } else if (e.key === "Escape") {
              setAbierta(false);
              setActiva(-1);
            }
          }}
          className="h-full min-w-0 flex-1 bg-transparent text-xl text-ink outline-none placeholder:text-ink-faint"
        />
        <BotonCarga type="submit" cargando={pendiente} disabled={pendiente} clase="inline-flex h-11 shrink-0 items-center justify-center rounded-xl bg-accent px-6 text-base font-semibold text-white hover:bg-accent-dark disabled:cursor-wait">
          {textoBoton}
        </BotonCarga>
      </div>

      {mostrar && (
        <ul
          id={`${id}-lista`}
          role="listbox"
          aria-label="Sugerencias"
          className="absolute left-0 right-0 top-full z-30 max-h-[60vh] overflow-y-auto rounded-b-xl border border-t-0 border-line bg-surface py-1.5 shadow-[0_18px_40px_rgba(15,27,46,0.14)]"
        >
          {visibles.length === 0 && (
            <li className="px-5 py-3 text-sm text-ink-faint">
              {fallo ? "No se pudieron cargar las sugerencias. Pulsa Buscar para ver los resultados." : `Sin coincidencias para «${q}». Pulsa Buscar para ver todos los resultados.`}
            </li>
          )}
          {visibles.map((s, i) => (
            <li
              key={s.clave}
              id={`${id}-${i}`}
              role="option"
              aria-selected={i === activa}
              // mousedown en vez de click: el blur del campo cerraría la lista antes del click.
              onMouseDown={(e) => {
                e.preventDefault();
                elegir(s);
              }}
              onMouseEnter={() => setActiva(i)}
              className={`flex cursor-pointer items-start gap-3 px-5 py-2.5 ${i === activa ? "bg-bg" : ""}`}
            >
              <svg className="mt-1 shrink-0 text-ink-faint" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                <circle cx="11" cy="11" r="7" />
                <path d="m20 20-3.5-3.5" />
              </svg>
              <span className="flex min-w-0 flex-col">
                <span className="truncate text-base text-ink-soft">
                  <Resaltado texto={s.texto} consulta={q} />
                </span>
                {s.detalle && <span className="text-xs text-ink-faint">{s.detalle}</span>}
                {s.fragmento && (
                  <span className="line-clamp-2 text-sm leading-snug text-ink-soft">
                    <Fragmento texto={s.fragmento} />
                  </span>
                )}
              </span>
            </li>
          ))}
        </ul>
      )}
    </form>
  );
}
