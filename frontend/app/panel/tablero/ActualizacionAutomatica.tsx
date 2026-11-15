"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

// CU-25 pide refrescar los indicadores al menos cada quince minutos: la página se vuelve a pedir sola
// (sin recargar el documento ni perder el periodo elegido).
export default function ActualizacionAutomatica({ minutos }: { minutos: number }) {
  const router = useRouter();
  useEffect(() => {
    const t = setInterval(() => router.refresh(), minutos * 60_000);
    return () => clearInterval(t);
  }, [minutos, router]);
  return null;
}
