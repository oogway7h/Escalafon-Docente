"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { createContext, useContext, useTransition } from "react";
import { claseBoton, type TamanoBoton, type VarianteBoton } from "./Boton";
import Loader from "./Loader";

// Navegación con la carga visible. Una fila clicable (FilaClicable) comparte su estado de espera con
// los EnlaceBoton y el IndicadorFila de su interior: haga clic en la fila o en el botón, el Loader
// aparece en el mismo sitio hasta que llega la página nueva.
const Espera = createContext<{ pendiente: boolean; ir: (href: string) => void } | null>(null);

function useNavegacion() {
  const router = useRouter();
  const [pendiente, iniciar] = useTransition();
  return { pendiente, ir: (href: string) => iniciar(() => router.push(href)) };
}

function esClicSimple(e: React.MouseEvent) {
  return !(e.defaultPrevented || e.metaKey || e.ctrlKey || e.shiftKey || e.button !== 0);
}

// Enlace con aspecto de botón que muestra el Loader dentro mientras navega.
export function EnlaceBoton({
  href,
  children,
  variante,
  tamano,
}: {
  href: string;
  children: React.ReactNode;
  variante?: VarianteBoton;
  tamano?: TamanoBoton;
}) {
  const propia = useNavegacion();
  const { pendiente, ir } = useContext(Espera) ?? propia;
  return (
    <Link
      href={href}
      aria-busy={pendiente}
      aria-disabled={pendiente}
      className={claseBoton(variante, tamano)}
      onClick={(e) => {
        if (!esClicSimple(e)) return;
        e.preventDefault();
        if (!pendiente) ir(href);
      }}
    >
      {pendiente ? (
        <span className="relative inline-flex items-center justify-center">
          <span className="invisible">{children}</span>
          <span className="absolute inset-0 flex items-center justify-center">
            <Loader mini />
          </span>
        </span>
      ) : (
        children
      )}
    </Link>
  );
}

// Fila de tabla que navega al hacer clic en cualquier parte (los enlaces y botones de dentro siguen
// funcionando por su cuenta). Con teclado se llega por el enlace del código, que sigue ahí.
export function FilaClicable({ href, children }: { href: string; children: React.ReactNode }) {
  const nav = useNavegacion();
  return (
    <Espera.Provider value={nav}>
      <tr
        aria-busy={nav.pendiente}
        onClick={(e) => {
          if ((e.target as HTMLElement).closest("a, button, input, select, textarea, label") || !esClicSimple(e)) return;
          if (!nav.pendiente) nav.ir(href);
        }}
        className={`cursor-pointer text-sm text-ink transition-colors hover:bg-bg ${nav.pendiente ? "cursor-progress opacity-60" : ""}`}
      >
        {children}
      </tr>
    </Espera.Provider>
  );
}

// Loader pequeño para filas sin botón: se muestra mientras la fila navega.
export function IndicadorFila() {
  const espera = useContext(Espera);
  if (!espera?.pendiente) return null;
  return (
    <span className="ml-2 inline-block align-middle">
      <Loader mini />
    </span>
  );
}

// Enlace de texto (p. ej. un documento) que navega sin recargar el scroll y muestra un Loader mini.
export function EnlaceSinScroll({ href, className, children }: { href: string; className?: string; children: React.ReactNode }) {
  const router = useRouter();
  const [pendiente, iniciar] = useTransition();
  return (
    <Link
      href={href}
      scroll={false}
      aria-busy={pendiente}
      className={className}
      onClick={(e) => {
        if (!esClicSimple(e)) return;
        e.preventDefault();
        if (!pendiente) iniciar(() => router.push(href, { scroll: false }));
      }}
    >
      {children}
      {pendiente && (
        <span className="ml-1.5 inline-block align-middle">
          <Loader mini />
        </span>
      )}
    </Link>
  );
}

// Lleva la vista al visor cuando aparece o cambia de documento.
export function AlVisor({ clave }: { clave: string }) {
  return (
    <span
      key={clave}
      ref={(el) => {
        el?.scrollIntoView({ behavior: "smooth", block: "start" });
      }}
      aria-hidden
    />
  );
}
