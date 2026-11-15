import Link from "next/link";
import { claseBoton } from "./Boton";
import BotonEnviar from "./BotonEnviar";

// Cancelar + submit al pie de los formularios del panel.
export default function AccionesFormulario({
  cancelarHref,
  textoEnviar,
  textoEnviando = "Guardando…",
}: {
  cancelarHref: string;
  textoEnviar: string;
  textoEnviando?: string;
}) {
  return (
    <div className="mt-1 flex justify-end gap-3">
      <Link href={cancelarHref} className={claseBoton("secundario")}>
        Cancelar
      </Link>
      <BotonEnviar textoEnviando={textoEnviando}>{textoEnviar}</BotonEnviar>
    </div>
  );
}
