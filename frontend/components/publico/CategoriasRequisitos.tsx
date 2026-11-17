"use client";

import Link from "next/link";
import { useRef, useState } from "react";
import { REQUISITOS_REGULARIZACION, textoPuntaje } from "@/lib/requisitos";

type CategoriaInfo = {
  id: number;
  nombre: string;
  puntajeDesde: number | null;
  puntajeHasta: number | null;
  requisitos: { id: number; nombre: string }[];
};

// La regularización no es una categoría del catálogo: es un trámite aparte con su propia lista de documentos.
const REGULARIZACION: CategoriaInfo = {
  id: -1,
  nombre: "Regularización docente",
  puntajeDesde: null,
  puntajeHasta: null,
  requisitos: REQUISITOS_REGULARIZACION.map((nombre, i) => ({ id: -(i + 1), nombre })),
};

// Tarjetas de dos caras: con ratón giran al pasar el cursor; al tocar o hacer clic se BLOQUEAN giradas
// (en móvil el toque es el único gesto) y otro clic las devuelve. Los datos son los del catálogo que
// administra el Director.
export default function CategoriasRequisitos({ categorias }: { categorias: CategoriaInfo[] }) {
  // Solo las categorías con puntaje (Adjunto a Emérito): son las que se alcanzan por categorización.
  const conPuntaje = categorias.filter((c) => c.puntajeDesde !== null);
  return (
    <div className="grid grid-cols-1 gap-7 sm:grid-cols-2 lg:grid-cols-3">
      {conPuntaje.map((c, i) => (
        <CartaCategoria key={c.id} categoria={c} indice={i} total={conPuntaje.length} />
      ))}
      <CartaCategoria categoria={REGULARIZACION} indice={conPuntaje.length} total={conPuntaje.length} regularizacion />
    </div>
  );
}

function CartaCategoria({
  categoria: c,
  indice,
  total,
  regularizacion = false,
}: {
  categoria: CategoriaInfo;
  indice: number;
  total: number;
  regularizacion?: boolean;
}) {
  const puntaje = textoPuntaje(c.puntajeDesde, c.puntajeHasta);
  const [encima, setEncima] = useState(false);
  const volteada = encima;
  const dialogo = useRef<HTMLDialogElement>(null);
  const idTitulo = `requisitos-${c.id}`;

  // Con ratón la tarjeta gira al pasar el cursor; al hacer clic (o tocar, en móvil) se abre el modal con la lista.
  const abrir = (e: React.MouseEvent) => {
    if ((e.target as HTMLElement).closest("a")) return; // el enlace de la cara trasera navega
    setEncima(false);
    dialogo.current?.showModal();
  };
  const cerrar = () => dialogo.current?.close();

  return (
    <>
    <div
      onClick={abrir}
      onPointerEnter={(e) => e.pointerType === "mouse" && setEncima(true)}
      onPointerLeave={(e) => e.pointerType === "mouse" && setEncima(false)}
      className={`h-[400px] cursor-pointer transition-transform duration-300 [perspective:1400px] motion-reduce:transition-none ${volteada ? "-translate-y-1" : ""}`}
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
                {regularizacion ? "Trámite aparte" : `Categoría ${indice + 1} de ${total}`}
              </span>
              <div className="flex h-[26px] items-end gap-1" aria-hidden="true">
                {Array.from({ length: regularizacion ? 0 : total }, (_, k) => (
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

          {regularizacion ? (
            <div className="flex items-baseline gap-2.5">
              <span className="font-mono text-[84px] font-medium leading-[0.9] tracking-[-0.04em] text-accent-ink">2</span>
              <span className="text-[15px] font-semibold leading-snug text-ink-soft">
                años o 4 semestres
                <br />
                sin cortes en la materia
              </span>
            </div>
          ) : puntaje ? (
            <div className="flex items-baseline gap-2.5">
              <span className="font-mono text-[56px] font-medium leading-[0.95] tracking-[-0.04em] text-accent-ink">{puntaje}</span>
              <span className="text-[15px] font-semibold leading-snug text-ink-soft">
                puntos de
                <br />
                categoría
              </span>
            </div>
          ) : null}

          <div className="flex items-center justify-between gap-3">
            <span className="rounded-full bg-bg px-3.5 py-2 text-[13.5px] font-bold text-ink">
              {c.requisitos.length} {c.requisitos.length === 1 ? "documento" : "documentos"}
            </span>
            <button
              type="button"
              aria-haspopup="dialog"
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
            <span className="shrink-0 rounded-full border border-white/25 px-3 py-1.5 text-[12px] font-bold text-white/80">Clic para ampliar</span>
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

      <dialog
        ref={dialogo}
        onClick={(e) => e.target === e.currentTarget && cerrar()} // clic en el fondo oscuro
        aria-labelledby={idTitulo}
        className="m-auto max-h-[calc(100vh-32px)] w-[min(640px,calc(100vw-32px))] rounded-2xl border border-line bg-surface p-0 text-ink backdrop:bg-black/30 backdrop:backdrop-blur-sm"
      >
        <div className="flex max-h-[calc(100vh-32px)] flex-col gap-5 p-6 sm:p-8">
          <header className="flex items-start justify-between gap-4">
            <div>
              <p className="text-[11.5px] font-bold uppercase tracking-[0.12em] text-ink-faint">Documentos exigidos</p>
              <h2 id={idTitulo} className="mt-1 font-display text-2xl font-extrabold tracking-[-0.01em]">
                {c.nombre}
              </h2>
              {puntaje && <p className="mt-1 text-sm font-semibold text-accent-ink">{puntaje} puntos de categoría</p>}
            </div>
            <button
              type="button"
              onClick={cerrar}
              aria-label="Cerrar"
              className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-line text-ink-soft hover:border-accent hover:text-accent-ink"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                <path d="M18 6 6 18" />
                <path d="m6 6 12 12" />
              </svg>
            </button>
          </header>

          {c.requisitos.length === 0 ? (
            <p className="text-sm text-ink-soft">Esta categoría no exige documentos.</p>
          ) : (
            <ol className="flex min-h-0 flex-col gap-3 overflow-y-auto pr-1">
              {c.requisitos.map((r, k) => (
                <li key={r.id} className="flex items-start gap-3 text-[15px] leading-snug">
                  <span className="mt-px flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent-soft text-[12px] font-bold text-accent-dark">{k + 1}</span>
                  {r.nombre}
                </li>
              ))}
            </ol>
          )}

          <Link
            href="/solicitud"
            className="inline-flex min-h-11 w-fit items-center rounded-lg bg-accent px-5 text-[14.5px] font-bold text-white transition-[transform,background] duration-100 hover:bg-accent-dark active:scale-[0.975]"
          >
            Iniciar mi solicitud
          </Link>
        </div>
      </dialog>
    </>
  );
}

function IconoGiro() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path d="M4 12a8 8 0 0114-5.3M20 12a8 8 0 01-14 5.3M18 3v4h-4M6 21v-4h4" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}
