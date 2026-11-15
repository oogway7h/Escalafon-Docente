import { listarPermisos, obtenerRol, exigirPermiso } from "@/lib/api";
import RolForm from "../../RolForm";
import { guardarRolAction } from "../../actions";

export default async function EditarRolPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  await exigirPermiso("ROL_GESTIONAR");

  const [rol, permisos] = await Promise.all([obtenerRol(id), listarPermisos()]);
  const accion = guardarRolAction.bind(null, rol.id);

  return (
    <RolForm
      accion={accion}
      permisos={permisos}
      esSistema={rol.esSistema}
      valoresIniciales={{
        nombre: rol.nombre,
        descripcion: rol.descripcion,
        permisoIds: rol.permisos.map((p) => p.id),
      }}
      titulo={`Editar rol — ${rol.nombre}`}
      textoBoton="Guardar rol"
    />
  );
}
