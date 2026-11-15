import type { Metadata } from "next";
import AuthDividido from "@/components/layout/AuthDividido";
import LoginForm from "./LoginForm";

export const metadata: Metadata = {
  title: "Iniciar sesión — Escalafón Docente",
};

export default async function LoginPage() {
  return (
    <AuthDividido
    
      heroTitulo="Gestión de flujos de trabajo y rastreo documental"
      heroTexto="Sube un peldaño a la vez: cada trámite, documento y legajo docente en un solo lugar, con trazabilidad completa."
      heroDetalle={
        <div className="flex flex-col gap-2">
          <div className="flex gap-1.5">
            <span className="h-1.5 w-[34px] rounded-[3px] bg-white/16" />
            <span className="h-1.5 w-[34px] rounded-[3px] bg-white/16" />
            <span className="h-1.5 w-[34px] rounded-[3px] bg-[#EF5B6E]" />
            <span className="h-1.5 w-[34px] rounded-[3px] bg-white/16" />
          </div>
          <p className="text-xs text-white/60">Categoría Titular · en revisión</p>
          
        </div>
      }
      titulo="Iniciar sesión"
      descripcion="Ingresa con tu correo institucional para continuar."
      pie={
        <>
          ¿Problemas para acceder? Escribe a <span className="font-semibold text-ink">cdotsuboa126@ficct.uagrm.edu.bo</span>
        </>
      }
    >
      <LoginForm />
    </AuthDividido>
    
  );
}
