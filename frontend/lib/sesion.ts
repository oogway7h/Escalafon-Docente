import "server-only";
import { cookies } from "next/headers";

const NOMBRE_COOKIE = "sesion";

export async function leerTokenSesion(): Promise<string | undefined> {
  const cookieStore = await cookies();
  return cookieStore.get(NOMBRE_COOKIE)?.value;
}

// Sin maxAgeSegundos la cookie es de sesión: el navegador la descarta al cerrarse.
export async function escribirTokenSesion(token: string, maxAgeSegundos?: number) {
  const cookieStore = await cookies();
  cookieStore.set(NOMBRE_COOKIE, token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production" && process.env.COOKIE_SECURE !== "false", // COOKIE_SECURE=false: despliegue con HTTP plano
    sameSite: "lax",
    path: "/",
    maxAge: maxAgeSegundos,
  });
}

export async function borrarTokenSesion() {
  const cookieStore = await cookies();
  cookieStore.delete(NOMBRE_COOKIE);
}
