import { listarRoles, exigirPermiso, obtenerUsuario } from "@/lib/api";
import UsuarioForm from "../../UsuarioForm";
import { editarUsuarioAction } from "../../actions";

export default async function EditarUsuarioPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  await exigirPermiso("USUARIO_GESTIONAR");

  const [usuario, roles] = await Promise.all([obtenerUsuario(id), listarRoles()]);
  const accion = editarUsuarioAction.bind(null, id);

  return (
    <UsuarioForm
      accion={accion}
      roles={roles}
      valoresIniciales={{
        nombres: usuario.nombres,
        apellidos: usuario.apellidos,
        correo: usuario.correo,
        telefono: usuario.telefono,
        rolId: usuario.rolId,
      }}
      titulo={`Editar usuario — ${usuario.nombres} ${usuario.apellidos}`}
      textoBoton="Guardar usuario"
    />
  );
}
