# UT-CU02-BE-02 — SolicitarRecuperacion y RestablecerContrasena

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU02-BE-02` |
| **Nombre** | SolicitarRecuperacion y RestablecerContrasena |
| **Caso de uso** | CU-02-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU02-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Restablecer debe dejar fuera a cualquiera que tuviera una sesión con la contraseña vieja.

## Objetivo

Casos de uso que emiten el enlace, validan el token, aplican la política y suben `version_sesion`.

## Alcance

- Respuesta neutra si el correo no existe.
- Mailer SMTP con plantilla y enlace a `APP_URL_FRONTEND`.
- Activación reutiliza el mismo flujo con propósito ACTIVACION.
- Tests: token vencido, usado, política, revocación.

Archivos esperados (crear/tocar):

- `backend/internal/application/auth/auth.go`
- `backend/internal/infrastructure/mail/mailer.go`
- `backend/internal/application/auth/auth_test.go`

## Fuera de alcance (no tocar)

- Handler (UT 03).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Casos de uso que emiten el enlace, validan el token, aplican la política y suben `version_sesion`.

## Validación

1. `go test ./internal/application/auth -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
