import type { Metadata } from "next";
import AuthTarjeta from "@/components/layout/AuthTarjeta";
import ContrasenaConTokenForm from "../ContrasenaConTokenForm";
import { activarCuenta } from "./actions";

export const metadata: Metadata = {
  title: "Activar cuenta — Escalafón",
};

export default async function ActivarPage({
  searchParams,
}: {
  searchParams: Promise<{ token?: string }>;
}) {
  const { token } = await searchParams;

  return (
    <AuthTarjeta
      titulo="Activa tu cuenta"
      descripcion="Elige tu contraseña para terminar de crear tu cuenta. El enlace vence a las 24 horas y solo sirve una vez."
    >
      <ContrasenaConTokenForm
        accion={activarCuenta}
        token={token ?? ""}
        etiquetaNueva="Crea tu contraseña"
        etiquetaConfirmar="Confirmar contraseña"
        textoBoton="Activar mi cuenta"
        textoEnviando="Activando…"
        sinToken="Este enlace no incluye un token válido. Pide a un Director que reenvíe tu correo de activación."
      />
    </AuthTarjeta>
  );
}
