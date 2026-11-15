"use client";

import { useSyncExternalStore } from "react";
import { CLAVE_TEMA } from "@/lib/tema";

const EVENTO = "escalafon:tema-cambiado";

function suscribirse(avisar: () => void) {
  window.addEventListener(EVENTO, avisar);
  return () => window.removeEventListener(EVENTO, avisar);
}

const esOscuro = () => document.documentElement.classList.contains("dark");

// Alterna modo claro/oscuro y lo guarda en localStorage. El script de app/layout.tsx
// lo aplica antes de pintar para que no haya parpadeo al recargar.
export default function BotonTema({ className = "" }: { className?: string }) {
  const oscuro = useSyncExternalStore(suscribirse, esOscuro, () => false);

  const alternar = () => {
    const siguiente = !oscuro;
    document.documentElement.classList.toggle("dark", siguiente);
    try {
      localStorage.setItem(CLAVE_TEMA, siguiente ? "oscuro" : "claro");
    } catch {}
    window.dispatchEvent(new Event(EVENTO));
  };

  return (
    <button
      type="button"
      onClick={alternar}
      aria-label={oscuro ? "Cambiar a modo claro" : "Cambiar a modo oscuro"}
      title={oscuro ? "Modo claro" : "Modo oscuro"}
      className={`flex h-8 w-8 items-center justify-center rounded-lg border border-line text-ink-soft transition-colors duration-100 hover:border-accent hover:text-accent-ink ${className}`}
    >
      {oscuro ? (
        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden="true">
          <circle cx="12" cy="12" r="4" />
          <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
        </svg>
      ) : (
        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
          <path d="M20.4 15.4A9 9 0 0 1 8.6 3.6 9 9 0 1 0 20.4 15.4Z" />
        </svg>
      )}
    </button>
  );
}
