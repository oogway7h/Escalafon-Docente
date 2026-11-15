"use client";

import { useActionState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import type { EstadoFormulario } from "@/lib/formularios";

type AccionEliminar = (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;

const estadoInicial: EstadoFormulario = {};

export default function EliminarRolForm({ accion, nombre }: { accion: AccionEliminar; nombre: string }) {
  const [estado, formAction] = useActionState(accion, estadoInicial);

  return (
    <form action={formAction} className="flex flex-col gap-4">
      <p className="text-sm text-ink-soft">Esta acción no se puede deshacer.</p>
      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
      <AccionesFormulario cancelarHref="/panel/roles" textoEnviar={`Eliminar ${nombre}`} textoEnviando="Eliminando…" />
    </form>
  );
}
