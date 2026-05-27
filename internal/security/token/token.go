package token

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Claims are the minimum auth facts the API needs from an access token.
type Claims struct {
	TokenID   string `json:"jti"`
	Type      string `json:"typ"`
	TenantID  int64  `json:"tenant_id"`
	UserID    int64  `json:"sub"`
	RoleID    int64  `json:"role_id"`
	RoleCode  string `json:"role_code"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

// Manager signs and verifies HMAC access tokens.
// This keeps the first auth implementation dependency-light while still using
// the standard JWT wire shape.
type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewManager(accessSecret, refreshSecret string) *Manager {
	return &Manager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (m *Manager) NewAccessToken(claims Claims, ttl time.Duration) (string, Claims, error) {
	now := time.Now().UTC()
	tokenID, err := NewOpaqueToken(24)
	if err != nil {
		return "", Claims{}, err
	}

	claims.TokenID = tokenID
	claims.Type = "access"
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(ttl).Unix()

	token, err := signJWT(claims, m.accessSecret)
	if err != nil {
		return "", Claims{}, err
	}
	return token, claims, nil
}

func (m *Manager) VerifyAccessToken(raw string) (Claims, error) {
	claims, err := verifyJWT(raw, m.accessSecret)
	if err != nil {
		return Claims{}, err
	}
	if claims.Type != "access" {
		return Claims{}, fmt.Errorf("invalid token type")
	}
	if time.Now().Unix() >= claims.ExpiresAt {
		return Claims{}, fmt.Errorf("access token expired")
	}
	return claims, nil
}

func (m *Manager) NewRefreshToken() (string, string, error) {
	raw, err := NewOpaqueToken(48)
	if err != nil {
		return "", "", err
	}
	return raw, m.HashRefreshToken(raw), nil
}

func (m *Manager) HashRefreshToken(raw string) string {
	return HashWithSecret(raw, m.refreshSecret)
}

func HashInvitationToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func HashWithSecret(raw string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func NewOpaqueToken(size int) (string, error) {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		return "", fmt.Errorf("generate secure token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func signJWT(claims Claims, secret []byte) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsPart := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsigned := headerPart + "." + claimsPart
	signature := sign(unsigned, secret)
	return unsigned + "." + signature, nil
}

func verifyJWT(raw string, secret []byte) (Claims, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 3 {
		return Claims{}, fmt.Errorf("invalid token format")
	}

	unsigned := parts[0] + "." + parts[1]
	expected := sign(unsigned, secret)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return Claims{}, fmt.Errorf("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("decode token payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, fmt.Errorf("parse token claims: %w", err)
	}
	return claims, nil
}

func sign(value string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
