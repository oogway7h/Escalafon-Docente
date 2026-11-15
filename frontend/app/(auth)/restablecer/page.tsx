import type { Metadata } from "next";
import AuthTarjeta from "@/components/layout/AuthTarjeta";
import ContrasenaConTokenForm from "../ContrasenaConTokenForm";
import { restablecerContrasena } from "./actions";

export const metadata: Metadata = {
  title: "Restablecer contraseña — Escalafón",
};

export default async function RestablecerPage({
  searchParams,
}: {
  searchParams: Promise<{ token?: string }>;
}) {
  const { token } = await searchParams;

  return (
    <AuthTarjeta
      titulo="Restablecer contraseña"
      descripcion="Elige una nueva contraseña. El enlace que recibiste por correo es válido por 30 minutos y de un solo uso."
    >
      <ContrasenaConTokenForm
        accion={restablecerContrasena}
        token={token ?? ""}
        etiquetaNueva="Nueva contraseña"
        etiquetaConfirmar="Confirmar nueva contraseña"
        textoBoton="Guardar nueva contraseña"
        textoEnviando="Guardando…"
        sinToken={
          <>
            Este enlace no incluye un token válido. Solicita uno nuevo desde{" "}
            <a href="/recuperar" className="font-semibold underline">
              recuperar contraseña
            </a>
            .
          </>
        }
      />
    </AuthTarjeta>
  );
}
