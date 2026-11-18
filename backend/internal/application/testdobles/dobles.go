// Package testdobles reúne dobles en memoria de los puertos y repositorios, usados por las
// pruebas unitarias de los casos de uso (sección 9, Fase 2: "dobles en memoria").
package testdobles

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/notificacion"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/domain/usuario"
)

// --- Clock -------------------------------------------------------------

type RelojFijo struct{ Momento time.Time }

func (r *RelojFijo) Ahora() time.Time { return r.Momento }

// --- PasswordHasher ------------------------------------------------------

// HasherFalso evita bcrypt en las pruebas: el "hash" es solo un prefijo reconocible.
type HasherFalso struct{}

func (HasherFalso) Hash(plano string) (string, error) { return "hash:" + plano, nil }
func (HasherFalso) Comparar(hash, plano string) bool  { return hash == "hash:"+plano }

// --- TokenService --------------------------------------------------------

type TokenServiceFalso struct {
	mu       sync.Mutex
	emitidos map[string]ports.ClaimsSesion
}

func NewTokenServiceFalso() *TokenServiceFalso {
	return &TokenServiceFalso{emitidos: map[string]ports.ClaimsSesion{}}
}

func (t *TokenServiceFalso) Emitir(claims ports.ClaimsSesion, _ time.Duration) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	tok := fmt.Sprintf("tok-%d", len(t.emitidos)+1)
	t.emitidos[tok] = claims
	return tok, nil
}

func (t *TokenServiceFalso) Verificar(tok string) (*ports.ClaimsSesion, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	c, ok := t.emitidos[tok]
	if !ok {
		return nil, fmt.Errorf("token no reconocido")
	}
	return &c, nil
}

// --- GeneradorToken --------------------------------------------------------

type GeneradorFalso struct {
	mu       sync.Mutex
	contador int
}

func (g *GeneradorFalso) Generar() (string, string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.contador++
	plano := fmt.Sprintf("plano-%d", g.contador)
	return plano, g.Hash(plano), nil
}

func (*GeneradorFalso) Hash(plano string) string { return "hash:" + plano }

// --- Mailer --------------------------------------------------------------

type CorreoEnviado struct {
	Tipo, Destino, Nombres, Enlace string
}

type MailerFalso struct {
	mu       sync.Mutex
	Enviados []CorreoEnviado
	// ErrCambioEstado simula una caída del SMTP en los avisos de CU-14.
	ErrCambioEstado error
}

func (m *MailerFalso) EnviarActivacion(destino, nombres, enlace string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enviados = append(m.Enviados, CorreoEnviado{"ACTIVACION", destino, nombres, enlace})
	return nil
}

func (m *MailerFalso) EnviarRecuperacion(destino, nombres, enlace string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enviados = append(m.Enviados, CorreoEnviado{"RECUPERACION", destino, nombres, enlace})
	return nil
}

// EnviarCambioEstado registra el aviso; si ErrCambioEstado no es nil lo devuelve sin registrar.
func (m *MailerFalso) EnviarCambioEstado(destino, nombres, codigo, estado, enlace string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ErrCambioEstado != nil {
		return m.ErrCambioEstado
	}
	m.Enviados = append(m.Enviados, CorreoEnviado{"CAMBIO_ESTADO:" + estado, destino, nombres, enlace})
	return nil
}

// --- SessionCache ----------------------------------------------------------

type SessionCacheFalso struct {
	mu        sync.Mutex
	sesiones  map[string]struct {
		version int32
		activo  bool
	}
	revocados map[string]bool
}

func NewSessionCacheFalso() *SessionCacheFalso {
	return &SessionCacheFalso{
		sesiones: map[string]struct {
			version int32
			activo  bool
		}{},
		revocados: map[string]bool{},
	}
}

func (s *SessionCacheFalso) ObtenerSesion(_ context.Context, usuarioID string) (int32, bool, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.sesiones[usuarioID]
	return v.version, v.activo, ok, nil
}

