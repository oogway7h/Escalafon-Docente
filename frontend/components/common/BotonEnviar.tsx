"use client";

import { useFormStatus } from "react-dom";
import BotonCarga from "./BotonCarga";
import type { TamanoBoton, VarianteBoton } from "./Boton";

// Submit que se deshabilita mientras el <form> que lo contiene se envía y muestra el Loader dentro.
// Con name/value (varios submit en un mismo form) solo carga el botón pulsado.
// sinLoader: para el inicio de sesión, que cubre su panel con su propio loader.
export default function BotonEnviar({
  children,
  textoEnviando,
  variante,
  tamano,
  clase,
  name,
  value,
  sinLoader = false,
}: {
  children: React.ReactNode;
  textoEnviando?: string;
  variante?: VarianteBoton;
  tamano?: TamanoBoton;
  clase?: string;
  name?: string;
  value?: string;
  sinLoader?: boolean;
}) {
  const { pending, data } = useFormStatus();
  const pulsado = pending && (name === undefined || data?.get(name) === value);
  return (
    <BotonCarga
      type="submit"
      name={name}
      value={value}
      disabled={pending}
      cargando={pulsado && !sinLoader}
      variante={variante}
      tamano={tamano}
      clase={clase}
    >
      {pulsado && sinLoader ? (textoEnviando ?? children) : children}
    </BotonCarga>
  );
}
