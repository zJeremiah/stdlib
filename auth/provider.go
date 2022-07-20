package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.rakops.com/rm/signal-api/stdlib/stats"
	"github.rakops.com/rm/signal-api/stdlib/xhttp"
)

var (
	// ErrInvalidLogin is an error that is returned when a given
	// username/password is invalid.
	ErrInvalidLogin = errors.New("Invalid login")

	// ErrNoJWTKeys is returned when the JWT keys were not able to be obtained.
	ErrNoJWTKeys = errors.New("could not get JWT signing keys")
)

type (
	// Token is a type alias for a raw JWT token.
	Token []byte

	// Provider authenticates a given username/password with
	// an external source.
	Provider interface {
		Login(username, password string) (Token, error)
		RequestSigningKeys(name, token string) (*SigningKeys, error)
	}

	// SigningKeys is a type that matches the response from the RD Auth Service
	SigningKeys struct {
		PublicKey  string `json:"public_key"`
		PrivateKey string `json:"private_key"`
	}

	providerImpl struct {
		hostname string
		client   xhttp.Client
	}

	fakeProvider struct {
		signingKey []byte
	}

	loginBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	loginResp struct {
		Token string `json:"token"`
	}
)

func (t Token) String() string {
	return string(t)
}

// NewProvider returns an implementation of Provider.
// This implementation uses RM Auth as an authentication
// source.
func NewProvider(authHostname string) Provider {
	return &providerImpl{
		hostname: authHostname,
		client:   xhttp.NewDefaultClient(),
	}
}

// NewProviderWithClient returns an implementation of Provider.
// This implementation uses RM Auth as an authentication source.
func NewProviderWithClient(authHostname string, client xhttp.Client) Provider {
	return &providerImpl{
		hostname: authHostname,
		client:   client,
	}
}

// NewProviderWithStats returns an implementation of Provider.
// This implementation uses RM Auth as an authentication
// source.
func NewProviderWithStats(authHostname string, stats stats.Client) Provider {
	return &providerImpl{
		hostname: authHostname,
		client:   xhttp.NewDefaultStatsClient(stats),
	}
}

// NewFakeProvider returns an implementation of Provider.
// This implementation always returns true and uses the
// given signingKey to sign the JWT token.
func NewFakeProvider(signingKey []byte) Provider {
	return &fakeProvider{
		signingKey: signingKey,
	}
}

func (f *fakeProvider) Login(username, password string) (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix()})
	signedToken, err := token.SignedString(f.signingKey)
	if err != nil {
		return nil, err
	}
	return []byte(signedToken), nil
}

func (f *fakeProvider) RequestSigningKeys(name, token string) (*SigningKeys, error) {
	return new(SigningKeys), nil
}

func (a *providerImpl) Login(username, password string) (Token, error) {
	login := &loginBody{
		Email:    username,
		Password: password,
	}

	b, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", a.makeURL("v1/login"), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrInvalidLogin
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respBody := new(loginResp)

	if err := json.Unmarshal(b, respBody); err != nil {
		return nil, err
	}

	return []byte(respBody.Token), nil
}

func (a *providerImpl) makeURL(path string) string {
	return fmt.Sprintf("%s/%s", a.hostname, path)
}

// RequestSigningKeys requests the signing keys used on the RD Auth Service to sign JWT keys.
// Once these keys have been obtained, one will be able to verify JWT tokens issued by
// the RD Auth Service.
func (a *providerImpl) RequestSigningKeys(name, token string) (*SigningKeys, error) {
	location := fmt.Sprintf("%s/v1/auth", strings.TrimSuffix(a.hostname, "/"))
	body := []byte(fmt.Sprintf(`{"service_name": "%s", "service_token": "%s"}`, name, token))
	req, err := http.NewRequest("POST", location, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrNoJWTKeys
	}

	keys := new(SigningKeys)
	if err := json.Unmarshal(b, keys); err != nil {
		return nil, err
	}

	return keys, nil
}
