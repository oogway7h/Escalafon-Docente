import { exigirPermiso } from "@/lib/api";
import { listarCategorias } from "@/lib/api-tramites";
import RadicarForm from "./RadicarForm";

export default async function RadicarPage() {
  await exigirPermiso("TRAMITE_RADICAR");
  const categorias = await listarCategorias();
  return <RadicarForm categorias={categorias.filter((c) => c.activa)} />;
}