func (s *SessionCacheFalso) GuardarSesion(_ context.Context, usuarioID string, version int32, activo bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sesiones[usuarioID] = struct {
		version int32
		activo  bool
	}{version, activo}
	return nil
}

func (s *SessionCacheFalso) InvalidarSesion(_ context.Context, usuarioID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sesiones, usuarioID)
	return nil
}

func (s *SessionCacheFalso) RevocarJTI(_ context.Context, jti string, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revocados[jti] = true
	return nil
}

func (s *SessionCacheFalso) EstaRevocado(_ context.Context, jti string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.revocados[jti], nil
}

// --- TxManager ---------------------------------------------------------

// TxManagerFalso ejecuta fn directamente: en memoria no hay nada que transaccionar.
type TxManagerFalso struct{}

func (TxManagerFalso) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

// --- usuario.Repository ------------------------------------------------

type UsuarioRepoFalso struct {
	mu        sync.Mutex
	porID     map[string]*usuario.Usuario
	siguiente int
}

func NewUsuarioRepoFalso() *UsuarioRepoFalso {
	return &UsuarioRepoFalso{porID: map[string]*usuario.Usuario{}}
}

func (r *UsuarioRepoFalso) BuscarPorID(_ context.Context, id string) (*usuario.Usuario, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.porID[id]
	if !ok {
		return nil, usuario.ErrNoEncontrado
	}
	copia := *u
	return &copia, nil
}

func (r *UsuarioRepoFalso) BuscarPorCorreo(_ context.Context, correo string) (*usuario.Usuario, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.porID {
		if u.Correo == correo {
			copia := *u
			return &copia, nil
		}
	}
	return nil, usuario.ErrNoEncontrado
}

func (r *UsuarioRepoFalso) Listar(_ context.Context, _ usuario.ListarFiltro) ([]usuario.Usuario, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []usuario.Usuario
	for _, u := range r.porID {
		out = append(out, *u)
	}
	return out, int64(len(out)), nil
}

func (r *UsuarioRepoFalso) correoDuplicado(correo, exceptoID string) bool {
	if correo == "" {
		return false
	}
	for id, u := range r.porID {
		if id != exceptoID && strings.EqualFold(u.Correo, correo) {
			return true
		}
	}
	return false
}

func (r *UsuarioRepoFalso) Crear(_ context.Context, u *usuario.Usuario) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.correoDuplicado(u.Correo, "") {
		return usuario.ErrCorreoDuplicado
	}
	if u.ID == "" {
		r.siguiente++
		u.ID = fmt.Sprintf("usuario-%d", r.siguiente)
	}
	copia := *u
	r.porID[u.ID] = &copia
	return nil
}

func (r *UsuarioRepoFalso) Actualizar(_ context.Context, u *usuario.Usuario) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.porID[u.ID]; !ok {
		return usuario.ErrNoEncontrado
	}
	if r.correoDuplicado(u.Correo, u.ID) {
		return usuario.ErrCorreoDuplicado
	}
	copia := *u
	r.porID[u.ID] = &copia
	return nil
}

func (r *UsuarioRepoFalso) IncrementarVersionSesionPorRol(_ context.Context, rolID int16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.porID {
		if u.RolID == rolID {
			u.VersionSesion++
		}
	}
	return nil
}

// Sembrar agrega un usuario directamente, para preparar el escenario de una prueba.
func (r *UsuarioRepoFalso) Sembrar(u usuario.Usuario) *usuario.Usuario {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u.ID == "" {
		r.siguiente++
		u.ID = fmt.Sprintf("usuario-%d", r.siguiente)
	}
	r.porID[u.ID] = &u
	return &u
}

// --- rol.Repository ------------------------------------------------------

type RolRepoFalso struct {
	mu             sync.Mutex
	roles          map[int16]*rol.Rol
	permisosPorRol map[int16][]string // rolID -> códigos asignados
	catalogo       map[string]int16   // código -> id, catálogo global de permisos (como la tabla permiso)
	usuariosPorRol map[int16][]string // rolID -> IDs de usuario, para simular ContarUsuarios/Eliminar
}

