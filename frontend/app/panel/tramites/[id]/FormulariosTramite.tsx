"use client";

import { useActionState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import Campo, { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import type { Analista } from "@/lib/api-tramites";
import type { EstadoFormulario } from "@/lib/formularios";
import {
  archivarAction,
  asignarAction,
  darDeBajaDocumentoAction,
  dictaminarAction,
  observarAction,
  periodoAction,
  regresoAction,
  remisionAction,
} from "../actions";

const estadoInicial: EstadoFormulario = {};
type Accion = (previo: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;

function useForm(accion: Accion) {
  const [estado, formAction] = useActionState(accion, estadoInicial);
  return { estado, formAction, errores: estado.camposInvalidos };
}

// CU-09: asignación inicial o reasignación (con motivo obligatorio).
export function AsignarForm({
  id,
  analistas,
  analistaActualId,
}: {
  id: string;
  analistas: Analista[];
  analistaActualId: string | null;
}) {
  const { estado, formAction, errores } = useForm(asignarAction.bind(null, id));
  const reasignar = analistaActualId !== null;
  return (
    <Tarjeta titulo={reasignar ? "Reasignar trámite" : "Asignar a un analista"} descripcion={reasignar ? "El motivo queda permanente en el historial." : "Elige según la carga vigente."}>
      <form action={formAction} className="flex flex-col gap-3.5">
        <Campo id="analistaId" etiqueta="Analista" error={errores?.analistaId}>
          <select id="analistaId" name="analistaId" defaultValue="" className={claseControl()}>
            <option value="" disabled>
              Selecciona un analista
            </option>
            {analistas
              .filter((a) => a.id !== analistaActualId)
              .map((a) => (
                <option key={a.id} value={a.id}>
                  {a.nombre} — {a.carga} activos
                </option>
              ))}
          </select>
        </Campo>
        {reasignar && (
          <Campo id="motivo" etiqueta="Motivo de la reasignación" error={errores?.motivo}>
            <textarea id="motivo" name="motivo" rows={3} required className={claseControl()} />
          </Campo>
        )}
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar tamano="sm" textoEnviando="Guardando…">
          {reasignar ? "Reasignar" : "Asignar"}
        </BotonEnviar>
      </form>
    </Tarjeta>
  );
}

// CU-10: observación con documento afectado, motivo y plazo.
export function ObservacionForm({ id, documentos, plazoMinimo }: { id: string; documentos: string[]; plazoMinimo: string }) {
  const { estado, formAction, errores } = useForm(observarAction.bind(null, id));
  return (
    <Tarjeta titulo="Registrar observación" descripcion="El trámite pasa a Observado hasta que se resuelvan todas.">
      <form action={formAction} className="flex flex-col gap-3.5">
        <Campo id="documento" etiqueta="Documento afectado" error={errores?.documento}>
          <input id="documento" name="documento" list="documentos-exigidos" className={claseControl()} />
          <datalist id="documentos-exigidos">
            {documentos.map((d) => (
              <option key={d} value={d} />
            ))}
          </datalist>
        </Campo>
        <Campo id="motivo" etiqueta="Motivo" error={errores?.motivo}>
          <textarea id="motivo" name="motivo" rows={3} className={claseControl()} />
        </Campo>
        <Campo id="plazo" etiqueta="Plazo para subsanar" error={errores?.plazo}>
          <input id="plazo" name="plazo" type="date" min={plazoMinimo} className={claseControl()} />
        </Campo>
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar tamano="sm" textoEnviando="Registrando…">
          Registrar observación
        </BotonEnviar>
      </form>
    </Tarjeta>
  );
}

// CU-11: dictamen formal.
export function DictamenForm({ id }: { id: string }) {
  const { estado, formAction, errores } = useForm(dictaminarAction.bind(null, id));
  return (
    <Tarjeta titulo="Emitir dictamen" descripcion="Solo es posible sin observaciones pendientes.">
      <form action={formAction} className="flex flex-col gap-3.5">
        <Campo etiqueta="Resultado" error={errores?.resultado}>
          <div className="flex gap-3">
            {(["FAVORABLE", "DESFAVORABLE"] as const).map((r) => (
              <label key={r} className="flex flex-1 cursor-pointer items-center gap-2 rounded-lg border border-line px-3.5 py-2.5 text-sm font-semibold has-[:checked]:border-accent has-[:checked]:bg-accent-soft">
                <input type="radio" name="resultado" value={r} className="accent-accent" />
                {r === "FAVORABLE" ? "Favorable" : "Desfavorable"}
              </label>
            ))}
          </div>
        </Campo>
        <Campo id="texto" etiqueta="Fundamentación" error={errores?.texto}>
          <textarea id="texto" name="texto" rows={4} className={claseControl()} />
        </Campo>
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar tamano="sm" textoEnviando="Emitiendo…">
          Emitir dictamen
        </BotonEnviar>
      </form>
    </Tarjeta>
  );
}

// CU-16: baja lógica de un adjunto cargado por error. El archivo se conserva y solo el Director
// puede restituirlo, por eso el motivo es obligatorio.
export function BajaDocumentoForm({
  id,
  documentos,
  dictaminado,
}: {
  id: string;
  documentos: { id: string; etiqueta: string }[];
  dictaminado: boolean;
}) {
  const { estado, formAction, errores } = useForm(darDeBajaDocumentoAction.bind(null, id));
  return (
    <Tarjeta titulo="Dar de baja un documento" descripcion="El archivo se conserva, pero deja de cubrir el requisito y de aparecer en las búsquedas.">
      <form action={formAction} className="flex flex-col gap-3.5">
        {dictaminado && (
          <Alerta tipo="aviso">Este expediente ya tiene un dictamen emitido: justifica la baja antes de confirmarla.</Alerta>
        )}
        <Campo id="documentoId" etiqueta="Documento" error={errores?.documentoId}>
          <select id="documentoId" name="documentoId" defaultValue="" className={claseControl()}>
            <option value="" disabled>
              Selecciona el documento
            </option>
            {documentos.map((d) => (
              <option key={d.id} value={d.id}>
                {d.etiqueta}
              </option>
            ))}
          </select>
        </Campo>
        <Campo id="motivoBaja" etiqueta="Motivo de la baja" error={errores?.motivo}>
          <textarea id="motivoBaja" name="motivo" rows={3} className={claseControl()} />
        </Campo>
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Dando de baja…">
          Dar de baja
        </BotonEnviar>
      </form>
    </Tarjeta>
  );
}

// Años que se verifican en el sistema antes de dictaminar.
export function PeriodoForm({ id, desde, hasta }: { id: string; desde: number; hasta: number }) {
  const { estado, formAction, errores } = useForm(periodoAction.bind(null, id));
  return (
    <Tarjeta titulo="Años a evaluar" descripcion="Verifica en el sistema el periodo que se evalúa. Es obligatorio para dictaminar.">
      <form action={formAction} className="flex flex-col gap-3.5">
        <div className="grid grid-cols-2 gap-3">
          <Campo id="anioDesde" etiqueta="Desde" error={errores?.anioDesde}>
            <input id="anioDesde" name="anioDesde" type="number" inputMode="numeric" defaultValue={desde || ""} placeholder="2016" className={claseControl()} />
          </Campo>
          <Campo id="anioHasta" etiqueta="Hasta" error={errores?.anioHasta}>
            <input id="anioHasta" name="anioHasta" type="number" inputMode="numeric" defaultValue={hasta || ""} placeholder="2021" className={claseControl()} />
          </Campo>
        </div>
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Guardando…">
          {desde ? "Actualizar años" : "Registrar años"}
        </BotonEnviar>
      </form>
    </Tarjeta>
  );
}

// Regreso del expediente desde la facultad, con la resolución recibida. Lo registra el analista o ventanilla.
export function RegresoForm({ id, hoy }: { id: string; hoy: string }) {
  const { estado, formAction, errores } = useForm(regresoAction.bind(null, id));
  return (
    <Tarjeta titulo="Registrar regreso de la facultad" descripcion="La carpeta volvió a Escalafón con la resolución de categorización. Vuelve al mismo analista para la revisión final.">
      <form action={formAction} className="flex flex-col gap-3.5">
        <Campo id="resolucionNumero" etiqueta="Número de resolución" error={errores?.resolucionNumero}>
          <input id="resolucionNumero" name="resolucionNumero" placeholder="RES-045/2026" className={claseControl()} />
        </Campo>
        <Campo id="resolucionFecha" etiqueta="Fecha de la resolución" error={errores?.resolucionFecha}>
          <input id="resolucionFecha" name="resolucionFecha" type="date" max={hoy} className={claseControl()} />
        </Campo>
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar tamano="sm" textoEnviando="Registrando…">
          Registrar regreso
        </BotonEnviar>
      </form>
    </Tarjeta>
  );
}

// Fin de la revisión final: el expediente sale hacia el Rectorado y queda a la espera de Archivo.
export function RemisionForm({ id }: { id: string }) {
  const { estado, formAction, errores } = useForm(remisionAction.bind(null, id));
  return (
    <Tarjeta titulo="Remitir al Rectorado" descripcion="Si la resolución está conforme, remite el expediente con los ejemplares de la categorización. Después solo queda que Archivo archive la carpeta.">
      <form action={formAction} className="flex flex-col gap-3.5">
        <Campo id="nota" etiqueta="Nota de remisión" error={errores?.nota}>
          <input id="nota" name="nota" placeholder="Nota 112/2026" className={claseControl()} />
        </Campo>
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar tamano="sm" textoEnviando="Remitiendo…">
          Remitir al Rectorado
        </BotonEnviar>
      </form>
    </Tarjeta>
  );
}

// Archivo físico: lo hace el Personal de Archivo cuando el trámite ya fue remitido y las carpetas volvieron.
export function ArchivarForm({ id }: { id: string }) {
  const { estado, formAction, errores } = useForm(archivarAction.bind(null, id));
  return (
    <Tarjeta titulo="Archivar carpeta" descripcion="Registra dónde queda archivada. Después el trámite es de solo lectura.">
      <form action={formAction} className="flex flex-col gap-3.5">
        <Campo id="ubicacion" etiqueta="Ubicación de la carpeta" error={errores?.ubicacion}>
          <input id="ubicacion" name="ubicacion" placeholder="Estante 3, caja 12" className={claseControl()} />
        </Campo>
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar tamano="sm" textoEnviando="Archivando…">
          Archivar trámite
        </BotonEnviar>
      </form>
    </Tarjeta>
  );
}
