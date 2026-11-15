import { exigirPermiso } from "@/lib/api";
import { obtenerCarpeta } from "@/lib/api-legajos";
import LegajoForm from "../../LegajoForm";
import { guardarLegajoAction } from "../../actions";

export default async function EditarLegajoPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  await exigirPermiso("LEGAJO_GESTIONAR");
  const carpeta = await obtenerCarpeta(id);
  return <LegajoForm accion={guardarLegajoAction.bind(null, id)} carpeta={carpeta} titulo={`Editar legajo ${carpeta.codigo}`} />;
}