func NewRolRepoFalso() *RolRepoFalso {
	return &RolRepoFalso{
		roles: map[int16]*rol.Rol{}, permisosPorRol: map[int16][]string{},
		catalogo: map[string]int16{}, usuariosPorRol: map[int16][]string{},
	}
}

// idParaCodigo asigna (o reutiliza) el id de un código de permiso en el catálogo. Se llama con
// el lock ya tomado.
func (r *RolRepoFalso) idParaCodigo(codigo string) int16 {
	if id, ok := r.catalogo[codigo]; ok {
		return id
	}
	id := int16(len(r.catalogo) + 1)
	r.catalogo[codigo] = id
	return id
}

// AsignarUsuarioARol prepara el escenario para ContarUsuarios/IDsUsuariosDelRol/Eliminar.
func (r *RolRepoFalso) AsignarUsuarioARol(rolID int16, usuarioID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.usuariosPorRol[rolID] = append(r.usuariosPorRol[rolID], usuarioID)
}

// Sembrar registra el rol y, de paso, da de alta en el catálogo cualquier código de permiso
// nuevo que traiga (como si ya existiera en la tabla permiso).
func (r *RolRepoFalso) Sembrar(rl rol.Rol, permisos []string) *rol.Rol {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, cod := range permisos {
		r.idParaCodigo(cod)
	}
	r.roles[rl.ID] = &rl
	r.permisosPorRol[rl.ID] = permisos
	return &rl
}

func (r *RolRepoFalso) conPermisos(rl *rol.Rol) rol.Rol {
	copia := *rl
	for _, cod := range r.permisosPorRol[rl.ID] {
		copia.Permisos = append(copia.Permisos, rol.Permiso{ID: r.catalogo[cod], Codigo: cod})
	}
	return copia
}

func (r *RolRepoFalso) Listar(context.Context) ([]rol.ConConteo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []rol.ConConteo
	for _, rl := range r.roles {
		out = append(out, rol.ConConteo{Rol: r.conPermisos(rl), UsuariosCount: int64(len(r.usuariosPorRol[rl.ID]))})
	}
	return out, nil
}

func (r *RolRepoFalso) BuscarPorID(_ context.Context, id int16) (*rol.Rol, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rl, ok := r.roles[id]
	if !ok {
		return nil, rol.ErrNoEncontrado
	}
	copia := r.conPermisos(rl)
	return &copia, nil
}

func (r *RolRepoFalso) BuscarPorNombre(_ context.Context, nombre string) (*rol.Rol, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rl := range r.roles {
		if rl.Nombre == nombre {
			copia := r.conPermisos(rl)
			return &copia, nil
		}
	}
	return nil, rol.ErrNoEncontrado
}

func (r *RolRepoFalso) Guardar(_ context.Context, rl *rol.Rol, permisoIDs []int16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, otro := range r.roles {
		if otro.Nombre == rl.Nombre && id != rl.ID {
			return rol.ErrDuplicado
		}
	}
	if rl.ID == 0 {
		rl.ID = int16(len(r.roles) + 1)
	}
	copia := *rl
	r.roles[rl.ID] = &copia

	idACodigo := map[int16]string{}
	for cod, id := range r.catalogo {
		idACodigo[id] = cod
	}
	codigos := make([]string, 0, len(permisoIDs))
	for _, id := range permisoIDs {
		if cod, ok := idACodigo[id]; ok {
			codigos = append(codigos, cod)
		}
	}
	r.permisosPorRol[rl.ID] = codigos
	return nil
}

func (r *RolRepoFalso) Eliminar(_ context.Context, id int16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n := len(r.usuariosPorRol[id]); n > 0 {
		return &rol.ErrConUsuarios{Cantidad: int64(n)}
	}
	delete(r.roles, id)
	return nil
}

func (r *RolRepoFalso) ContarUsuarios(_ context.Context, id int16) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return int64(len(r.usuariosPorRol[id])), nil
}

