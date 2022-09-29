package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zitadel/oidc/pkg/oidc"
	"net/http"
	"os"
	"strings"
)

var key *rsa.PublicKey

func initBackend() {
	rawKey := os.Getenv("ACME_RESTFUL_PUBLIC_KEY")
	if rawKey == "" {
		panic("no public key found. please set ACME_RESTFUL_PUBLIC_KEY")
	}

	pemBlock, _ := pem.Decode([]byte(strings.Join([]string{"-----BEGIN RSA PUBLIC KEY-----", rawKey, "-----END RSA PUBLIC KEY-----"}, "\n")))
	if pemBlock == nil {
		panic("failed to decode public key")
	}
	pKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}
	key = pKey.(*rsa.PublicKey)
}

func Middleware() gin.HandlerFunc {
	initBackend()
	return Protect
}

func Protect(c *gin.Context) {
	ok, rawToken := checkToken(c)
	if !ok {
		return
	}
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusForbidden, err)
		return
	}
	if !token.Valid {
		_ = c.AbortWithError(http.StatusForbidden, fmt.Errorf("token is not valid"))
		return
	}
	c.Set("user", token)
}

func checkToken(c *gin.Context) (bool, string) {
	auth := c.GetHeader("authorization")
	if auth == "" {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("auth header missing"))
		return false, ""
	}
	if !strings.HasPrefix(auth, oidc.PrefixBearer) {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid header"))
		return false, ""
	}
	return true, strings.TrimPrefix(auth, oidc.PrefixBearer)
}
