---
globs: ../backend/**/*
---

# Arquitectura Backend (Go)

Repo `backend/`: API REST en Go con Fiber v2, GORM + Postgres (Neon), Redis, JWT HS256 + bcrypt, Atlas
para migraciones, SMTP (Mailpit en desarrollo), S3 y `ia-service` (OCR). Detalle vivo en
`backend/docs/estado-del-backend.md`.

## Capas (Clean Architecture)

```
interfaces/http  →  application  →  domain
infrastructure   →  application  →  domain
```

| Capa | Qué va | Ejemplo |
|---|---|---|
| `internal/domain/<agregado>` | Entidades, reglas, errores de dominio, puertos de repositorio. Sin GORM/Fiber/Redis. | `domain/tramite/tramite.go` (máquina de estados) |
| `internal/application/<módulo>` | **Un struct por caso de uso** con `Ejecutar(ctx, entrada)`. Puertos externos en `ports`. | `application/tramites/flujo.go` |
| `internal/infrastructure` | Implementaciones: `persistence` (GORM), `cache` (Redis), `security`, `mail`, `storage` (S3), `iaservice`, `reportes`. | `persistence/repositorio_tramite.go` |
| `internal/interfaces/http` | Router, middlewares (`Autenticar`, `RequierePermiso`, `RequiereAlguno`), handlers delgados, DTO, `apierror`. | `handlers/tramites.go` |

`cmd/api/main.go` es el **único composition root**. `make arch` falla si `domain` o `application`
importan capas externas.

## Convenciones que una UT de backend debe respetar

- Handler = traducir HTTP ↔ caso de uso. Nada de reglas de negocio en el handler.
- Errores de negocio con código estable `{codigo, mensaje, detalles?}` (`apierror.Mapear`).
- Autorización **por permiso** (`rol.Permiso*`), nunca por nombre de rol. Un permiso nuevo va al seed.
- Toda acción que cambia estado deja **auditoría**. Las bitácoras son de solo anexión (trigger).
- Listados paginan con `?pagina=N&tamano=` → `{datos,total,pagina,tamano}`.
- Migración nueva: `make migrate-diff name=<snake_case>`; revisar que no borre índices manuales.
- Tests: caso de uso con dobles de `application/testdobles`; repositorio contra `DATABASE_URL_TEST`.

Validación estándar de una UT de backend: `go test ./internal/<paquete> -count=1` + `make arch`.
