import { exigirPermiso } from "@/lib/api";
import LegajoForm from "../LegajoForm";
import { guardarLegajoAction } from "../actions";

export default async function NuevoLegajoPage() {
  await exigirPermiso("LEGAJO_GESTIONAR");
  return <LegajoForm accion={guardarLegajoAction.bind(null, null)} titulo="Nuevo legajo" />;
}
