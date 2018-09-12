package config

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	JWT struct {
		Secret     string
		Expiry     int64
		DebugToken bool
	}
)

var jwt *JWT

func (c *JWT) Validate() error {
	if c == nil {
		return nil
	}
	if c.Secret == "" {
		return errors.New("JWT Secret not set for AUTH")
	}
	return nil
}

func (*JWT) Init(prefix ...string) *JWT {
	if jwt != nil {
		return jwt
	}

	jwt := new(JWT)
	flag.StringVar(&jwt.Secret, "auth-jwt-secret", "", "JWT Secret")
	flag.Int64Var(&jwt.Expiry, "auth-jwt-expiry", 3600, "JWT Expiration in minutes")
	flag.BoolVar(&jwt.DebugToken, "auth-jwt-debug", false, "Generate debug JWT key")
	return jwt
}