import { exigirPermiso } from "@/lib/api";
import CategoriaForm from "../CategoriaForm";
import { guardarCategoriaAction } from "../actions";

export default async function NuevaCategoriaPage() {
  await exigirPermiso("CATEGORIA_GESTIONAR");
  return <CategoriaForm accion={guardarCategoriaAction.bind(null, null)} titulo="Nueva categoría" />;
}
