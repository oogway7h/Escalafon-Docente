import type { Metadata } from "next";
import AuthDividido from "@/components/layout/AuthDividido";
import RecuperarForm from "./RecuperarForm";

export const metadata: Metadata = {
  title: "Recuperar contraseña — Escalafón",
};

export default async function RecuperarPage() {
  return (
    <AuthDividido
      heroTitulo="¿Olvidaste tu contraseña?"
      heroTexto="Te enviamos un enlace de un solo uso, válido por 30 minutos, para que elijas una nueva."
      titulo="Recuperar contraseña"
      descripcion="Ingresa tu correo institucional y te enviaremos el enlace."
      pie={
        <a href="/login" className="font-semibold text-accent-ink hover:text-accent-dark hover:underline">
          Volver a iniciar sesión
        </a>
      }
    >
      <RecuperarForm />
    </AuthDividido>
  );
}
