import AcuerdoUsuario from "@/components/common/AcuerdoUsuario";
import PanelShell from "@/components/layout/PanelShell";
import { obtenerSesion } from "@/lib/api";

export default async function PanelLayout({ children }: { children: React.ReactNode }) {
  const sesion = await obtenerSesion();
  // Primer inicio de sesión (o nueva versión del acuerdo): sin aceptarlo no se ve el panel.
  if (!sesion.terminosAceptados) return <AcuerdoUsuario nombre={sesion.usuario.nombres} />;
  return <PanelShell sesion={sesion}>{children}</PanelShell>;
}
