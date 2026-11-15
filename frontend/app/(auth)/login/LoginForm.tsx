"use client";

import { useActionState } from "react";
import { useFormStatus } from "react-dom";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import Campo, { claseControl } from "@/components/common/Campo";
import Loader from "@/components/common/Loader";
import type { EstadoFormulario } from "@/lib/formularios";
import { iniciarSesion } from "./actions";

const estadoInicial: EstadoFormulario = {};

// Cubre todo el panel del formulario (título, campos y pie), no la pantalla completa:
// se posiciona respecto al contenedor con position:relative de AuthDividido (el lado derecho).
function CargandoSesion() {
  const { pending } = useFormStatus();
  return pending ? (
    <div className="absolute inset-0 z-50 flex items-center justify-center bg-bg/60 backdrop-blur-sm">
      <Loader compacto />
    </div>
  ) : null;
}

export default function LoginForm() {
  const [estado, formAction] = useActionState(iniciarSesion, estadoInicial);

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

      <Campo id="contrasena" etiqueta="Contraseña" error={estado.camposInvalidos?.contrasena}>
        <input id="contrasena" name="contrasena" type="password" required className={claseControl({ grande: true })} />
      </Campo>

      <div className="mt-0.5 flex items-center justify-between">
        <label className="flex items-center gap-2 text-sm text-ink-soft">
          <input
            type="checkbox"
            name="recordarme"
            defaultChecked
            className="h-4 w-4 rounded border-[1.5px] border-line accent-accent"
          />
          Recordarme
        </label>
        <a href="/recuperar" className="text-sm font-semibold text-accent-ink hover:text-accent-dark hover:underline">
          ¿Olvidaste tu contraseña?
        </a>
      </div>

      <CargandoSesion />

      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

      <BotonEnviar tamano="lg" textoEnviando="Iniciando sesión…" sinLoader>
        Iniciar sesión
      </BotonEnviar>
    </form>
  );
}
