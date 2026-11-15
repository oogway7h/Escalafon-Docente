# Unidades de trabajo — CU-23-MOB

**Caso de uso:** Archivo registra la devolución con el estado material.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-23 - Registrar Devolución de Carpeta (Check-in)`](../../../CasosDeUso/CU-23%20-%20Registrar%20Devoluci%C3%B3n%20de%20Carpeta%20%28Check-in%29.md)
**Recurso HTTP / contrato:** `POST /api/legajos/:id/devolver`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /legajos/:id/devolver `{estadoMaterial, observacion}` | 200 | 422 `OBSERVACION_REQUERIDA`, `ESTADO_MATERIAL_INVALIDO`; 409 `LEGAJO_NO_PRESTADO` | `CUSTODIA_REGISTRAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU23-MOB-01](./UT-CU23-MOB-01.md) | `devolver` en `CarpetasService` | — |
| [UT-CU23-MOB-02](./UT-CU23-MOB-02.md) | Pestaña Prestadas y hoja de devolución | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Devolución | 01, 02 |
| CA02 | Observación obligatoria | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Foto de evidencia (se quitó del CU).
