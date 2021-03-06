package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"time"

	"github.com/neuronlabs/neuron/store"
)

// Tokener is the interface used for the authorization with the token.
type Tokener interface {
	// InspectToken extracts claims from the token.
	InspectToken(ctx context.Context, token string) (claims Claims, err error)
	// Token creates the token for provided options.
	Token(ctx context.Context, account Account, options ...TokenOption) (Token, error)
	// RevokeToken revokes provided 'token'
	RevokeToken(ctx context.Context, token string) error
}

// Token is the authorization token structure.
type Token struct {
	// AccessToken is the string access token.
	AccessToken string
	// RefreshToken defines the token.
	RefreshToken string
	// ExpiresIn defines the expiration time for given access token.
	ExpiresIn int
	// TokenType defines the token type.
	TokenType string
}

// TokenOptions is the options used to create the token.
type TokenOptions struct {
	// ExpirationTime is the expiration time of the token.
	ExpirationTime time.Duration
	// RefreshExpirationTime is the expiration time for refresh token
	RefreshExpirationTime time.Duration
	// RefreshToken is the optional refresh token used on token creation, when the refresh token is still valid (optional).
	RefreshToken string

	// Optional settings.
	//
	// Scope contains space separated authorization scopes that the token is available for (optional).
	Scope string
	// Audience is the audience of the token.
	Audience string
	// Issuer is the token issuer name.
	Issuer string
	// NotBefore is an option that sets the token to be valid not before provided time.
	NotBefore time.Time
}

// TokenOption is the token options changer function.
type TokenOption func(o *TokenOptions)

// TokenExpirationTime sets the expiration time for the token.
func TokenExpirationTime(d time.Duration) TokenOption {
	return func(o *TokenOptions) {
		o.ExpirationTime = d
	}
}

// TokenRefreshExpirationTime sets the expiration time for the token.
func TokenRefreshExpirationTime(d time.Duration) TokenOption {
	return func(o *TokenOptions) {
		o.RefreshExpirationTime = d
	}
}

// TokenRefreshToken sets the refresh token for the token creation.
func TokenRefreshToken(refreshToken string) TokenOption {
	return func(o *TokenOptions) {
		o.RefreshToken = refreshToken
	}
}

// TokenScope sets the space separated scopes where the token should have an access.
func TokenScope(scope string) TokenOption {
	return func(o *TokenOptions) {
		o.Scope = scope
	}
}

// TokenWithAudience sets the token audience.
func TokenWithAudience(audience string) TokenOption {
	return func(o *TokenOptions) {
		o.Audience = audience
	}
}

// TokenWithIssuer is the token option that sets up the issuer.
func TokenWithIssuer(issuer string) TokenOption {
	return func(o *TokenOptions) {
		o.Issuer = issuer
	}
}

// TokenWithNotBefore is the token option that sets up the not before option.
func TokenWithNotBefore(notBefore time.Time) TokenOption {
	return func(o *TokenOptions) {
		o.NotBefore = notBefore
	}
}

// AccessClaims is an interface used for the access token claims. It should store the whole user account.
type AccessClaims interface {
	// GetAccount gets the account stored in given token.
	GetAccount() Account
	Claims
}

// Claims is an interface used for the tokens.
type Claims interface {
	// Subject should contain account id string value.
	Subject() string
	// ExpiresIn should define when (in seconds) the claims will expire.
	ExpiresIn() int64
	// Valid validates the claims.
	Valid() error
}

// NotBeforer is an interface that allows to get Token's NotBefore (nbf) value.
type NotBeforer interface {
	NotBefore() int64
}

// Scoper is an interface that allows to get Token's authorization scope value. This should return all of the scopes
// for which the token is authorized, space separated.
type Scoper interface {
	Scope() string
}

// Audiencer is an interface that allows to get token's optional audience value.
type Audiencer interface {
	Audience() string
}

// Issuer is an interface that allows to get the token issuer.
type Issuer interface {
	Issuer() string
}

// SigningMethod is an interface used for signing and verify the string.
// This interface is equal to the Signing method of github.com/dgrijalva/jwt-go.
type SigningMethod interface {
	Verify(signingString, signature string, key interface{}) error
	Sign(signingString string, key interface{}) (string, error)
	Alg() string
}

// TokenerOptions are the options that defines the settings for the Tokener.
type TokenerOptions struct {
	// Model is the account model used by the tokener.
	Model Account
	// Store is a store used for some authenticator implementations.
	Store store.Store
	// Secret is the authorization secret.
	Secret []byte
	// RsaPrivateKey is used for encoding the token using RSA methods.
	RsaPrivateKey *rsa.PrivateKey
	// EcdsaPrivateKey is used for encoding the token using ECDSA methods.
	EcdsaPrivateKey *ecdsa.PrivateKey
	// TokenExpiration is the default token expiration time.
	TokenExpiration time.Duration
	// RefreshTokenExpiration is the default refresh token expiration time,.
	RefreshTokenExpiration time.Duration
	// SigningMethod is the token signing method.
	SigningMethod SigningMethod
	// TimeFunc sets the time function for given tokener.
	TimeFunc func() time.Time
}

// TokenerOption is a function that sets the TokenerOptions.
type TokenerOption func(o *TokenerOptions)

// TokenerAccount sets the account for the tokener.
func TokenerAccount(model Account) TokenerOption {
	return func(o *TokenerOptions) {
		o.Model = model
	}
}

// TokenerSecret is an option that sets Secret in the auth options.
func TokenerSecret(secret []byte) TokenerOption {
	return func(o *TokenerOptions) {
		o.Secret = secret
	}
}

// TokenerRsaPrivateKey is an option that sets RsaPrivateKey in the auth options.
func TokenerRsaPrivateKey(key *rsa.PrivateKey) TokenerOption {
	return func(o *TokenerOptions) {
		o.RsaPrivateKey = key
	}
}

// TokenerEcdsaPrivateKey is an option that sets EcdsaPrivateKey in the auth options.
func TokenerEcdsaPrivateKey(key *ecdsa.PrivateKey) TokenerOption {
	return func(o *TokenerOptions) {
		o.EcdsaPrivateKey = key
	}
}

// TokenerTokenExpiration is an option that sets TokenExpiration in the auth options.
func TokenerTokenExpiration(op time.Duration) TokenerOption {
	return func(o *TokenerOptions) {
		o.TokenExpiration = op
	}
}

// TokenerRefreshTokenExpiration is an option that sets RefreshTokenExpiration in the auth options.
func TokenerRefreshTokenExpiration(op time.Duration) TokenerOption {
	return func(o *TokenerOptions) {
		o.RefreshTokenExpiration = op
	}
}

// TokenerSigningMethod is an option that sets SigningMethod in the auth options.
func TokenerSigningMethod(op SigningMethod) TokenerOption {
	return func(o *TokenerOptions) {
		o.SigningMethod = op
	}
}

// TokenerStore sets the store for the tokener.
func TokenerStore(s store.Store) TokenerOption {
	return func(o *TokenerOptions) {
		o.Store = s
	}
}

// TokenerTimeFunc sets the default time function for the tokener.
func TokenerTimeFunc(tf func() time.Time) TokenerOption {
	return func(o *TokenerOptions) {
		o.TimeFunc = tf
	}
}
