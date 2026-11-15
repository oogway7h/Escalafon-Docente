import { listarRoles, exigirPermiso } from "@/lib/api";
import UsuarioForm from "../UsuarioForm";
import { crearUsuarioAction } from "../actions";

export default async function NuevoUsuarioPage() {
  await exigirPermiso("USUARIO_GESTIONAR");

  const roles = await listarRoles();

  return <UsuarioForm accion={crearUsuarioAction} roles={roles} titulo="Nuevo usuario" textoBoton="Crear usuario" />;
}
