import AcuerdoUsuario from "@/components/common/AcuerdoUsuario";
import PanelShell from "@/components/layout/PanelShell";
import { obtenerSesion, obtenerVersionFoto } from "@/lib/api";

export default async function PanelLayout({ children }: { children: React.ReactNode }) {
  const sesion = await obtenerSesion();
  // Primer inicio de sesión (o nueva versión del acuerdo): sin aceptarlo no se ve el panel.
  if (!sesion.terminosAceptados) return <AcuerdoUsuario nombre={sesion.usuario.nombres} />;
  const fotoVersion = await obtenerVersionFoto();
  return (
    <PanelShell sesion={sesion} fotoVersion={fotoVersion}>
      {children}
    </PanelShell>
  );
}
