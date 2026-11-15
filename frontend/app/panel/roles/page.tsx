import Link from "next/link";
import { claseBoton } from "@/components/common/Boton";
import { LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Tabla from "@/components/common/Tabla";
import Tarjeta from "@/components/common/Tarjeta";
import { IconoEditar, IconoEliminar } from "@/components/layout/IconosPanel";
import { listarPermisos, listarRoles, exigirPermiso } from "@/lib/api";

export default async function RolesPage() {
  const sesion = await exigirPermiso("ROL_GESTIONAR");

  const [roles, permisos] = await Promise.all([listarRoles(), listarPermisos()]);

  const porModulo = new Map<string, number>();
  for (const permiso of permisos) {
    porModulo.set(permiso.modulo, (porModulo.get(permiso.modulo) ?? 0) + 1);
  }

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Roles y permisos"
          descripcion={<>{roles.length} roles · {permisos.length} permisos en el catálogo.</>}
          accion={
            <LinkCarga href="/panel/roles/nuevo" className={claseBoton("primario", "sm")}>
              + Nuevo rol
            </LinkCarga>
          }
        />

        <div className="flex flex-col gap-5 lg:flex-row lg:items-start">
          <ZonaCarga className="flex-[1.4]">
          <Tabla
            filas={roles}
            clave={(rol) => rol.id}
            pie="Los roles de sistema (Director, Analista, Personal de Archivo) no se pueden eliminar ni renombrar. Un rol con usuarios asignados tampoco se puede eliminar."
            columnas={[
              {
                titulo: "Rol",
                celda: (rol) => (
                  <>
                    <div className="font-semibold">{rol.nombre}</div>
                    <div className="text-xs text-ink-faint">
                      {rol.esSistema ? "Rol de sistema · protegido" : "Rol personalizado"}
                    </div>
                  </>
                ),
              },
              {
                titulo: "Usuarios",
                celda: (rol) => (
                  <>
                    {rol.usuariosCount} {rol.usuariosCount === 1 ? "usuario" : "usuarios"}
                    {rol.nombre === sesion.rol && (
                      <span className="ml-1.5 text-xs text-ink-faint">(incluye tu cuenta)</span>
                    )}
                  </>
                ),
              },
              {
                className: "w-24 text-right",
                celda: (rol) => (
                  <div className="flex justify-end gap-1.5">
                    <Link
                      href={`/panel/roles/${rol.id}/editar`}
                      className="flex h-8 w-8 items-center justify-center rounded-lg border border-line text-ink-soft hover:border-accent hover:text-accent-ink"
                      aria-label={`Editar ${rol.nombre}`}
                    >
                      <IconoEditar className="h-3.5 w-3.5" />
                    </Link>
                    {!rol.esSistema && rol.usuariosCount === 0 && (
                      <Link
                        href={`/panel/roles/${rol.id}/eliminar`}
                        className="flex h-8 w-8 items-center justify-center rounded-lg border border-accent-soft text-accent-ink hover:bg-accent-soft"
                        aria-label={`Eliminar ${rol.nombre}`}
                      >
                        <IconoEliminar className="h-3.5 w-3.5" />
                      </Link>
                    )}
                  </div>
                ),
              },
            ]}
          />
          </ZonaCarga>

          <Tarjeta titulo="Catálogo de permisos" descripcion="Se usa al armar cualquier rol." className="max-w-[360px] flex-1">
            <div className="flex flex-col gap-2.5">
              {Array.from(porModulo.entries()).map(([modulo, cantidad]) => (
                <div key={modulo} className="flex justify-between text-sm">
                  <span className="text-ink">{modulo}</span>
                  <span className="text-ink-faint">
                    {cantidad} {cantidad === 1 ? "permiso" : "permisos"}
                  </span>
                </div>
              ))}
            </div>
          </Tarjeta>
        </div>
      </div>
    </ListadoConCarga>
  );
}
