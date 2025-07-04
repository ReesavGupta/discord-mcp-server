package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type AuthManager struct {
	jwtSecret []byte
	apiKeys   map[string]bool
	logger    *logrus.Logger
	auditor   *AuditLogger
}

type Claims struct {
	UserID      string   `json:"user_id"`
	Permissions []string `json:"permissions"`
	BotID       string   `json:"bot_id"`
	jwt.RegisteredClaims
}

type AuditLogger struct {
	logger *logrus.Logger
}

func NewAuthManager(jwtSecret string, apiKeys []string, logger *logrus.Logger, enableAudit bool, auditPath string) (*AuthManager, error) {
	keyMap := make(map[string]bool)
	for _, key := range apiKeys {
		keyMap[key] = true
	}

	var auditor *AuditLogger
	if enableAudit {
		auditor = &AuditLogger{logger: logger}
	}

	return &AuthManager{
		jwtSecret: []byte(jwtSecret),
		apiKeys:   keyMap,
		logger:    logger,
		auditor:   auditor,
	}, nil
}

func (am *AuthManager) ValidateAPIKey(apiKey string) bool {
	valid := am.apiKeys[apiKey]
	if am.auditor != nil {
		am.auditor.LogAuth("api_key", apiKey, valid)
	}
	return valid
}

func (am *AuthManager) GenerateToken(userID string, permissions []string, botID string) (string, error) {
	claims := Claims{
		UserID:      userID,
		Permissions: permissions,
		BotID:       botID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(am.jwtSecret)
}

func (am *AuthManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return am.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if am.auditor != nil {
			am.auditor.LogAuth("jwt_token", claims.UserID, true)
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (am *AuthManager) GenerateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (al *AuditLogger) LogAuth(method, identifier string, success bool) {
	al.logger.WithFields(logrus.Fields{
		"method":     method,
		"identifier": identifier,
		"success":    success,
		"timestamp":  time.Now(),
	}).Info("Authentication attempt")
}

func (al *AuditLogger) LogOperation(operation, userID string, data interface{}) {
	al.logger.WithFields(logrus.Fields{
		"operation": operation,
		"user_id":   userID,
		"data":      data,
		"timestamp": time.Now(),
	}).Info("Operation performed")
}
