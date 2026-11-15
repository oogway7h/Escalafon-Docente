import { redirect } from "next/navigation";

// Atender solicitudes ahora vive en la pantalla única de solicitudes y custodia.
export default function SolicitudesPendientesPage() {
  redirect("/panel/custodia");
}
