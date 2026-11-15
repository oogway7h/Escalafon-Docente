import os
import shutil

import pytest
from fastapi.testclient import TestClient
from PIL import Image, ImageDraw

os.environ["IA_SERVICE_TOKEN"] = "secreto"
os.environ["S3_BUCKET"] = "bucket-test"

from app import main, procesar  # noqa: E402

cliente = TestClient(main.app)
H = {"X-Token": "secreto"}


def test_exige_token():
    assert cliente.post("/trabajos", json={}).status_code == 401
    assert cliente.get("/trabajos/x", headers={"X-Token": "malo"}).status_code == 401


def test_rechaza_objetos_ajenos():
    base = {"bucket": "bucket-test", "clave": "docentes/1/a.pdf", "tipoMime": "application/pdf"}
    assert cliente.post("/trabajos", headers=H, json={**base, "bucket": "otro"}).status_code == 422
    assert cliente.post("/trabajos", headers=H, json={**base, "clave": "secretos/a.pdf"}).status_code == 422
    assert cliente.post("/trabajos", headers=H, json={**base, "tipoMime": "text/html"}).status_code == 422
    assert cliente.get("/trabajos/no-existe", headers=H).status_code == 404


@pytest.mark.skipif(shutil.which("tesseract") is None, reason="tesseract no instalado")
def test_ocr_de_imagen(tmp_path):
    ruta = tmp_path / "t.png"
    img = Image.new("RGB", (900, 200), "white")
    ImageDraw.Draw(img).text((20, 80), "RESOLUCION 2026", fill="black", font_size=48)
    img.save(ruta)
    assert "RESOLUCION" in procesar.texto_de_imagen(str(ruta)).upper()


def test_dividir_paginas_conserva_el_orden_y_descarta_el_cierre():
    assert procesar.dividir_paginas("uno\x0cdos\x0c") == ["uno", "dos"]
    assert procesar.dividir_paginas("solo una\x0c") == ["solo una"]
    assert procesar.dividir_paginas("") == []
    assert procesar.dividir_paginas("a\x0c\x0cc\x0c") == ["a", "", "c"]  # una página en blanco no mueve las demás
