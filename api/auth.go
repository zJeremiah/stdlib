package api

import (
	"encoding/base64"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.rakops.com/rm/signal-api/stdlib/auth"
	"github.rakops.com/rm/signal-api/stdlib/xhttp"
)

type (
	AuthOptions struct {
		AuthServer       string
		AppName          string
		AppToken         string
		Devmode          bool
		AdditionalFields []string
	}

	RMAuthJWTConfig struct {
		middleware.JWTConfig

		AdditionalFields []string
	}
)

// DevSigningKey is a hardcoded key that's used to sign
// JWT keys in development mode.
var DevSigningKey = []byte("signingkey")

var DefaultRMAuthJWTConfig = RMAuthJWTConfig{
	JWTConfig: middleware.DefaultJWTConfig,
}

// SetupAuth contacts the given auth servers and obtains its signing keys. Then it uses these keys to validate incoming JWT values.
func SetupAuth(opts *AuthOptions, client xhttp.Client, skipper middleware.Skipper) (*auth.SigningKeys, echo.MiddlewareFunc, error) {
	var signingKeys *auth.SigningKeys

	if !opts.Devmode {
		sk, err := auth.NewProviderWithClient(opts.AuthServer, client).RequestSigningKeys(opts.AppName, opts.AppToken)
		if err != nil {
			return nil, nil, err
		}

		signingKeys = sk
	}

	jwtConf := RMAuthJWTConfig{
		AdditionalFields: opts.AdditionalFields,
		JWTConfig:        middleware.DefaultJWTConfig,
	}

	jwtConf.JWTConfig.Skipper = skipper

	if signingKeys != nil {
		b, err := base64.StdEncoding.DecodeString(signingKeys.PrivateKey)
		if err != nil {
			return nil, nil, err
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
		if err != nil {
			return nil, nil, err
		}

		jwtConf.JWTConfig.SigningMethod = "RS256"
		jwtConf.JWTConfig.SigningKey = &privateKey.PublicKey
	} else {
		jwtConf.JWTConfig.SigningMethod = "HS256"
		jwtConf.JWTConfig.SigningKey = DevSigningKey
	}

	return signingKeys, RMAuthJWT(jwtConf), nil
}

func RMAuthJWT(config RMAuthJWTConfig) echo.MiddlewareFunc {
	jwtMiddleware := middleware.JWTWithConfig(config.JWTConfig)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := jwtMiddleware(fakeHandler)(c); err != nil {
				return err
			}

			if x := c.Get("user"); x != nil {
				if user, ok := x.(*jwt.Token); ok {
					if claims, ok := user.Claims.(jwt.MapClaims); ok {
						if username, ok := claims["username"]; ok {
							c.Set("username", username)
						} else {
							return echo.NewHTTPError(http.StatusUnauthorized)
						}

						if roles, ok := claims["roles"]; ok {
							c.Set("roles", roles)
						} else {
							return echo.NewHTTPError(http.StatusUnauthorized)
						}

						if id, ok := claims["id"]; ok {
							c.Set("id", id)
						} else {
							return echo.NewHTTPError(http.StatusUnauthorized)
						}

						for _, field := range config.AdditionalFields {
							if value, ok := claims[field]; ok {
								c.Set(field, value)
							}
						}
					}
				}
			}

			return next(c)
		}
	}
}

func fakeHandler(c echo.Context) error {
	return nil
}
