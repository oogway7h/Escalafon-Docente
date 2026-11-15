"use client";

import { useActionState } from "react";
import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import BotonEnviar from "@/components/common/BotonEnviar";
import Campo, { claseControl } from "@/components/common/Campo";
import type { EstadoFormulario } from "@/lib/formularios";

const estadoInicial: EstadoFormulario = {};

// Formulario "elige tu contraseña" a partir de un token de correo (activar cuenta, restablecer contraseña).
export default function ContrasenaConTokenForm({
  accion,
  token,
  etiquetaNueva,
  etiquetaConfirmar,
  textoBoton,
  textoEnviando,
  sinToken,
}: {
  accion: (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;
  token: string;
  etiquetaNueva: string;
  etiquetaConfirmar: string;
  textoBoton: string;
  textoEnviando: string;
  sinToken: React.ReactNode;
}) {
  const [estado, formAction] = useActionState(accion, estadoInicial);

  if (estado.mensaje) {
    return (
      <div className="flex flex-col gap-4">
        <Alerta tipo="aviso">{estado.mensaje}</Alerta>
        <a href="/login" className={claseBoton("primario", "lg")}>
          Iniciar sesión
        </a>
      </div>
    );
  }

  if (!token) return <Alerta tipo="aviso">{sinToken}</Alerta>;

  return (
    <form action={formAction} className="flex flex-col gap-[14px]">
      <input type="hidden" name="token" value={token} />

      <Campo id="nueva" etiqueta={etiquetaNueva} error={estado.camposInvalidos?.nueva}>
        <input
          id="nueva"
          name="nueva"
          type="password"
          required
          placeholder="Mínimo 8 caracteres"
          className={claseControl({ grande: true })}
        />
      </Campo>

      <Campo id="confirmar" etiqueta={etiquetaConfirmar} error={estado.camposInvalidos?.confirmar}>
        <input
          id="confirmar"
          name="confirmar"
          type="password"
          required
          placeholder="Repite la contraseña"
          className={claseControl({ grande: true })}
        />
      </Campo>

      <p className="text-xs text-ink-faint">Debe incluir al menos una mayúscula y un dígito.</p>

      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

      <BotonEnviar tamano="lg" textoEnviando={textoEnviando}>
        {textoBoton}
      </BotonEnviar>
    </form>
  );
}
