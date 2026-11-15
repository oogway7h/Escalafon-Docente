"use client";

import { useActionState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import type { EstadoFormulario } from "@/lib/formularios";

type AccionEliminar = (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;

const estadoInicial: EstadoFormulario = {};

export default function EliminarCategoriaForm({ accion, nombre }: { accion: AccionEliminar; nombre: string }) {
  const [estado, formAction] = useActionState(accion, estadoInicial);

  return (
    <form action={formAction} className="flex flex-col gap-4">
      <p className="text-sm text-ink-soft">
        Se borrará de la base de datos junto con sus documentos exigidos. Esta acción no se puede deshacer.
      </p>
      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
      <AccionesFormulario cancelarHref="/panel/categorias" textoEnviar={`Eliminar ${nombre}`} textoEnviando="Eliminando…" />
    </form>
  );
}
