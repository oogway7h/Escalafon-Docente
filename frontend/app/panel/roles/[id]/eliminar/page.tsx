import Link from "next/link";
import { notFound } from "next/navigation";
import { claseBoton } from "@/components/common/Boton";
import Tarjeta from "@/components/common/Tarjeta";
import { listarRoles, exigirPermiso } from "@/lib/api";
import { eliminarRolAction } from "../../actions";
import EliminarRolForm from "../../EliminarRolForm";

export default async function EliminarRolPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  await exigirPermiso("ROL_GESTIONAR");

  const roles = await listarRoles();
  const rol = roles.find((r) => r.id === Number(id));
  if (!rol) notFound();

  const bloqueado = rol.esSistema || rol.usuariosCount > 0;

  return (
    <Tarjeta
      nivel="h1"
      titulo={bloqueado ? `No se puede eliminar "${rol.nombre}"` : `Eliminar "${rol.nombre}"`}
      descripcion={
        bloqueado &&
        (rol.esSistema
          ? "Es un rol de sistema: no se puede eliminar."
          : `Tiene ${rol.usuariosCount} ${rol.usuariosCount === 1 ? "usuario asignado" : "usuarios asignados"}. Reasigna o desactiva a esas personas primero.`)
      }
      className="max-w-md"
    >
      {bloqueado ? (
        <Link href="/panel/roles" className={claseBoton("secundario")}>
          Entendido
        </Link>
      ) : (
        <EliminarRolForm accion={eliminarRolAction.bind(null, rol.id)} nombre={rol.nombre} />
      )}
    </Tarjeta>
  );
}
