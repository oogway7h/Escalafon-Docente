import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Tarjeta from "@/components/common/Tarjeta";
import { obtenerPerfil, obtenerSesion } from "@/lib/api";
import ContactoForm from "./ContactoForm";
import ContrasenaForm from "./ContrasenaForm";

export default async function PerfilPage() {
  const [sesion, perfil] = await Promise.all([obtenerSesion(), obtenerPerfil()]);

  const datos = [
    ["Nombre", `${sesion.usuario.nombres} ${sesion.usuario.apellidos}`],
    ["Correo", sesion.usuario.correo],
    ["Rol", sesion.rol],
  ];

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina titulo="Mi perfil" />

      <div className="flex flex-col gap-5 lg:flex-row lg:items-start">
        <Tarjeta
          titulo="Datos de contacto"
          descripcion="El correo y el rol solo puede cambiarlos un Director desde Usuarios."
          className="flex-1"
        >
          <dl className="mb-5 flex flex-col gap-3 text-sm">
            {datos.map(([etiqueta, valor]) => (
              <div key={etiqueta}>
                <dt className="text-xs font-semibold uppercase tracking-wide text-ink-faint">{etiqueta}</dt>
                <dd className="text-ink">{valor}</dd>
              </div>
            ))}
          </dl>
          <ContactoForm telefono={perfil.telefono} />
        </Tarjeta>

        <Tarjeta
          titulo="Cambiar contraseña"
          descripcion="Mínimo 8 caracteres, una mayúscula y un dígito. Se cierran tus demás sesiones activas."
          className="flex-1"
        >
          <ContrasenaForm />
        </Tarjeta>
      </div>
    </div>
  );
}
