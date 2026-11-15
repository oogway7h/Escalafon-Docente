"use client";

import { useActionState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import Campo, { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import type { RolConConteo } from "@/lib/api";
import type { EstadoFormulario } from "@/lib/formularios";

type AccionUsuario = (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;

const estadoInicial: EstadoFormulario = {};

export default function UsuarioForm({
  accion,
  roles,
  valoresIniciales,
  titulo,
  textoBoton,
}: {
  accion: AccionUsuario;
  roles: RolConConteo[];
  valoresIniciales?: { nombres: string; apellidos: string; correo: string; telefono: string; rolId?: number };
  titulo: string;
  textoBoton: string;
}) {
  const [estado, formAction] = useActionState(accion, estadoInicial);
  const errores = estado.camposInvalidos;

  return (
    <Tarjeta nivel="h1" titulo={titulo} className="max-w-xl">
      <form action={formAction} className="flex flex-col gap-4">
        <div className="flex gap-4">
          <Campo id="nombres" etiqueta="Nombres" error={errores?.nombres} className="flex-1">
            <input id="nombres" name="nombres" defaultValue={valoresIniciales?.nombres} className={claseControl()} />
          </Campo>
          <Campo id="apellidos" etiqueta="Apellidos" error={errores?.apellidos} className="flex-1">
            <input id="apellidos" name="apellidos" defaultValue={valoresIniciales?.apellidos} className={claseControl()} />
          </Campo>
        </div>

        <Campo
          id="correo"
          etiqueta="Correo institucional"
          ayuda="Debe ser único (sin distinguir mayúsculas)."
          error={errores?.correo}
        >
          <input id="correo" name="correo" type="email" defaultValue={valoresIniciales?.correo} className={claseControl()} />
        </Campo>

        <div className="flex gap-4">
          <Campo id="telefono" etiqueta="Teléfono" opcional className="flex-1">
            <input
              id="telefono"
              name="telefono"
              defaultValue={valoresIniciales?.telefono}
              placeholder="+591 700 00000"
              className={claseControl()}
            />
          </Campo>

          {valoresIniciales?.rolId === undefined ? (
            <Campo id="rolId" etiqueta="Rol" error={errores?.rolId} className="flex-1">
              <select id="rolId" name="rolId" defaultValue="" className={claseControl()}>
                <option value="" disabled>
                  Selecciona un rol
                </option>
                {roles.map((rol) => (
                  <option key={rol.id} value={rol.id}>
                    {rol.nombre}
                  </option>
                ))}
              </select>
            </Campo>
          ) : (
            <Campo etiqueta="Rol" ayuda='Cámbialo desde "Cambiar rol" en el listado.' className="flex-1">
              <div className={claseControl({ bloqueado: true })}>
                {roles.find((r) => r.id === valoresIniciales.rolId)?.nombre ?? "—"}
              </div>
            </Campo>
          )}
        </div>

        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

        <AccionesFormulario cancelarHref="/panel/usuarios" textoEnviar={textoBoton} />
      </form>
    </Tarjeta>
  );
}