func (r *RolRepoFalso) IDsUsuariosDelRol(_ context.Context, id int16) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.usuariosPorRol[id], nil
}

func (r *RolRepoFalso) ListarPermisos(context.Context) ([]rol.Permiso, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]rol.Permiso, 0, len(r.catalogo))
	for cod, id := range r.catalogo {
		out = append(out, rol.Permiso{ID: id, Codigo: cod})
	}
	return out, nil
}

func (r *RolRepoFalso) PermisosExisten(_ context.Context, ids []int16) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	validos := map[int16]bool{}
	for _, id := range r.catalogo {
		validos[id] = true
	}
	for _, id := range ids {
		if !validos[id] {
			return false, nil
		}
	}
	return true, nil
}

func (r *RolRepoFalso) PermisosPorRol(_ context.Context, rolID int16) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.permisosPorRol[rolID], nil
}

// --- token.Repository ------------------------------------------------------

type TokenRepoFalso struct {
	mu       sync.Mutex
	porHash  map[string]*token.TokenRecuperacion
	siguiente int
}

func NewTokenRepoFalso() *TokenRepoFalso {
	return &TokenRepoFalso{porHash: map[string]*token.TokenRecuperacion{}}
}

func (t *TokenRepoFalso) Crear(_ context.Context, tk *token.TokenRecuperacion) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if tk.ID == "" {
		t.siguiente++
		tk.ID = fmt.Sprintf("token-%d", t.siguiente)
	}
	copia := *tk
	t.porHash[tk.TokenHash] = &copia
	return nil
}

func (t *TokenRepoFalso) BuscarPorHash(_ context.Context, hash string) (*token.TokenRecuperacion, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	tk, ok := t.porHash[hash]
	if !ok {
		return nil, token.ErrInvalido
	}
	copia := *tk
	return &copia, nil
}

func (t *TokenRepoFalso) MarcarUsado(_ context.Context, id string, cuando time.Time) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, tk := range t.porHash {
		if tk.ID == id {
			tk.UsadoEn = &cuando
			return nil
		}
	}
	return token.ErrInvalido
}

func (t *TokenRepoFalso) InvalidarPendientes(_ context.Context, usuarioID string, proposito token.Proposito, cuando time.Time) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, tk := range t.porHash {
		if tk.UsuarioID == usuarioID && tk.Proposito == proposito && tk.UsadoEn == nil {
			tk.UsadoEn = &cuando
		}
	}
	return nil
}

func (t *TokenRepoFalso) ExisteTokenUsado(_ context.Context, usuarioID string, proposito token.Proposito) (bool, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, tk := range t.porHash {
		if tk.UsuarioID == usuarioID && tk.Proposito == proposito && tk.UsadoEn != nil {
			return true, nil
		}
	}
	return false, nil
}

func (t *TokenRepoFalso) ContarDesde(_ context.Context, usuarioID string, proposito token.Proposito, desde time.Time) (int64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	var n int64
	for _, tk := range t.porHash {
		if tk.UsuarioID == usuarioID && tk.Proposito == proposito && !tk.CreadoEn.Before(desde) {
			n++
		}
	}
	return n, nil
}

// --- auditoria.Repository --------------------------------------------------

// AuditoriaRepoFalso implementa el registro (auditoria.Repository) y la consulta
// (auditoria.Consulta). Autores simula el JOIN con usuario: id -> datos del autor.
type AuditoriaRepoFalso struct {
	mu      sync.Mutex
	Eventos []auditoria.Evento
	Autores map[string]auditoria.Autor
}

func (a *AuditoriaRepoFalso) Registrar(_ context.Context, e *auditoria.Evento) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if e.ID == 0 {
		e.ID = int64(len(a.Eventos) + 1)
	}
	a.Eventos = append(a.Eventos, *e)
	return nil
}

