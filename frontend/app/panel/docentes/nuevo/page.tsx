import { exigirPermiso } from "@/lib/api";
import { listarCategorias, listarFacultades } from "@/lib/api-tramites";
import DocenteForm from "../DocenteForm";
import { guardarDocenteAction } from "../actions";

export default async function NuevoDocentePage() {
  await exigirPermiso("DOCENTE_GESTIONAR");
  const [categorias, facultades] = await Promise.all([listarCategorias(), listarFacultades()]);
  return (
    <DocenteForm accion={guardarDocenteAction.bind(null, null)} categorias={categorias} facultades={facultades} titulo="Nuevo docente" />
  );
}
