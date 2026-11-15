package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"escalafon/backend/internal/application/ports"
)

// claimsJWT es la forma exacta de los claims descrita en la sección 6.
type claimsJWT struct {
	Rol   int16    `json:"rol"`
	Perms []string `json:"perms"`
	Ver   int32    `json:"ver"`
	jwt.RegisteredClaims
}

// ServicioJWT implementa ports.TokenService con HS256 (sección 6).
type ServicioJWT struct {
	secreto []byte
	emisor  string
}

func NewServicioJWT(secreto, emisor string) *ServicioJWT {
	return &ServicioJWT{secreto: []byte(secreto), emisor: emisor}
}

func (s *ServicioJWT) Emitir(claims ports.ClaimsSesion, vigencia time.Duration) (string, error) {
	ahora := time.Now()
	c := claimsJWT{
		Rol:   claims.RolID,
		Perms: claims.Permisos,
		Ver:   claims.Version,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.UsuarioID,
			ID:        claims.JTI,
			Issuer:    s.emisor,
			IssuedAt:  jwt.NewNumericDate(ahora),
			ExpiresAt: jwt.NewNumericDate(ahora.Add(vigencia)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return tok.SignedString(s.secreto)
}

func (s *ServicioJWT) Verificar(tokenStr string) (*ports.ClaimsSesion, error) {
	var c claimsJWT
	tok, err := jwt.ParseWithClaims(tokenStr, &c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return s.secreto, nil
	}, jwt.WithIssuer(s.emisor), jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !tok.Valid {
		return nil, errors.New("token inválido")
	}
	var expira time.Time
	if c.ExpiresAt != nil {
		expira = c.ExpiresAt.Time
	}
	return &ports.ClaimsSesion{
		UsuarioID: c.Subject,
		RolID:     c.Rol,
		Permisos:  c.Perms,
		Version:   c.Ver,
		JTI:       c.ID,
		Expira:    expira,
	}, nil
}
