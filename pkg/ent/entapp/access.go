package entapp

import (
    "entando_go_tools/pkg/support/keycloak"

    "entando_go_tools/pkg/util/validators"
)

type ApiAccessInfo struct {
    Url   string
    Token string
}

func ObtainApiAccess(apiUrl string) ApiAccessInfo {
    // Client Credentials
    creds, _ := QueryClientCredentials(string(Params.TheAppname) + "-server-secret")

    tokenEndpoint, _ := keycloak.QueryTokenEndpoint(apiUrl)

    validators.FATAL_IF(tokenEndpoint == "", "Unable to determine authorization endpoint")

    //log.DEF.Warn(0, "> %s", tokenEndpoint)
    var token, _ = keycloak.Authenticate(tokenEndpoint, creds)

    validators.FATAL_IF(token == "", "Unable to obtain the access token")

    return ApiAccessInfo{apiUrl, token}
}
