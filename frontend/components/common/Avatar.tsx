"use client";

import { useState } from "react";

// Foto de perfil o, si no hay (o no carga), las iniciales sobre el azul del panel.
// version: la que devuelve la API; cambia con cada foto y evita servir una vieja desde la caché del navegador.
export default function Avatar({
  nombres,
  apellidos,
  version,
  tamano = 32,
  vistaPrevia,
}: {
  nombres: string;
  apellidos: string;
  version: number | null;
  tamano?: number;
  vistaPrevia?: string; // blob: de una foto elegida y aún no guardada
}) {
  const [fallo, setFallo] = useState<number | null>(null);
  const iniciales = `${nombres[0] ?? ""}${apellidos[0] ?? ""}`.toUpperCase() || "?";
  const src = vistaPrevia ?? (version !== null && fallo !== version ? `/panel/perfil/foto?v=${version}` : null);

  return src ? (
    // eslint-disable-next-line @next/next/no-img-element -- ruta propia autenticada: no pasa por el optimizador de next/image
    <img
      src={src}
      alt=""
      width={tamano}
      height={tamano}
      onError={() => setFallo(version)}
      className="shrink-0 rounded-full object-cover"
      style={{ width: tamano, height: tamano }}
    />
  ) : (
    <span
      aria-hidden="true"
      className="flex shrink-0 items-center justify-center rounded-full bg-sidebar font-semibold text-white"
      style={{ width: tamano, height: tamano, fontSize: tamano * 0.38 }}
    >
      {iniciales}
    </span>
  );
}
