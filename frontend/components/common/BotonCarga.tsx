import type { ComponentProps } from "react";
import { claseBoton, type TamanoBoton, type VarianteBoton } from "./Boton";
import Loader from "./Loader";

// Botón único del sistema para acciones con espera: mientras `cargando`, se deshabilita y muestra
// el Loader dentro (en lugar del texto, sin cambiar el tamaño del botón). Lo usan BotonEnviar (formularios) y BotonFiltrar (listados).
export default function BotonCarga({
  cargando = false,
  variante,
  tamano,
  clase,
  children,
  ...props
}: Omit<ComponentProps<"button">, "className"> & {
  cargando?: boolean;
  variante?: VarianteBoton;
  tamano?: TamanoBoton;
  clase?: string;
}) {
  return (
    <button {...props} aria-busy={cargando} className={clase ?? claseBoton(variante, tamano)}>
      {cargando ? (
        // El texto sigue ocupando su espacio (invisible) y el Loader se superpone: el botón no cambia de tamaño.
        <span className="relative inline-flex items-center justify-center">
          <span className="invisible">{children}</span>
          <span className="absolute inset-0 flex items-center justify-center">
            <Loader mini />
          </span>
        </span>
      ) : (
        children
      )}
    </button>
  );
}
