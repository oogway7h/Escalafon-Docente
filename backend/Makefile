include .env
export

.PHONY: dev test arch lint migrate-diff migrate-new migrate-hash migrate-apply migrate-status seed import-docentes import-docentes-dry

dev:
	go run ./cmd/api

test: arch
	go test ./... -count=1

arch:
	@! go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/domain/... \
	  | grep -E 'gorm|fiber|redis|jwt|/infrastructure|/interfaces|/application' \
	  || (echo "Violación: domain importa una capa externa" && exit 1)
	@! go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/application/... \
	  | grep -E 'gorm|fiber|redis|/infrastructure|/interfaces' \
	  || (echo "Violación: application importa una capa externa" && exit 1)

lint:
	go vet ./...

migrate-diff:
	@if [ -z "$(name)" ]; then echo "Falta name=: describe el cambio en snake_case, p. ej. make migrate-diff name=agregar_tabla_notificaciones"; exit 1; fi
	atlas migrate diff $(name) --env local

migrate-new:
	@if [ -z "$(name)" ]; then echo "Falta name=: describe el cambio en snake_case, p. ej. make migrate-diff name=agregar_tabla_notificaciones"; exit 1; fi
	atlas migrate new $(name) --dir file://migrations

migrate-hash:
	atlas migrate hash --dir file://migrations

migrate-apply:
	atlas migrate apply --env local --revisions-schema public

migrate-status:
	atlas migrate status --env local --revisions-schema public

seed:
	go run ./cmd/seed

# Carga única del padrón docente desde un CSV: make import-docentes file=../archivos/docsAyuda/datos.csv
import-docentes:
	go run ./cmd/importdocentes -file="$(file)"

# Igual pero sin escribir en la base de datos: solo informa y genera reporte-importacion.txt
import-docentes-dry:
	go run ./cmd/importdocentes -file="$(file)" -dry-run
