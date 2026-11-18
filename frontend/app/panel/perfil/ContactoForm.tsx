"use client";

import { useActionState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import Campo, { claseControl } from "@/components/common/Campo";
import { actualizarContacto } from "./actions";
import type { EstadoFormulario } from "@/lib/formularios";

const estadoInicial: EstadoFormulario = {};

export default function ContactoForm({ nombre, correo, telefono }: { nombre: string; correo: string; telefono: string }) {
  const [estado, formAction] = useActionState(actualizarContacto, estadoInicial);

  return (
    <form action={formAction} className="flex flex-col gap-[14px]">
      <div className="flex flex-wrap gap-[14px]">
        <Campo id="nombre" etiqueta="Nombre" className="min-w-[220px] flex-1">
          <input id="nombre" value={nombre} disabled readOnly className={claseControl({ bloqueado: true })} />
        </Campo>
        <Campo id="correo" etiqueta="Correo" className="min-w-[220px] flex-1">
          <input id="correo" value={correo} disabled readOnly className={claseControl({ bloqueado: true })} />
        </Campo>
      </div>

      <Campo id="telefono" etiqueta="Teléfono" className="max-w-sm" error={estado.camposInvalidos?.telefono}>
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
