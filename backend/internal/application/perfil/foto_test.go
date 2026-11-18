package perfil_test

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"testing"
	"time"

	"escalafon/backend/internal/application/perfil"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/usuario"
)

func pngDe(w, h int) []byte {
	var b bytes.Buffer
	_ = png.Encode(&b, image.NewRGBA(image.Rect(0, 0, w, h)))
	return b.Bytes()
}

func TestSubirFoto_ValidaElContenidoRealYGuarda(t *testing.T) {
	fotos, almacen := testdobles.NewFotoRepoFalso(), testdobles.NewAlmacenFalso()
	caso := &perfil.SubirFoto{Fotos: fotos, Almacen: almacen, Auditoria: &testdobles.AuditoriaRepoFalso{}, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}
	ctx := context.Background()

	f, err := caso.Ejecutar(ctx, "u1", pngDe(10, 10))
	if err != nil || f.TipoMime != "image/png" || f.Clave != "usuarios/u1/foto" {
		t.Fatalf("PNG válido: foto=%+v err=%v", f, err)
	}
	if len(almacen.Guardados["usuarios/u1/foto"]) == 0 {
		t.Fatal("la foto no se guardó en el almacén")
	}

	for nombre, datos := range map[string][]byte{
		"texto disfrazado de imagen": []byte("no soy una imagen"),
		"vacío":                      nil,
	} {
		if _, err := caso.Ejecutar(ctx, "u1", datos); !errors.Is(err, usuario.ErrFotoInvalida) {
			t.Errorf("%s: esperaba ErrFotoInvalida, obtuve %v", nombre, err)
		}
	}
	grande := append(pngDe(10, 10), make([]byte, usuario.FotoMaxBytes)...)
	if _, err := caso.Ejecutar(ctx, "u1", grande); !errors.Is(err, usuario.ErrFotoGrande) {
		t.Errorf("tamaño: esperaba ErrFotoGrande, obtuve %v", err)
	}
	if _, err := caso.Ejecutar(ctx, "u1", pngDe(usuario.FotoMaxLado+1, 1)); !errors.Is(err, usuario.ErrFotoInvalida) {
		t.Errorf("dimensiones: esperaba ErrFotoInvalida, obtuve %v", err)
	}
}

func TestQuitarFoto_BorraReferenciaYObjetoYEsIdempotente(t *testing.T) {
	fotos, almacen := testdobles.NewFotoRepoFalso(), testdobles.NewAlmacenFalso()
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	aud := &testdobles.AuditoriaRepoFalso{}
	_, _ = (&perfil.SubirFoto{Fotos: fotos, Almacen: almacen, Auditoria: aud, Reloj: reloj}).Ejecutar(context.Background(), "u1", pngDe(4, 4))

	quitar := &perfil.QuitarFoto{Fotos: fotos, Almacen: almacen, Auditoria: aud, Reloj: reloj}
	if err := quitar.Ejecutar(context.Background(), "u1"); err != nil {
		t.Fatal(err)
	}
	if _, ok := fotos.Fotos["u1"]; ok || len(almacen.Guardados) != 0 {
		t.Fatal("la foto debía quedar borrada de la base y del almacén")
	}
	if err := quitar.Ejecutar(context.Background(), "u1"); err != nil {
		t.Fatalf("quitar sin foto no debe fallar: %v", err)
	}
}

func TestObtenerFoto_SinFotoDevuelveNoEncontrada(t *testing.T) {
	caso := &perfil.ObtenerFoto{Fotos: testdobles.NewFotoRepoFalso(), Almacen: testdobles.NewAlmacenFalso()}
	if _, err := caso.Ejecutar(context.Background(), "u1"); !errors.Is(err, usuario.ErrFotoNoEncontrada) {
		t.Fatalf("esperaba ErrFotoNoEncontrada, obtuve %v", err)
	}
}
