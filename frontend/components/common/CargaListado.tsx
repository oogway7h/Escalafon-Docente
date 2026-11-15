"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { createContext, useContext, useEffect, useRef, useState, useTransition, type ComponentProps } from "react";
import type { TamanoBoton, VarianteBoton } from "./Boton";
import BotonCarga from "./BotonCarga";
import { claseControl } from "./Campo";
import Loader from "./Loader";
import TablaEsqueleto from "./TablaEsqueleto";

// Listados que cambian por URL (búsqueda, página). Mientras llega el resultado nuevo se conserva la
// tabla anterior, difuminada, con el Loader encima; así no se pierde el contexto ni parpadea el layout.
//
//   <ListadoConCarga>
//     <BuscadorCarga action="/panel/x" … />
//     <ZonaCarga><Tabla … /></ZonaCarga>
//     <Paginacion … />            (sus enlaces usan LinkCarga)
//   </ListadoConCarga>
//
// Sin JavaScript todo sigue funcionando como enlaces y formularios normales.

type Navegacion = {
  pendiente: boolean;
  navegar: (href: string) => void;
  ejecutar: (accion: () => void | Promise<void>) => void;
};

const Contexto = createContext<Navegacion | null>(null);

// Para componentes propios (p. ej. BuscadorSugerencias) que navegan con la misma transición del listado.
export function useCargaListado(): Navegacion | null {
  return useContext(Contexto);
}

