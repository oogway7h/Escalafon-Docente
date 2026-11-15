"use client";

import { useActionState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import Campo, { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import type { Categoria, Docente, Facultad } from "@/lib/api-tramites";
import type { EstadoFormulario } from "@/lib/formularios";

const estadoInicial: EstadoFormulario = {};

export default function DocenteForm({
  accion,
  categorias,
  facultades,
  docente,
  titulo,
}: {
  accion: (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;
  categorias: Categoria[];
  facultades: Facultad[];
  docente?: Docente;
  titulo: string;
}) {
  const [estado, formAction] = useActionState(accion, estadoInicial);
  const errores = estado.camposInvalidos;

  return (
    <Tarjeta nivel="h1" titulo={titulo} className="max-w-xl">
      <form action={formAction} className="flex flex-col gap-4">
        <div className="flex gap-4">
          <Campo id="codigo" etiqueta="Código de docente" ayuda="El que otorga la universidad. Es único." error={errores?.codigo} className="flex-1">
            <input id="codigo" name="codigo" inputMode="numeric" defaultValue={docente?.codigo} className={claseControl()} />
          </Campo>
          <Campo
            id="cedula"
            etiqueta="Cédula de identidad"
            opcional
            ayuda={docente?.cedula ? "Una vez asignada, no se puede cambiar." : undefined}
            error={errores?.cedula}
            className="flex-1"
          >
            {docente?.cedula ? (
              <div className={claseControl({ bloqueado: true })}>{docente.cedula}</div>
            ) : (
              <input id="cedula" name="cedula" inputMode="numeric" className={claseControl()} />
            )}
          </Campo>
        </div>

        <div className="flex gap-4">
          <Campo id="nombres" etiqueta="Nombres" error={errores?.nombres} className="flex-1">
            <input id="nombres" name="nombres" defaultValue={docente?.nombres} className={claseControl()} />
          </Campo>
          <Campo id="apellidos" etiqueta="Apellidos" error={errores?.apellidos} className="flex-1">
            <input id="apellidos" name="apellidos" defaultValue={docente?.apellidos} className={claseControl()} />
          </Campo>
        </div>

        <div className="flex gap-4">
          <Campo id="correo" etiqueta="Correo" opcional error={errores?.correo} className="flex-1">
            <input id="correo" name="correo" type="email" defaultValue={docente?.correo} className={claseControl()} />
          </Campo>
          <Campo id="telefono" etiqueta="Teléfono" opcional className="flex-1">
            <input id="telefono" name="telefono" defaultValue={docente?.telefono} className={claseControl()} />
          </Campo>
        </div>

        <Campo id="facultadId" etiqueta="Facultad / carrera" opcional>
          <select id="facultadId" name="facultadId" defaultValue={docente?.facultadId ?? ""} className={claseControl()}>
            <option value="">Sin facultad</option>
            {facultades.map((f) => (
              <option key={f.id} value={f.id}>
                {f.nombre}
              </option>
            ))}
          </select>
        </Campo>

        <Campo id="categoriaId" etiqueta="Categoría actual" opcional>
          <select id="categoriaId" name="categoriaId" defaultValue={docente?.categoriaId ?? ""} className={claseControl()}>
            <option value="">Sin categoría</option>
            {categorias.filter((c) => c.activa || c.id === docente?.categoriaId).map((c) => (
              <option key={c.id} value={c.id}>
                {c.nombre}
              </option>
            ))}
          </select>
        </Campo>

        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

        <AccionesFormulario cancelarHref="/panel/docentes" textoEnviar="Guardar docente" />
      </form>
    </Tarjeta>
  );
}
