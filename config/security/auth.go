package security

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/env"
)

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func validateTokenAndParseAuthorities(token string) (string, int) {
	client := gocloak.NewClient(env.KeycloakHost())

	ctx := context.Background()
	_, claims, err := client.DecodeAccessToken(ctx, token, "fdk")

	authorities := ""
	errStatus := http.StatusOK

	if err != nil {
		errStatus = http.StatusUnauthorized
	} else if claims == nil {
		errStatus = http.StatusForbidden
	} else {
		validError := claims.Valid()
		if validError != nil {
			errStatus = http.StatusForbidden
		}

		if !claims.VerifyAudience(env.SecurityValues.TokenAudience, true) {
			errStatus = http.StatusForbidden
		}

		authClaim := (*claims)["authorities"]
		if authClaim != nil {
			authorities = authClaim.(string)
		}
	}

	return authorities, errStatus
}

func hasOrganizationWriteRole(authorities string, org string) bool {
	sysAdminAuth := env.SecurityValues.SysAdminAuth
	orgAdminAuth := fmt.Sprintf("%s:%s:%s", env.SecurityValues.OrgType, org, env.SecurityValues.AdminPermission)
	orgWriteAuth := fmt.Sprintf("%s:%s:%s", env.SecurityValues.OrgType, org, env.SecurityValues.WritePermission)
	if strings.Contains(authorities, sysAdminAuth) {
		return true
	} else if strings.Contains(authorities, orgAdminAuth) {
		return true
	} else if strings.Contains(authorities, orgWriteAuth) {
		return true
	} else {
		return false
	}
}

func hasOrganizationReadRole(authorities string, org string) bool {
	orgReadAuth := fmt.Sprintf("%s:%s:%s", env.SecurityValues.OrgType, org, env.SecurityValues.ReadPermission)
	if hasOrganizationWriteRole(authorities, org) {
		return true
	} else if strings.Contains(authorities, orgReadAuth) {
		return true
	} else {
		return false
	}
}

func RequireWriteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorities, status := validateTokenAndParseAuthorities(c.GetHeader("Authorization"))

		if status != http.StatusOK {
			respondWithError(c, status, http.StatusText(status))
		} else if !hasOrganizationWriteRole(authorities, c.Param("catalogId")) {
			respondWithError(c, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		}

		c.Next()
	}
}

func RequireReadAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorities, status := validateTokenAndParseAuthorities(c.GetHeader("Authorization"))

		if status != http.StatusOK {
			respondWithError(c, status, http.StatusText(status))
		} else if !hasOrganizationReadRole(authorities, c.Param("catalogId")) {
			respondWithError(c, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		}

		c.Next()
	}
}
