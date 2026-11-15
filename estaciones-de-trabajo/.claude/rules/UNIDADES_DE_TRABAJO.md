---
globs: UnidadesTrabajo/**/*
---

# Unidades de Trabajo (UT) — las "estaciones de trabajo"

Una UT es la **unidad mínima de trabajo que se le pide a una persona o a una IA**: una tarea acotada,
con alcance y fuera-de-alcance explícitos y una validación binaria. Es el **plan** del ciclo SDD: el
CU dice *qué*, la UT dice *dónde y cómo*, el código la cumple y el commit la cita.

Cuantas más decisiones queden tomadas en la UT, menos errores por interpretación al implementarla (y
menos alucinaciones cuando la implementa una IA).

## Cuándo se crean

Solo cuando el CU está escrito y completo en `CasosDeUso/`. Las UT de la web se escribieron de forma
**retroactiva** (Estado: Hecho, citan archivos que ya existen); las de mobile, **antes del código**.

## Estructura

```
UnidadesTrabajo/
└── CU-{NN} - {Nombre}/            <- mismo nombre que el archivo del CU
    ├── UT-BACKEND/   README.md + UT-CU{NN}-BE-01.md …   (repo backend/, Go)
    ├── UT-FRONTEND/  README.md + UT-CU{NN}-FE-01.md …   (repo frontend/, Next.js)
    └── UT-MOBILE/    README.md + UT-CU{NN}-MOB-01.md …  (repo ded_Mobile/, Flutter) — solo si el CU toca la app
```

Una carpeta por capa, con numeración propia. Si un CU no toca una capa, la carpeta no existe.

## Código

`UT-CU{NN}-{BE|FE|MOB}-{NN}`, único en todo el proyecto; se cita en el commit
(`feat: UT-CU12-MOB-02 pantalla de seguimiento`). No se reusa un número.

## README de cada capa

```markdown
# Unidades de trabajo — CU-NN-{CAPA}

**Caso de uso:** resumen de una línea en esta capa
**Proyecto:** `backend` | `frontend` | `ded_Mobile`
**Fuente CU:** [link](../../../CasosDeUso/CU-NN%20-%20Nombre.md)
**Recurso HTTP / contrato:** endpoints o pantallas
**Estado:** Hecho | Pendiente

## Antes de implementar cualquier UT
## Decisiones de diseño         (solo si hubo que decidir algo que el CU no resolvía)
## Contrato de datos / recurso HTTP
## Orden de ejecución           (tabla Código | Nombre | Depende de  +  **Cadena:** 01 → 02 → 03)
## Cobertura del caso de uso    (CA o flujo → UTs)
## Follow-ups
## Qué NO incluyen estas UT
```

## Cada `UT-CU{NN}-{CAPA}-{NN}.md`

```markdown
# UT-CUNN-CAPA-NN — Nombre

| Campo | Valor |
|---|---|
| **Código** | `UT-CUNN-CAPA-NN` |
| **Nombre** | … |
| **Caso de uso** | CU-NN-CAPA |
| **Proyecto** | backend / frontend / ded_Mobile |
| **Depende de** | UT del mismo CU, o "—" |
| **Estado** | Hecho / Pendiente |

## Antes de implementar   (rules puntuales para ESTA tarea)
## Solicitud de negocio     (el porqué, 1-2 frases)
## Objetivo                 (qué queda construido, 1 frase)
## Alcance                  (exacto; termina con "Archivos esperados (crear/tocar)")
## Fuera de alcance (no tocar)
## Resultado esperado
## Validación               (pasos verificables: comando, test o acción con resultado)
```

## Reglas estrictas

1. **Una UT nunca construye trabajo de otro CU**, aunque sea trivial. El CU futuro *edita* el archivo
   (lo declara en su Alcance como "tocar"), no lo encuentra ya hecho. Así cada CU se verifica solo.
2. **Alcance y Fuera de alcance son obligatorios.** Sin "fuera de alcance" la UT está incompleta.
3. **Granularidad:** 1 a 4 archivos estrechamente relacionados. Si toca dominio, persistencia y HTTP a
   la vez, son varias UT.
4. **Dependencias** solo hacia UT del mismo CU o hacia algo que ya existe en el repo.
5. **Toda UT termina en algo verificable**: `go test …`, `flutter test`, `pnpm build`, o pasos manuales
   con su resultado esperado.

## Cómo se le pide una UT a la IA

```
Implementa UT-CU12-MOB-02. Lee primero la UT, el README de su capa y las rules que cita.
Respeta "Fuera de alcance". Termina corriendo la Validación y reporta el resultado.
```

Cada pedido y su resultado se anota en [BITACORA_IA.md](../../BITACORA_IA.md).
