package keycloak

import (
    "time"

    "entando_go_tools/pkg/util/jnav"

    "entando_go_tools/pkg/util/http"
)

// Obtains the token endpoint from keycloak
func QueryTokenEndpoint(authUrl string) (string, error) {
    url := authUrl + "realms/entando/.well-known/openid-configuration"
    jnavRes, err := http.RequestJson(url, 15*time.Second)
    return jnav.MkObjValueReturn(jnavRes, "token_endpoint", err)
}

// Authenticates and returns the related access token
func Authenticate(tokenEndpoint string, creds ClientCredentials) (string, error) {
    mapRes, err := http.RequestJson(
        tokenEndpoint,
        http.V{Verb: "POST"},
        http.H{Name: "User-Agent", Value: "entando-ent/0.0.0"},
        http.H{Name: "Accept", Value: "application/json"},
        http.H{Name: "Content-Type", Value: "application/x-www-form-urlencoded"},
        http.U{User: creds.ClientId, Pass: creds.ClientSecret},
        http.D{Data: "grant_type=client_credentials"},
        http.Timeout{Timeout: 15 * time.Second},
    )
    return jnav.MkObjValueReturn(mapRes, "access_token", err)
}

type ClientCredentials struct {
    ClientId     string
    ClientSecret string
}
