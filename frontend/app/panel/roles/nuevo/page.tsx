import { listarPermisos, exigirPermiso } from "@/lib/api";
import RolForm from "../RolForm";
import { guardarRolAction } from "../actions";

export default async function NuevoRolPage() {
  await exigirPermiso("ROL_GESTIONAR");

  const permisos = await listarPermisos();
  const accion = guardarRolAction.bind(null, null);

  return <RolForm accion={accion} permisos={permisos} titulo="Nuevo rol" textoBoton="Crear rol" />;
}
