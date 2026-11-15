import BotonEnviar from "@/components/common/BotonEnviar";
import Link from "next/link";
import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import { BotonFiltrar, FiltrosCarga, LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import { claseControl } from "@/components/common/Campo";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Paginacion from "@/components/common/Paginacion";
import Insignia from "@/components/common/Insignia";
import Tabla from "@/components/common/Tabla";
import { IconoMas } from "@/components/layout/IconosPanel";
import { listarRoles, listarUsuarios, exigirPermiso } from "@/lib/api";
import { mensajeError } from "@/lib/formularios";
import { cambiarEstadoAction, reenviarActivacionAction } from "./actions";

const claseItemMenu = "block w-full rounded-lg px-3 py-2 text-left text-sm font-semibold";

export default async function UsuariosPage({
  searchParams,
}: {
  searchParams: Promise<{ buscar?: string; rolId?: string; activo?: string; pagina?: string; error?: string; exito?: string }>;
}) {
  const sesion = await exigirPermiso("USUARIO_GESTIONAR");

  const sp = await searchParams;
  const [lista, roles] = await Promise.all([
    listarUsuarios({ buscar: sp.buscar, rolId: sp.rolId, activo: sp.activo, pagina: sp.pagina }),
    listarRoles(),
  ]);
  const rolPorId = new Map(roles.map((rol) => [rol.id, rol.nombre]));
  const totalPaginas = Math.max(1, Math.ceil(lista.total / lista.tamano));

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Usuarios"
          descripcion={<>{lista.total} personas con acceso al sistema.</>}
          accion={
            <LinkCarga href="/panel/usuarios/nuevo" className={claseBoton("primario", "sm")}>
              + Nuevo usuario
            </LinkCarga>
          }
        />

        {sp.error && (
          <Alerta tipo="error">{mensajeError(sp.error, "No se pudo completar la acción.")}</Alerta>
        )}
        {sp.exito === "REENVIADO" && (
          <Alerta tipo="exito">Correo de activación reenviado.</Alerta>
        )}

        <FiltrosCarga action="/panel/usuarios" className="flex gap-2.5">
          <input
            type="text"
            name="buscar"
            defaultValue={sp.buscar}
            placeholder="Buscar por nombre o correo"
            className={`${claseControl()} max-w-[280px] flex-1`}
          />
          <select
            name="rolId"
            defaultValue={sp.rolId ?? ""}
            className={claseControl()}
          >
            <option value="">Todos los roles</option>
            {roles.map((rol) => (
              <option key={rol.id} value={rol.id}>
                {rol.nombre}
              </option>
            ))}
          </select>
          <select
            name="activo"
            defaultValue={sp.activo ?? ""}
            className={claseControl()}
          >
            <option value="">Todos los estados</option>
            <option value="true">Activos</option>
            <option value="false">Inactivos</option>
          </select>
          <BotonFiltrar>Filtrar</BotonFiltrar>
        </FiltrosCarga>

        <ZonaCarga>
        <Tabla
          filas={lista.datos}
          clave={(usuario) => usuario.id}
          vacio="No hay usuarios que coincidan con el filtro."
          columnas={[
            {
              titulo: "Usuario",
              celda: (usuario) => (
                <div className="flex items-center gap-2.5">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-sidebar-2 text-xs font-semibold text-white">
                    {usuario.nombres[0]?.toUpperCase()}
                  </div>
                  <div>
                    <div className="font-semibold">
                      {usuario.nombres} {usuario.apellidos}
                      {usuario.id === sesion.usuario.id && (
                        <span className="ml-1.5 text-xs font-normal text-ink-faint">(tú)</span>
                      )}
                    </div>
                    <div className="text-xs text-ink-faint">{usuario.correo}</div>
                  </div>
                </div>
              ),
            },
            {
              titulo: "Rol",
              celda: (usuario) => (
                <Insignia>{rolPorId.get(usuario.rolId) ?? usuario.rolId}</Insignia>
              ),
            },
            {
              titulo: "Estado",
              celda: (usuario) =>
                usuario.activo ? <Insignia tono="exito">Activo</Insignia> : <Insignia tono="neutro">Inactivo</Insignia>,
            },
            {
              className: "w-12 text-right",
              celda: (usuario) => (
                <details className="group relative inline-block text-left">
                  <summary className="flex h-8 w-8 cursor-pointer list-none items-center justify-center rounded-lg text-ink-soft hover:bg-bg">
                    <IconoMas className="h-[15px] w-[15px]" />
                  </summary>
                  <div className="absolute right-0 z-10 mt-1 w-56 rounded-xl border border-line bg-surface p-1.5 shadow-lg">
                    <Link
                      href={`/panel/usuarios/${usuario.id}/editar`}
                      className={`${claseItemMenu} text-ink hover:bg-bg`}
                    >
                      Editar usuario
                    </Link>
                    <Link
                      href={`/panel/usuarios/${usuario.id}/rol`}
                      className={`${claseItemMenu} text-ink hover:bg-bg`}
                    >
                      Cambiar rol
                    </Link>
                    <form action={reenviarActivacionAction}>
                      <input type="hidden" name="id" value={usuario.id} />
                      <BotonEnviar clase={`${claseItemMenu} text-ink hover:bg-bg`} textoEnviando="Reenviando…">
                        Reenviar activación
                      </BotonEnviar>
                    </form>
                    <div className="my-1 h-px bg-line" />
                    <form action={cambiarEstadoAction}>
                      <input type="hidden" name="id" value={usuario.id} />
                      <input type="hidden" name="activo" value={(!usuario.activo).toString()} />
                      <BotonEnviar clase={`${claseItemMenu} text-accent-ink hover:bg-accent-soft`} textoEnviando="Guardando…">
                        {usuario.activo ? "Desactivar usuario" : "Activar usuario"}
                      </BotonEnviar>
                    </form>
                  </div>
                </details>
              ),
            },
          ]}
        />
        </ZonaCarga>

        <Paginacion pagina={lista.pagina} totalPaginas={totalPaginas} searchParams={sp} />
      </div>
    </ListadoConCarga>
  );
}
