"use client";

import { useActionState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import Campo, { claseControl } from "@/components/common/Campo";
import { solicitarRecuperacion } from "./actions";
import type { EstadoFormulario } from "@/lib/formularios";

const estadoInicial: EstadoFormulario = {};

export default function RecuperarForm() {
  const [estado, formAction] = useActionState(solicitarRecuperacion, estadoInicial);

  if (estado.mensaje) return <Alerta tipo="aviso">{estado.mensaje}</Alerta>;

  return (
    <form action={formAction} className="flex flex-col gap-[18px]">
      <Campo id="correo" etiqueta="Correo institucional" error={estado.camposInvalidos?.correo}>
        <input
          id="correo"
          name="correo"
          type="email"
          required
          placeholder="nombre.apellido@uagrm.edu.bo"
          className={claseControl({ grande: true })}
        />
      </Campo>

      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

      <BotonEnviar tamano="lg" textoEnviando="Enviando…">
        Enviar enlace
      </BotonEnviar>
    </form>
  );
}