export function ListadoConCarga({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [pendiente, iniciar] = useTransition();
  const navegar = (href: string) => iniciar(() => router.push(href));
  // ejecutar mantiene la carga hasta que termina la acción de servidor y llega la vista actualizada.
  const ejecutar = (accion: () => void | Promise<void>) =>
    iniciar(async () => {
      await accion();
    });
  return <Contexto.Provider value={{ pendiente, navegar, ejecutar }}>{children}</Contexto.Provider>;
}

export function ZonaCarga({ children, className = "" }: { children: React.ReactNode; className?: string }) {
  const pendiente = useContext(Contexto)?.pendiente ?? false;
  const caja = useRef<HTMLDivElement>(null);
  // Forma de la tabla que hay ahora en pantalla. Se mide del DOM en vez de pedirla por props
  // para no tener que tocar las catorce páginas que usan ZonaCarga, y porque así la silueta
  // tiene exactamente las columnas y las filas que el usuario acaba de ver.
  const [forma, setForma] = useState<{ columnas: number; filas: number } | null>(null);

  useEffect(() => {
    if (pendiente) return; // durante la espera los hijos no están montados: no hay nada que medir
    const tabla = caja.current?.querySelector("table");
    const columnas = tabla?.querySelectorAll("thead th").length ?? 0;
    if (!tabla || columnas === 0) {
      setForma((previa) => (previa === null ? previa : null));
      return;
    }
    const filas = Math.min(Math.max(tabla.querySelectorAll("tbody tr").length, 3), 12);
    // Comparar antes de asignar: children cambia de identidad en cada render y un objeto nuevo
    // cada vez volvería a disparar este efecto en bucle.
    setForma((previa) =>
      previa && previa.columnas === columnas && previa.filas === filas ? previa : { columnas, filas },
    );
  }, [pendiente, children]);

  // Listado con tabla: se sustituye por su silueta, que conserva la altura de la página.
  if (pendiente && forma) {
    return (
      <div className={className} role="status" aria-busy="true" aria-label="Cargando resultados">
        <TablaEsqueleto columnas={forma.columnas} filas={forma.filas} />
      </div>
    );
  }

  // Sin tabla (tablero, custodia, buscador de documentos) no hay silueta que valga: ahí se
  // mantiene el contenido anterior difuminado con el Loader encima.
  return (
    <div className={`relative ${className}`} aria-busy={pendiente} ref={caja}>
      <div
        className={`min-h-56 transition-[filter,opacity] duration-200 ${pendiente ? "pointer-events-none select-none opacity-60 blur-[3px]" : ""}`}
        inert={pendiente}
      >
        {children}
      </div>
      {pendiente && (
        <div className="animate-materializar absolute inset-0 z-10">
          {/* sticky: el loader sigue en el centro de la ventana aunque la tabla sea más alta que ella */}
          <div className="sticky top-[40vh]">
            <Loader compacto />
          </div>
        </div>
      )}
    </div>
  );
}

// Enlace que, dentro de un ListadoConCarga, navega con transición para que se vea la carga.
export function LinkCarga({ href, onClick, ...props }: ComponentProps<typeof Link> & { href: string }) {
  const nav = useContext(Contexto);
  return (
    <Link
      href={href}
      onClick={(e) => {
        onClick?.(e);
        if (!nav || e.defaultPrevented || e.metaKey || e.ctrlKey || e.shiftKey || e.button !== 0) return;
        e.preventDefault();
        nav.navegar(href);
      }}
      {...props}
    />
  );
}

// Formulario con una acción de servidor (dar de baja, reactivar…): mientras se ejecuta, el listado
// se difumina con el Loader. Sin ListadoConCarga funciona como un formulario normal.
// Submit de búsqueda o filtros: mientras el listado carga, muestra el Loader dentro del botón.
export function BotonFiltrar({
  children,
  variante = "secundario",
  tamano = "sm",
}: {
  children: React.ReactNode;
  variante?: VarianteBoton;
  tamano?: TamanoBoton;
}) {
  const pendiente = useContext(Contexto)?.pendiente ?? false;
  return (
    <BotonCarga type="submit" cargando={pendiente} disabled={pendiente} variante={variante} tamano={tamano}>
      {children}
    </BotonCarga>
  );
}

export function FormularioCarga({
  action,
  className,
  children,
}: {
  action: (formData: FormData) => void | Promise<void>;
  className?: string;
  children: React.ReactNode;
}) {
  const nav = useContext(Contexto);
  return (
    <form className={className} action={nav ? (formData: FormData) => nav.ejecutar(() => action(formData)) : action}>
      {children}
    </form>
  );
}

// Formulario GET de varios filtros: navega con la carga visible y omite los campos vacíos de la URL.
// Al aplicar filtros se vuelve a la primera página (no se envía ?pagina).
export function FiltrosCarga({
  action,
  className,
  children,
}: {
  action: string;
  className?: string;
  children: React.ReactNode;
}) {
  const nav = useContext(Contexto);
  return (
    <form
      action={action}
      method="get"
      className={className}
      onSubmit={(e) => {
        if (!nav) return;
        e.preventDefault();
        const qs = new URLSearchParams();
        for (const [clave, valor] of new FormData(e.currentTarget)) {
          if (typeof valor === "string" && valor.trim()) qs.set(clave, valor.trim());
        }
        const consulta = qs.toString();
        nav.navegar(consulta ? `${action}?${consulta}` : action);
      }}
    >
      {children}
    </form>
  );
}

export function BuscadorCarga({
  action,
  defaultValue,
  placeholder,
}: {
  action: string;
  defaultValue?: string;
  placeholder: string;
}) {
  const nav = useContext(Contexto);
  return (
    <form
      action={action}
      method="get"
      className="flex gap-2.5"
      onSubmit={(e) => {
        if (!nav) return;
        e.preventDefault();
        const texto = String(new FormData(e.currentTarget).get("buscar") ?? "").trim();
        nav.navegar(texto ? `${action}?buscar=${encodeURIComponent(texto)}` : action);
      }}
    >
      <input
        type="search"
        name="buscar"
        defaultValue={defaultValue}
        placeholder={placeholder}
        className={`${claseControl()} max-w-[320px] flex-1`}
      />
      <BotonFiltrar>Buscar</BotonFiltrar>
    </form>
  );
}
