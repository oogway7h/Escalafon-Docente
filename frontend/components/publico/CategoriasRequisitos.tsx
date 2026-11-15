"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";

type CategoriaInfo = { id: number; nombre: string; aniosAntiguedad: number; requisitos: { id: number; nombre: string }[] };

// Tarjetas de dos caras: con ratón giran al pasar el cursor; al tocar o hacer clic se BLOQUEAN giradas
// (en móvil el toque es el único gesto) y otro clic las devuelve. Los datos son los del catálogo que
// administra el Director.
export default function CategoriasRequisitos({ categorias }: { categorias: CategoriaInfo[] }) {
  return (
    <div className="grid grid-cols-1 gap-7 sm:grid-cols-2 lg:grid-cols-3">
      {categorias.map((c, i) => (
        <CartaCategoria key={c.id} categoria={c} indice={i} total={categorias.length} />
      ))}
    </div>
  );
}

function CartaCategoria({ categoria: c, indice, total }: { categoria: CategoriaInfo; indice: number; total: number }) {
  const [bloqueada, setBloqueada] = useState(false);
  const [encima, setEncima] = useState(false);
  const volteada = bloqueada || encima;
  const verRef = useRef<HTMLButtonElement>(null);
  const cerrarRef = useRef<HTMLButtonElement>(null);
  const primera = useRef(true);

  // Al bloquear o desbloquear con el teclado el botón pulsado queda inerte (cambia de cara): el foco pasa al de la otra cara.
  useEffect(() => {
    if (primera.current) {
      primera.current = false;
      return;
    }
    (bloqueada ? cerrarRef : verRef).current?.focus({ preventScroll: true });
  }, [bloqueada]);

  const alternar = (e: React.MouseEvent) => {
    if ((e.target as HTMLElement).closest("a")) return; // el enlace de la cara trasera navega
    setBloqueada((b) => !b);
    setEncima(false); // si el cursor sigue encima, no se vuelve a girar hasta que salga y vuelva a entrar
  };

  return (
    <div
      onClick={alternar}
      onPointerEnter={(e) => e.pointerType === "mouse" && setEncima(true)}
      onPointerLeave={(e) => e.pointerType === "mouse" && setEncima(false)}
      className={`h-[400px] cursor-pointer transition-transform duration-300 [perspective:1400px] motion-reduce:transition-none ${volteada && !bloqueada ? "-translate-y-1" : ""}`}
    >
      <div
        className={`relative h-full w-full transition-transform duration-700 ease-[cubic-bezier(.2,.8,.2,1)] [transform-style:preserve-3d] motion-reduce:transition-none ${volteada ? "[transform:rotateY(180deg)]" : ""}`}
      >
        {/* Cara frontal */}
        <div
          inert={volteada}
          className="absolute inset-0 flex flex-col justify-between overflow-hidden rounded-[22px] border-[1.5px] border-line bg-surface p-7 shadow-[0_18px_40px_-26px_rgba(20,33,61,0.45)] [backface-visibility:hidden]"
        >
          <div className="flex flex-col gap-2.5">
            <div className="flex items-center justify-between">
              <span className="text-[11.5px] font-bold uppercase tracking-[0.12em] text-ink-faint">
                Categoría {indice + 1} de {total}
              </span>
              <div className="flex h-[26px] items-end gap-1" aria-hidden="true">
                {Array.from({ length: total }, (_, k) => (
                  <span
                    key={k}
                    className={`block w-2 rounded-[3px] transition-transform duration-500 ${k <= indice ? "bg-accent" : "bg-line"} ${encima ? "scale-y-110" : ""}`}
                    style={{ height: 8 + k * (18 / Math.max(total - 1, 1)), transitionDelay: `${k * 40}ms` }}
                  />
                ))}
              </div>
            </div>
            <h3 className="font-display text-[28px] font-extrabold leading-[1.1] tracking-[-0.02em] text-ink">{c.nombre}</h3>
          </div>

          {c.aniosAntiguedad > 0 ? (
            <div className="flex items-baseline gap-2.5">
              <span className="font-mono text-[84px] font-medium leading-[0.9] tracking-[-0.04em] text-accent-ink">{c.aniosAntiguedad}</span>
              <span className="text-[15px] font-semibold leading-snug text-ink-soft">
                años de
                <br />
                antigüedad mínima
              </span>
            </div>
          ) : (
            <p className="text-[18px] font-bold leading-snug text-ink-soft">Sin antigüedad mínima</p>
          )}

          <div className="flex items-center justify-between gap-3">
            <span className="rounded-full bg-bg px-3.5 py-2 text-[13.5px] font-bold text-ink">
              {c.requisitos.length} {c.requisitos.length === 1 ? "documento" : "documentos"}
            </span>
            <button
              ref={verRef}
              type="button"
              aria-pressed={bloqueada}
              aria-label={`Ver los requisitos de ${c.nombre}`}
              className="inline-flex min-h-11 items-center gap-2 rounded-[10px] border-[1.5px] border-line bg-surface px-4 text-[14px] font-bold text-ink transition-colors hover:border-accent hover:text-accent-ink"
            >
              Ver requisitos
              <IconoGiro />
            </button>
          </div>
        </div>

        {/* Cara trasera */}
        <div
          inert={!volteada}
          className="absolute inset-0 flex flex-col gap-4 overflow-hidden rounded-[22px] p-7 text-white shadow-[0_18px_40px_-22px_rgba(20,33,61,0.7)] [backface-visibility:hidden] [transform:rotateY(180deg)]"
          style={{ background: "var(--color-sidebar)" }}
        >
          <div className="flex items-start justify-between gap-3">
            <div>
              <p className="text-[11.5px] font-bold uppercase tracking-[0.12em] text-white/60">Documentos exigidos</p>
              <h3 className="mt-1.5 font-display text-[22px] font-extrabold leading-tight">{c.nombre}</h3>
            </div>
            <button
              ref={cerrarRef}
              type="button"
              aria-label={`Volver a la tarjeta de ${c.nombre}`}
              className="flex h-11 w-11 shrink-0 items-center justify-center rounded-[10px] border-[1.5px] border-white/25 text-white transition-colors hover:bg-white hover:text-[#14213d]"
            >
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path d="M6 6l12 12M18 6L6 18" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" />
              </svg>
            </button>
          </div>

          {c.requisitos.length === 0 ? (
            <p className="flex-1 text-[14.5px] text-white/75">Esta categoría no exige documentos.</p>
          ) : (
            <ul className="flex flex-1 flex-col gap-2.5 overflow-y-auto pr-1">
              {c.requisitos.map((r) => (
                <li key={r.id} className="flex items-center gap-3 text-[15px] font-semibold leading-snug">
                  <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent">
                    <svg width="13" height="13" viewBox="0 0 14 14" fill="none" aria-hidden="true">
                      <path d="M2.5 7.5l3 3 6-7" stroke="#fff" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />
                    </svg>
                  </span>
                  {r.nombre}
                </li>
              ))}
            </ul>
          )}

          <Link
            href="/solicitud"
            className="inline-flex min-h-11 w-fit items-center rounded-[10px] bg-accent px-5 text-[14.5px] font-extrabold text-white transition-[transform,background] duration-100 hover:bg-accent-dark active:scale-[0.975]"
          >
            Iniciar mi solicitud
          </Link>
        </div>
      </div>
    </div>
  );
}

function IconoGiro() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path d="M4 12a8 8 0 0114-5.3M20 12a8 8 0 01-14 5.3M18 3v4h-4M6 21v-4h4" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}
