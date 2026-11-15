"use client";

import { useActionState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import Campo, { claseControl } from "@/components/common/Campo";
import { cambiarContrasenaPerfil } from "./actions";
import type { EstadoFormulario } from "@/lib/formularios";

const estadoInicial: EstadoFormulario = {};

export default function ContrasenaForm() {
  const [estado, formAction] = useActionState(cambiarContrasenaPerfil, estadoInicial);

  return (
    <form action={formAction} className="flex flex-col gap-[14px]">
      <Campo id="actual" etiqueta="Contraseña actual" error={estado.camposInvalidos?.actual}>
        <input id="actual" name="actual" type="password" className={claseControl()} />
      </Campo>

      <Campo id="nueva" etiqueta="Nueva contraseña" error={estado.camposInvalidos?.nueva}>
        <input id="nueva" name="nueva" type="password" placeholder="Nueva contraseña" className={claseControl()} />
      </Campo>

      <Campo id="confirmar" etiqueta="Confirmar nueva contraseña" error={estado.camposInvalidos?.confirmar}>
        <input id="confirmar" name="confirmar" type="password" placeholder="Repite la contraseña" className={claseControl()} />
      </Campo>

      <p className="text-xs text-ink-faint">Mínimo 8 caracteres, una mayúscula y un dígito.</p>

      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
      {estado.mensaje && <Alerta tipo="exito">{estado.mensaje}</Alerta>}

      <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Actualizando…">
        Actualizar contraseña
      </BotonEnviar>
    </form>
  );
}
