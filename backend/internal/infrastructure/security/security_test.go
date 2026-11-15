package security

import (
	"testing"
	"time"

	"escalafon/backend/internal/application/ports"
)

func TestBcryptHasher_HashYComparar(t *testing.T) {
	h := NewBcryptHasher()
	hash, err := h.Hash("ClaveSegura1")
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if !h.Comparar(hash, "ClaveSegura1") {
		t.Fatalf("debería aceptar la contraseña correcta")
	}
	if h.Comparar(hash, "Incorrecta1") {
		t.Fatalf("no debería aceptar una contraseña incorrecta")
	}
}

func TestGeneradorAleatorio_GeneraYHashCoinciden(t *testing.T) {
	g := NewGeneradorAleatorio()
	plano, hash, err := g.Generar()
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if plano == "" || hash == "" {
		t.Fatalf("plano y hash no deberían estar vacíos")
	}
	if g.Hash(plano) != hash {
		t.Fatalf("Hash(plano) debería coincidir con el hash devuelto por Generar")
	}
	plano2, _, _ := g.Generar()
	if plano == plano2 {
		t.Fatalf("dos llamadas a Generar no deberían producir el mismo token")
	}
}

func TestServicioJWT_EmitirYVerificar(t *testing.T) {
	s := NewServicioJWT("un-secreto-de-al-menos-32-bytes!!", "escalafon-uagrm")
	claims := ports.ClaimsSesion{UsuarioID: "u1", RolID: 1, Permisos: []string{"USUARIO_GESTIONAR"}, Version: 3, JTI: "jti-1"}

	tok, err := s.Emitir(claims, time.Hour)
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	verificados, err := s.Verificar(tok)
	if err != nil {
		t.Fatalf("no esperaba error al verificar: %v", err)
	}
	if verificados.UsuarioID != claims.UsuarioID || verificados.RolID != claims.RolID ||
		verificados.Version != claims.Version || verificados.JTI != claims.JTI {
		t.Fatalf("claims no coinciden: %+v", verificados)
	}

	if _, err := s.Verificar(tok + "manipulado"); err == nil {
		t.Fatalf("un token manipulado no debería ser válido")
	}
}

func TestServicioJWT_RechazaTokenExpirado(t *testing.T) {
	s := NewServicioJWT("un-secreto-de-al-menos-32-bytes!!", "escalafon-uagrm")
	tok, _ := s.Emitir(ports.ClaimsSesion{UsuarioID: "u1"}, -time.Minute)
	if _, err := s.Verificar(tok); err == nil {
		t.Fatalf("un token expirado no debería ser válido")
	}
}
