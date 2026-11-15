# Unidades de trabajo — CU-02-FE

**Caso de uso:** Pantallas de recuperar, restablecer y activar.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-02 - Recuperar Contraseña`](../../../CasosDeUso/CU-02%20-%20Recuperar%20Contrase%C3%B1a.md)
**Recurso HTTP / contrato:** `/recuperar`, `/restablecer?token=`, `/activar?token=`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /recuperar | Pide el correo y muestra mensaje neutro | público |
| /restablecer | Nueva contraseña con el token del enlace | público |
| /activar | Primera contraseña de una cuenta nueva | público |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU02-FE-01](./UT-CU02-FE-01.md) | Pantalla de recuperación | — |
| [UT-CU02-FE-02](./UT-CU02-FE-02.md) | Restablecer y activar con token | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Pedir y restablecer | 01, 02 |
| A2-A4 | Mensajes de error | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Cambio de contraseña del perfil (CU-03).
