import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Solo comprueba que exista la cookie: la validación real (firma, vigencia, versión de sesión)
// la hace la API en Go en cada petición (sección 6.1 de Inicializacion.md).
export function proxy(request: NextRequest) {
  if (!request.cookies.has("sesion")) {
    return NextResponse.redirect(new URL("/login", request.url));
  }
  return NextResponse.next();
}

export const config = {
  matcher: ["/panel/:path*"],
};
