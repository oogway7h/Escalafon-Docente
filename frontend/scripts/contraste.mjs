// Comprueba que los pares texto/fondo del sistema de color cumplan el piso de contraste de
// WCAG 2.1 AA (4.5:1 texto normal, 3:1 texto grande y bordes de control), en los dos modos.
// Lee los valores reales de app/globals.css, así que falla si alguien baja un token.
//
//   node scripts/contraste.mjs
//
// Los pares son los que el código usa de verdad: si aparece una combinación nueva
// (p. ej. bg-X text-Y), añádela aquí.
import { readFileSync } from "node:fs";
import assert from "node:assert/strict";

const css = readFileSync(new URL("../app/globals.css", import.meta.url), "utf8");

// El bloque :root trae el modo claro; .dark lo sobreescribe.
function tokens(selector) {
  const bloque = css.match(new RegExp(`${selector}\\s*\\{([^}]*)\\}`))?.[1];
  assert.ok(bloque, `no se encontró el bloque ${selector} en globals.css`);
  return Object.fromEntries([...bloque.matchAll(/--color-([\w-]+):\s*(#[0-9a-fA-F]{6})/g)].map(([, k, v]) => [k, v]));
}

const claro = tokens(":root");
const oscuro = { ...claro, ...tokens("\\.dark") };

const canal = (c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4);
const luminancia = (hex) =>
  [1, 3, 5]
    .map((i) => canal(parseInt(hex.slice(i, i + 2), 16) / 255))
    .reduce((suma, v, i) => suma + [0.2126, 0.7152, 0.0722][i] * v, 0);

function contraste(a, b) {
  const [x, y] = [luminancia(a), luminancia(b)];
  return (Math.max(x, y) + 0.05) / (Math.min(x, y) + 0.05);
}

// [texto, fondo, piso, dónde se usa, modos?]. "#ffffff" es literal (text-white sobre bg-accent).
// modos limita el par a un modo cuando la clase solo se aplica ahí (p. ej. `dark:` lo cambia).
const pares = [
  ["ink", "bg", 4.5, "texto principal sobre el lienzo"],
  ["ink", "surface", 4.5, "texto principal en tarjetas"],
  ["ink-soft", "bg", 4.5, "descripciones (EncabezadoPagina, Tarjeta)"],
  ["ink-soft", "surface", 4.5, "descripciones en tarjetas"],
  ["ink-faint", "bg", 4.5, "Insignia neutro, pie de la portada"],
  ["ink-faint", "surface", 4.5, "cabeceras de Tabla, ayuda de Campo, estado vacío"],
  ["accent-ink", "surface", 4.5, "enlaces y errores en tarjetas"],
  ["accent-ink", "bg", 4.5, "enlaces sobre el lienzo"],
  ["accent-ink", "accent-soft", 4.5, "Alerta error/aviso, Insignia acento"],
  // Alerta aviso usa `text-accent-dark dark:text-accent-ink`: este par solo existe en claro.
  ["accent-dark", "accent-soft", 4.5, "Alerta aviso", ["claro"]],
  ["success", "surface", 4.5, "Alerta éxito, estados en tabla"],
  ["success", "bg", 4.5, "estados sobre el lienzo"],
  ["#ffffff", "accent", 4.5, "texto del botón primario"],
  ["#ffffff", "accent-dark", 4.5, "botón primario en hover"],
  ["#ffffff", "sidebar", 4.5, "barra lateral (extremo oscuro del degradado)"],
  ["#ffffff", "sidebar-2", 4.5, "barra lateral (extremo claro del degradado)"],
  ["accent", "surface", 3, "borde de control enfocado (focus:border-accent)"],
];

const fallos = [];
for (const [modo, t] of [
  ["claro", claro],
  ["oscuro", oscuro],
]) {
  for (const [textoClave, fondoClave, piso, donde, modos] of pares) {
    if (modos && !modos.includes(modo)) continue;
    const texto = textoClave.startsWith("#") ? textoClave : t[textoClave];
    const fondo = fondoClave.startsWith("#") ? fondoClave : t[fondoClave];
    assert.ok(texto, `token --color-${textoClave} ausente en modo ${modo}`);
    assert.ok(fondo, `token --color-${fondoClave} ausente en modo ${modo}`);
    const r = contraste(texto, fondo);
    const linea = `${modo.padEnd(6)} ${textoClave} sobre ${fondoClave}`.padEnd(46);
    if (r < piso) fallos.push(`✗ ${linea} ${r.toFixed(2)}:1 < ${piso}:1  (${donde})`);
    else console.log(`✓ ${linea} ${r.toFixed(2)}:1`);
  }
}

assert.deepEqual(fallos, [], `\n${fallos.join("\n")}\n`);
console.log("\nTodos los pares del sistema de color cumplen WCAG 2.1 AA.");
