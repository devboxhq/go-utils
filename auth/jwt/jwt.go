package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// Validator is a validator that can be provided to add custom validation
type Validator func(claims jwt.Claims) bool

// Manager is a jwt manager
type Manager struct {
	signingKey      interface{}       // key used to sign token (vary based on signing method)
	headerScheme    string            // authentication header scheme
	signingMethod   jwt.SigningMethod // signing method for the token
	claims          jwt.Claims        // template for the claims implementation
	customValidate  bool
	customValidator Validator
}

// NewJwtManager creates jwt manager to generate and verify token
func NewJwtManager(signingKey interface{}, headerScheme string,
	signingMethod jwt.SigningMethod, claims jwt.Claims) *Manager {
	return &Manager{
		signingKey:    signingKey,
		headerScheme:  headerScheme,
		signingMethod: signingMethod,
		claims:        claims,
	}
}

func (m *Manager) WithCustomValidator(validator Validator) *Manager {
	m.customValidator = validator
	m.customValidate = true
	return m
}

// GetHeaderScheme returns authentication header scheme to retrieve the token from
func (m *Manager) GetHeaderScheme() string {
	return m.headerScheme
}

// Generate returns token string based on claims provided
func (m *Manager) Generate(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(m.signingMethod, claims)
	signedToken, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", fmt.Errorf("jwt token sign error: %v", err)
	}

	return signedToken, nil
}

// Verify verifies the token provided and returns claims based on provided template
func (m *Manager) Verify(token string) (jwt.Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		m.claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != m.signingMethod {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return m.signingKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("token parse error: %v", err)
	}

	if !jwtToken.Valid {
		return nil, fmt.Errorf("token validation failed")
	}

	if m.customValidate {
		if !m.customValidator(jwtToken.Claims) {
			return nil, fmt.Errorf("token custom validation failed")
		}
	}

	return jwtToken.Claims, nil
}
