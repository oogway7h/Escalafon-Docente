"use client";

import { useActionState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import Campo, { claseControl } from "@/components/common/Campo";
import { actualizarContacto } from "./actions";
import type { EstadoFormulario } from "@/lib/formularios";

const estadoInicial: EstadoFormulario = {};

export default function ContactoForm({ telefono }: { telefono: string }) {
  const [estado, formAction] = useActionState(actualizarContacto, estadoInicial);

  return (
    <form action={formAction} className="flex flex-col gap-[14px]">
      <Campo id="telefono" etiqueta="Teléfono" error={estado.camposInvalidos?.telefono}>
        <input id="telefono" name="telefono" defaultValue={telefono} placeholder="+591 700 00000" className={claseControl()} />
      </Campo>

      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
      {estado.mensaje && <Alerta tipo="exito">{estado.mensaje}</Alerta>}

      <BotonEnviar tamano="sm" textoEnviando="Guardando…">
        Guardar cambios
      </BotonEnviar>
    </form>
  );
}
