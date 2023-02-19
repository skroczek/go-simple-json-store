package oicd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/oidc/pkg/oidc"
	"net/http"
	"strings"
)

type Oicd struct {
	provider rs.ResourceServer
}

func NewOicd(provider rs.ResourceServer) Oicd {
	return Oicd{provider: provider}
}

func (o *Oicd) Middleware() gin.HandlerFunc {
	return o.Protect
}

func (o *Oicd) Protect(c *gin.Context) {
	ok, token := o.checkToken(c)
	if !ok {
		return
	}
	resp, err := rs.Introspect(c, o.provider, token)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}
	if !resp.IsActive() {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("token is not active"))
		return
	}
	c.Set("user", resp)
}
func (o *Oicd) checkToken(c *gin.Context) (bool, string) {
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
