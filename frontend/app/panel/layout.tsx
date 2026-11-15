import PanelShell from "@/components/layout/PanelShell";
import { obtenerSesion } from "@/lib/api";

export default async function PanelLayout({ children }: { children: React.ReactNode }) {
  const sesion = await obtenerSesion();
  return <PanelShell sesion={sesion}>{children}</PanelShell>;
}
