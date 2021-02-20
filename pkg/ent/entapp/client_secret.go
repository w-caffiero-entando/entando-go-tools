package entapp

import (
    "encoding/base64"
    "fmt"
    "strings"

    "entando_go_tools/pkg/util/sys/spawn"

    "entando_go_tools/pkg/support/keycloak"

    "entando_go_tools/pkg/support/kube"
)

func QueryClientCredentials(clientSecretName string) (keycloak.ClientCredentials, error) {
    res, _ := kube.Kubectl(nil,
        "get",
        "secret",
        clientSecretName,
        spawn.SOA("o", "jsonpath", "{.data.clientId}:{.data.clientSecret}"),
    )
    arr := strings.Split(res.Stdout, ":")
    if len(arr) != 2 {
        return keycloak.ClientCredentials{}, fmt.Errorf("illegal format detected while querying the client secret")
    }
    clientId, _ := base64.StdEncoding.DecodeString(arr[0])
    clientSecret, _ := base64.StdEncoding.DecodeString(arr[1])

    return keycloak.ClientCredentials{ClientId: string(clientId), ClientSecret: string(clientSecret)}, nil
}
