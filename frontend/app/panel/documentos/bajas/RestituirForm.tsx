"use client";

import Link from "next/link";
import { useActionState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import { claseControl } from "@/components/common/Campo";
import type { EstadoFormulario } from "@/lib/formularios";
import { restituirDocumentoAction } from "../actions";

const estadoInicial: EstadoFormulario = {};

export default function RestituirForm({ docId, cancelar }: { docId: string; cancelar: string }) {
  const [estado, accion] = useActionState(restituirDocumentoAction.bind(null, docId), estadoInicial);

  return (
    <form action={accion} className="flex flex-col gap-4">
      <label className="flex flex-col gap-1.5 text-sm font-semibold text-ink">
        Motivo de la restitución
        <span className="font-normal text-ink-faint">Queda en el registro de auditoría.</span>
        <textarea name="motivo" rows={4} required minLength={5} maxLength={255} className={`${claseControl()} resize-y font-normal`} />
      </label>
      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
      <div className="flex flex-wrap gap-2">
        <BotonEnviar textoEnviando="Restituyendo…">Restituir y reindexar</BotonEnviar>
        <Link href={cancelar} className="inline-flex h-11 items-center rounded-xl border border-line px-5 text-sm font-semibold text-ink hover:border-accent hover:text-accent-ink">
          Cancelar
        </Link>
      </div>
    </form>
  );
}
