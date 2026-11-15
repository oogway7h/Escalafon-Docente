"use client";

import { useActionState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import Campo, { claseControl } from "@/components/common/Campo";
import Insignia from "@/components/common/Insignia";
import Tarjeta from "@/components/common/Tarjeta";
import type { RolConConteo, UsuarioListado } from "@/lib/api";
import type { EstadoFormulario } from "@/lib/formularios";

type AccionRol = (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;

const estadoInicial: EstadoFormulario = {};

export default function CambiarRolForm({
  accion,
  usuario,
  roles,
  rolActual,
}: {
  accion: AccionRol;
  usuario: UsuarioListado;
  roles: RolConConteo[];
  rolActual?: RolConConteo;
}) {
  const [estado, formAction] = useActionState(accion, estadoInicial);

  return (
    <Tarjeta
      nivel="h1"
      titulo="Cambiar rol"
      descripcion={`${usuario.nombres} ${usuario.apellidos} — al cambiar el rol se cierra su sesión activa.`}
      className="max-w-md"
    >
      <form action={formAction} className="flex flex-col gap-4">
        <Campo etiqueta="Rol actual">
          <Insignia tono="neutro" className="w-fit">
            {rolActual?.nombre ?? "—"}
          </Insignia>
        </Campo>

        <Campo id="rolId" etiqueta="Nuevo rol" error={estado.camposInvalidos?.rolId}>
          <select id="rolId" name="rolId" defaultValue={usuario.rolId} className={claseControl()}>
            {roles.map((rol) => (
              <option key={rol.id} value={rol.id}>
                {rol.nombre}
              </option>
            ))}
          </select>
        </Campo>

        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

        <AccionesFormulario cancelarHref="/panel/usuarios" textoEnviar="Confirmar cambio" />
      </form>
    </Tarjeta>
  );
}