func (a *AuditoriaRepoFalso) Listar(_ context.Context, f auditoria.Filtro) ([]auditoria.Registro, int64, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var coinciden []auditoria.Registro
	for _, e := range a.Eventos {
		switch {
		case f.UsuarioID != "" && (e.UsuarioID == nil || *e.UsuarioID != f.UsuarioID),
			f.Accion != "" && e.Accion != f.Accion,
			f.Desde != nil && e.Fecha.Before(*f.Desde),
			f.Hasta != nil && !e.Fecha.Before(*f.Hasta):
			continue
		}
		reg := auditoria.Registro{Evento: e}
		if e.UsuarioID != nil {
			if autor, ok := a.Autores[*e.UsuarioID]; ok {
				reg.Autor = &autor
			}
		}
		coinciden = append(coinciden, reg)
	}
	sort.Slice(coinciden, func(i, j int) bool {
		if !coinciden[i].Fecha.Equal(coinciden[j].Fecha) {
			return coinciden[i].Fecha.After(coinciden[j].Fecha)
		}
		return coinciden[i].ID > coinciden[j].ID
	})
	total := int64(len(coinciden))
	desde := min((f.Pagina-1)*f.Tamano, len(coinciden))
	hasta := min(desde+f.Tamano, len(coinciden))
	return coinciden[desde:hasta], total, nil
}

func (a *AuditoriaRepoFalso) ListarAutores(context.Context) ([]auditoria.Autor, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	vistos := map[string]bool{}
	var out []auditoria.Autor
	for _, e := range a.Eventos {
		if e.UsuarioID == nil || vistos[*e.UsuarioID] {
			continue
		}
		if autor, ok := a.Autores[*e.UsuarioID]; ok {
			vistos[*e.UsuarioID] = true
			out = append(out, autor)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Apellidos+out[i].Nombres < out[j].Apellidos+out[j].Nombres })
	return out, nil
}

// --- notificacion.Repository ---------------------------------------------

type NotificacionRepoFalso struct{ Envios []*notificacion.Envio }

func (r *NotificacionRepoFalso) Crear(_ context.Context, e *notificacion.Envio) error {
	e.ID = fmt.Sprintf("n-%d", len(r.Envios)+1)
	c := *e
	r.Envios = append(r.Envios, &c)
	return nil
}

func (r *NotificacionRepoFalso) Actualizar(_ context.Context, e *notificacion.Envio) error {
	for _, x := range r.Envios {
		if x.ID == e.ID {
			*x = *e
		}
	}
	return nil
}

func (r *NotificacionRepoFalso) ListarPendientes(_ context.Context, limite int) ([]notificacion.Envio, error) {
	var out []notificacion.Envio
	for _, x := range r.Envios {
		if x.Estado == notificacion.Pendiente && len(out) < limite {
			out = append(out, *x)
		}
	}
	return out, nil
}

func (r *NotificacionRepoFalso) Listar(_ context.Context, f notificacion.Filtro, _, _ int) ([]notificacion.Envio, int64, error) {
	var out []notificacion.Envio
	for _, x := range r.Envios {
		if f.Estado == "" || x.Estado == f.Estado {
			out = append(out, *x)
		}
	}
	return out, int64(len(out)), nil
}

// --- usuario.FotoRepository ------------------------------------------------------------

type FotoRepoFalso struct{ Fotos map[string]*usuario.Foto }

func NewFotoRepoFalso() *FotoRepoFalso { return &FotoRepoFalso{Fotos: map[string]*usuario.Foto{}} }

func (r *FotoRepoFalso) Obtener(_ context.Context, usuarioID string) (*usuario.Foto, error) {
	f, ok := r.Fotos[usuarioID]
	if !ok {
		return nil, usuario.ErrFotoNoEncontrada
	}
	return f, nil
}

func (r *FotoRepoFalso) Guardar(_ context.Context, f *usuario.Foto) error {
	r.Fotos[f.UsuarioID] = f
	return nil
}

func (r *FotoRepoFalso) Eliminar(_ context.Context, usuarioID string) error {
	delete(r.Fotos, usuarioID)
	return nil
}
