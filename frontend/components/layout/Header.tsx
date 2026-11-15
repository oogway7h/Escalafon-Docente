"use client";

import { cerrarSesion } from "@/app/panel/actions";
import BotonTema from "@/components/common/BotonTema";
import type { Sesion } from "@/lib/api";
import { IconoColapsar, IconoMenuMovil, IconoSalir } from "./IconosPanel";

export default function Header({
  sesion,
  colapsado,
  onAlternarColapso,
  onAbrirMovil,
}: {
  sesion: Sesion;
  colapsado: boolean;
  onAlternarColapso: () => void;
  onAbrirMovil: () => void;
}) {
  return (
    <header className="sticky top-0 z-30 flex items-center justify-between gap-4 border-b border-line bg-surface px-5 py-3">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={onAbrirMovil}
          className="flex h-8 w-8 items-center justify-center rounded-lg text-ink-soft hover:bg-bg md:hidden"
          aria-label="Abrir menú"
        >
          <IconoMenuMovil className="h-[18px] w-[18px]" />
        </button>
        <button
          type="button"
          onClick={onAlternarColapso}
          className="hidden h-8 w-8 items-center justify-center rounded-lg text-ink-soft hover:bg-bg md:flex"
          aria-label={colapsado ? "Expandir menú" : "Contraer menú"}
        >
          <IconoColapsar className="h-[18px] w-[18px]" colapsado={colapsado} />
        </button>
        <div className="hidden flex-col sm:flex">
          <p className="text-sm font-semibold text-ink">Hola, {sesion.usuario.nombres}</p>
          <p className="text-xs text-ink-faint">{sesion.rol}</p>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <BotonTema />
        <form action={cerrarSesion}>
          <button
            type="submit"
            className="flex items-center gap-1.5 rounded-lg border border-line px-3 py-1.5 text-sm font-semibold text-ink-soft transition-colors duration-100 hover:border-accent hover:text-accent-ink"
          >
            <IconoSalir className="h-[15px] w-[15px]" />
            Salir
          </button>
        </form>
      </div>
    </header>
  );
}
