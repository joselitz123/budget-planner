package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

// JWTClient wraps Clerk JWT verification using JWKS
type JWTClient struct {
	jwks         *keyfunc.JWKS
	clerkDomain  string
}

// NewJWTClient creates a new JWT client for token verification using JWKS
func NewJWTClient(clerkDomain string) (*JWTClient, error) {
	if clerkDomain == "" {
		return nil, fmt.Errorf("Clerk domain is required")
	}

	// Construct JWKS URL for Clerk
	// Format: https://<domain>.clerk.accounts.dev/.well-known/jwks.json
	jwksURL := fmt.Sprintf("https://%s.clerk.accounts.dev/.well-known/jwks.json", clerkDomain)

	// Create JWKS with refresh interval of 1 hour
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS: %w", err)
	}

	return &JWTClient{
		jwks:        jwks,
		clerkDomain: clerkDomain,
	}, nil
}

// VerifyToken verifies a JWT token and returns the user ID (sub claim)
// Uses JWKS to verify RS256 signatures from Clerk
func (c *JWTClient) VerifyToken(tokenString string) (string, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and verify token
	token, err := jwt.Parse(tokenString, c.jwks.Keyfunc)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return "", fmt.Errorf("invalid token: token is not valid")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token: could not parse claims")
	}

	// Verify issuer
	issuer, ok := claims["iss"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token: missing iss claim")
	}

	expectedIssuer := fmt.Sprintf("https://%s.clerk.accounts.dev", c.clerkDomain)
	if issuer != expectedIssuer {
		return "", fmt.Errorf("invalid token: invalid issuer, expected %s, got %s", expectedIssuer, issuer)
	}

	// Verify expiration
	if exp, ok := claims["exp"].(float64); ok {
		if float64(time.Now().Unix()) > exp {
			return "", fmt.Errorf("invalid token: token has expired")
		}
	}

	// Extract subject (user ID)
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token: missing sub claim")
	}

	return sub, nil
}

// ExtractTokenFromRequest extracts the JWT token from the Authorization header
func ExtractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	return token, nil
}

// GenerateToken creates a new JWT token for testing purposes
// Note: This is only for testing and should not be used in production
func (c *JWTClient) GenerateToken(userID string) (string, error) {
	// This method is not supported with JWKS-based verification
	// In production, tokens should only come from Clerk
	return "", fmt.Errorf("GenerateToken is not supported with JWKS verification - use Clerk to generate tokens")
}
