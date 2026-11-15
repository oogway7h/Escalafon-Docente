import "server-only";
import { cookies } from "next/headers";

const NOMBRE_COOKIE = "sesion";

export async function leerTokenSesion(): Promise<string | undefined> {
  const cookieStore = await cookies();
  return cookieStore.get(NOMBRE_COOKIE)?.value;
}

export async function escribirTokenSesion(token: string, maxAgeSegundos: number) {
  const cookieStore = await cookies();
  cookieStore.set(NOMBRE_COOKIE, token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: maxAgeSegundos,
  });
}

export async function borrarTokenSesion() {
  const cookieStore = await cookies();
  cookieStore.delete(NOMBRE_COOKIE);
}
