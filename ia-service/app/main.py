"""Servicio de OCR y compresión de documentos del Departamento de Escalafón Docente (CU-15).

El backend (Go) pide un trabajo con POST /trabajos y consulta su avance con GET /trabajos/{id}.
El archivo no viaja por HTTP: el servicio lo lee y escribe directamente en S3.
"""
import hmac
import os
import time
import uuid
from concurrent.futures import ThreadPoolExecutor

from fastapi import Depends, FastAPI, Header, HTTPException
from pydantic import BaseModel

from . import procesar

app = FastAPI(title="ia-service")
# ponytail: trabajos en memoria y un pool pequeño; si el servicio reinicia, el backend reanuda los pendientes.
trabajos: dict[str, dict] = {}
TTL_TRABAJO = 3600
pool = ThreadPoolExecutor(max_workers=int(os.getenv("MAX_TRABAJOS", "1")))


def autenticar(x_token: str = Header(default="")):
    esperado = os.getenv("IA_SERVICE_TOKEN", "")
    if not esperado or not hmac.compare_digest(x_token, esperado):
        raise HTTPException(401, "token inválido")


class Solicitud(BaseModel):
    bucket: str
    clave: str
    tipoMime: str
    claveSalida: str | None = None


def ejecutar(id_: str, s: Solicitud):
    trabajos[id_]["estado"] = "PROCESANDO"
    try:
        trabajos[id_].update(procesar.procesar(s.bucket, s.clave, s.tipoMime, s.claveSalida), estado="LISTO")
    except Exception as e:  # el motivo se devuelve al backend, que lo muestra al analista
        trabajos[id_].update(estado="ERROR", error=str(e)[:290])
    trabajos[id_]["terminado"] = time.monotonic()


def purgar():
    # Los trabajos terminados retienen el texto completo: se descartan pasada una hora.
    for k in [k for k, v in trabajos.items() if time.monotonic() - v.get("terminado", time.monotonic()) > TTL_TRABAJO]:
        trabajos.pop(k, None)


@app.get("/salud")
def salud():
    return {"ok": True}


@app.post("/trabajos", status_code=202, dependencies=[Depends(autenticar)])
def crear(s: Solicitud):
    # Solo se procesan objetos del bucket configurado y de la carpeta de docentes.
    if s.bucket != os.getenv("S3_BUCKET") or not s.clave.startswith("docentes/"):
        raise HTTPException(422, "objeto no permitido")
    if s.tipoMime not in procesar.TIPOS:
        raise HTTPException(422, "tipo de archivo no permitido")
    if s.claveSalida and not (s.claveSalida.startswith("docentes/") and s.claveSalida.endswith("-comprimido.pdf")):
        raise HTTPException(422, "destino no permitido")  # no se puede pisar el original de otro documento
    purgar()
    id_ = str(uuid.uuid4())
    trabajos[id_] = {"estado": "EN_COLA", "texto": "", "paginas": [], "error": "", "hashSha256": "", "tamanoOriginal": 0, "tamanoFinal": 0, "claveComprimida": ""}
    pool.submit(ejecutar, id_, s)
    return {"id": id_}


@app.get("/trabajos/{id_}", dependencies=[Depends(autenticar)])
def estado(id_: str):
    if id_ not in trabajos:
        raise HTTPException(404, "trabajo no encontrado")
    return trabajos[id_]
