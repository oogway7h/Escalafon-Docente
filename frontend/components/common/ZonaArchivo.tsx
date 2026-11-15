"use client";

import { useRef, useState } from "react";

const MB = 1024 * 1024;

function tamanoLegible(bytes: number) {
  return bytes >= MB ? `${(bytes / MB).toFixed(1)} MB` : `${Math.max(1, Math.round(bytes / 1024))} KB`;
}

// Campo de archivo con arrastrar y soltar. Es un <label> sobre un <input type="file"> nativo: el clic y el
// teclado abren el selector sin código extra y el formulario envía el archivo como cualquier otro campo
// (required incluido). Al soltar, los archivos se asignan al propio input.
// maxMB puede depender del tipo (los PDF escaneados admiten más que las imágenes).
export default function ZonaArchivo({
  id,
  name,
  accept = "application/pdf,image/jpeg,image/png",
  maxMB = 8,
  maxArchivos = 1,
  required = false,
  disabled = false,
  etiqueta = "PDF, JPG o PNG",
}: {
  id?: string;
  name: string;
  accept?: string;
  maxMB?: number | ((tipo: string) => number);
  maxArchivos?: number;
  required?: boolean;
  disabled?: boolean;
  etiqueta?: string;
}) {
  const entrada = useRef<HTMLInputElement>(null);
  const [archivos, setArchivos] = useState<File[]>([]);
  const [encima, setEncima] = useState(false);
  const [error, setError] = useState("");
  const multiple = maxArchivos > 1;
  const tipos = accept.split(",").map((t) => t.trim());
  const limite = (tipo: string) => (typeof maxMB === "function" ? maxMB(tipo) : maxMB);

  // Valida y deja en el input solo lo aceptado (un FileList no se puede editar: se arma con DataTransfer).
  function asignar(lista: FileList | File[]) {
    const elegidos = [...lista].slice(0, maxArchivos);
    const malo = elegidos.find((a) => !tipos.includes(a.type));
    const grande = elegidos.find((a) => a.size > limite(a.type) * MB);
    const problema = malo
      ? `«${malo.name}» no es un archivo admitido (${etiqueta}).`
      : grande
        ? `«${grande.name}» supera los ${limite(grande.type)} MB.`
        : [...lista].length > maxArchivos
          ? `Solo se admiten ${maxArchivos} archivos: se tomaron los primeros.`
          : "";
    const validos = malo || grande ? [] : elegidos;
    const dt = new DataTransfer();
    validos.forEach((a) => dt.items.add(a));
    if (entrada.current) entrada.current.files = dt.files;
    setArchivos(validos);
    setError(problema);
  }

  function soltar(e: React.DragEvent) {
    e.preventDefault();
    setEncima(false);
    if (!disabled && e.dataTransfer.files.length > 0) asignar(e.dataTransfer.files);
  }

  return (
    <div className="flex flex-col gap-1.5">
      <label
        onDragOver={(e) => {
          e.preventDefault();
          if (!disabled) setEncima(true);
        }}
        onDragLeave={() => setEncima(false)}
        onDrop={soltar}
        className={`flex cursor-pointer flex-col items-center gap-1 rounded-lg border border-dashed px-4 py-5 text-center transition-colors duration-150 focus-within:border-accent ${
          encima ? "border-accent bg-accent-soft" : "border-line bg-surface hover:border-ink-faint"
        } ${disabled ? "cursor-not-allowed opacity-60" : ""}`}
      >
        {/* sr-only (no hidden): el input sigue siendo enfocable y conserva la validación required del navegador. */}
        <input
          ref={entrada}
          id={id}
          name={name}
          type="file"
          accept={accept}
          multiple={multiple}
          required={required}
          disabled={disabled}
          onChange={(e) => asignar(e.target.files ?? [])}
          className="sr-only"
        />
        {archivos.length === 0 ? (
          <>
            <span className="text-sm text-ink">
              Arrastra {multiple ? "los archivos" : "el archivo"} aquí o <span className="font-semibold text-accent-ink underline">elige {multiple ? "archivos" : "uno"}</span>
            </span>
            <span className="text-xs text-ink-faint">
              {etiqueta}
              {typeof maxMB === "number" && ` · hasta ${maxMB} MB`}
              {multiple && ` · máximo ${maxArchivos}`}
            </span>
          </>
        ) : (
          <ul className="flex w-full flex-col gap-1 text-left text-sm">
            {archivos.map((a) => (
              <li key={`${a.name}-${a.size}`} className="flex items-baseline justify-between gap-3">
                <span className="min-w-0 truncate font-semibold text-ink">{a.name}</span>
                <span className="shrink-0 text-xs text-ink-faint">{tamanoLegible(a.size)}</span>
              </li>
            ))}
            <li className="text-xs text-ink-faint">Haz clic o arrastra de nuevo para cambiar{multiple ? "los" : "lo"}.</li>
          </ul>
        )}
      </label>
      {error && (
        <p role="alert" className="text-xs text-accent-ink">
          {error}
        </p>
      )}
    </div>
  );
}
