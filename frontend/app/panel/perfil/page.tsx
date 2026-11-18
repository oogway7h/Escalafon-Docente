import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Tarjeta from "@/components/common/Tarjeta";
import { obtenerPerfil, obtenerSesion, obtenerVersionFoto } from "@/lib/api";
import ContactoForm from "./ContactoForm";
import ContrasenaForm from "./ContrasenaForm";
import FotoPerfil from "./FotoPerfil";

export default async function PerfilPage() {
  const [sesion, perfil, fotoVersion] = await Promise.all([obtenerSesion(), obtenerPerfil(), obtenerVersionFoto()]);

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina titulo="Mi perfil" />

      <div className="flex flex-col gap-5 lg:flex-row lg:items-start">
        <Tarjeta titulo="Foto de perfil" className="lg:w-[340px] lg:shrink-0">
          <FotoPerfil
            nombres={sesion.usuario.nombres}
            apellidos={sesion.usuario.apellidos}
            rol={sesion.rol}
            correo={sesion.usuario.correo}
            version={fotoVersion}
          />
        </Tarjeta>

        <div className="flex min-w-0 flex-1 flex-col gap-5">
          <Tarjeta
            titulo="Datos de contacto"
            descripcion="El nombre, el correo y el rol solo puede cambiarlos un Director desde Usuarios."
          >
            <ContactoForm
              nombre={`${sesion.usuario.nombres} ${sesion.usuario.apellidos}`}
              correo={sesion.usuario.correo}
              telefono={perfil.telefono}
            />
          </Tarjeta>

          <Tarjeta
            titulo="Cambiar contraseña"
            descripcion="Mínimo 8 caracteres, una mayúscula y un dígito. Se cierran tus demás sesiones activas."
          >
            <ContrasenaForm />
          </Tarjeta>
        </div>
      </div>
    </div>
  );
}
