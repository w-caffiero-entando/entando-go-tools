package entapp

import (
    "entando_go_tools/pkg/util/sys/spawn"

    "entando_go_tools/pkg/ent/misc"

    "entando_go_tools/pkg/support/kube"
)

type IngressesInfo struct {
    Main string
    ECR  string
    KC   string
}

func QueryIngressesInfo() IngressesInfo {
    jsonpath := `{range .items[?(@.metadata.labels.EntandoApp)]}`                     // selector
    jsonpath += `{.spec.rules[0].host}{.spec.rules[0].http.paths[2].path}{"\n"}{end}` // host+path

    res, _ := kube.Kubectl(nil,
        "get",
        "ingress",
        spawn.SOA("o", "custom-columns", "NAME:.metadata.name,HOST:.spec.rules[0].host"),
    )

    var urlMain, urlECI, urlKC string

    for _, ingr := range misc.FieldsToMap(res.Stdout) {
        switch ingr["NAME"] {
        case string(Params.TheAppname) + "-ingress":
            urlMain = misc.NormalizeBasePath(ingr["HOST"]) + "app-builder/"
            urlECI = misc.NormalizeBasePath(ingr["HOST"]) + "digital-exchange/"
        //case string(Params.TheAppname) + "-eci-ingress":
        //    urlECI = util.NormalizeBasePath(ingr["HOST"]) + "digital-exchange/"
        case string(Params.TheAppname) + "-kc-ingress":
            urlKC = misc.NormalizeBasePath(ingr["HOST"]) + "auth/"
        }
    }

    return IngressesInfo{urlMain, urlECI, urlKC}
}
