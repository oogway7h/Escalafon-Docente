import { listarRoles, exigirPermiso, obtenerUsuario } from "@/lib/api";
import CambiarRolForm from "../../CambiarRolForm";
import { cambiarRolAction } from "../../actions";

export default async function CambiarRolPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  await exigirPermiso("USUARIO_GESTIONAR");

  const [usuario, roles] = await Promise.all([obtenerUsuario(id), listarRoles()]);
  const rolActual = roles.find((rol) => rol.id === usuario.rolId);
  const accion = cambiarRolAction.bind(null, id);

  return <CambiarRolForm accion={accion} usuario={usuario} roles={roles} rolActual={rolActual} />;
}
