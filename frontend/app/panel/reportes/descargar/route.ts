import { NextResponse } from "next/server";
import { generarReporte } from "@/lib/api-indicadores";

// Pide el reporte a la API y lo reenvía como descarga (el navegador nunca habla con la API).
// Desde el formulario con JavaScript (cabecera x-descarga) los errores vuelven como JSON; sin JavaScript,
// se redirige a /panel/reportes con el código de error en la URL.
export async function POST(request: Request) {
  const datos = await request.formData();
  const texto = (clave: string) => {
    const valor = datos.get(clave);
    return typeof valor === "string" && valor.trim() ? valor.trim() : undefined;
  };
  const tipo = texto("tipo") ?? "";
  const categoria = Number.parseInt(texto("categoriaId") ?? "", 10);

  const respuesta = await generarReporte({
    tipo,
    formato: texto("formato") ?? "",
    desde: texto("desde"),
    hasta: texto("hasta"),
    estado: texto("estado"),
    tipoTramite: texto("tipoTramite"),
    categoriaId: Number.isFinite(categoria) ? categoria : undefined,
    estadoLegajo: texto("estadoLegajo"),
  });

  if (!respuesta.ok) {
    const cuerpo = await respuesta.json().catch(() => null);
    const codigo: string = cuerpo?.codigo ?? "ERROR_REPORTE";
    if (request.headers.get("x-descarga") === "fetch") {
      return NextResponse.json({ codigo }, { status: respuesta.status });
    }
    const destino = new URL("/panel/reportes", request.url);
    destino.searchParams.set("error", codigo);
    if (tipo) destino.searchParams.set("tipo", tipo);
    return NextResponse.redirect(destino, 303);
  }

  const cabeceras = new Headers({ "cache-control": "private, no-store" });
  for (const nombre of ["content-type", "content-disposition", "content-length", "x-reporte-filas"]) {
    const valor = respuesta.headers.get(nombre);
    if (valor) cabeceras.set(nombre, valor);
  }
  return new Response(respuesta.body, { status: 200, headers: cabeceras });
}
