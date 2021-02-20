package profile

import (
    "fmt"

    . "entando_go_tools/pkg/util/log/def"

    "entando_go_tools/pkg/util/sugar"

    "entando_go_tools/pkg/support/kube"

    "entando_go_tools/pkg/ent/entapp"
)

func PrintCurrentProfileInfo(verbose bool) {
    if entapp.Params.TheProfile != "" {
        if verbose {
            BizLog.Info("Current application profile:")
            fmt.Printf(" - PROFILE NAME:  ${%s}\n", string(entapp.Params.TheProfile))
        } else {
            BizLog.Info("Current application profile: ${THIS_APP_PROFILE}")
        }
    } else {
        BizLog.Info("Currently using no application profile")
    }

    if verbose {
        fmt.Printf(" - APPNAME:       %s", sugar.FNN(string(entapp.Params.TheAppname), "{NONE}"))
        fmt.Printf(" - NAMESPACE:     %s", sugar.FNN(string(kube.Params.Namespace), "{NONE}"))
        fmt.Printf(" - K8S CONTEXT:   %s", sugar.FNN(string(kube.Params.Kubectx), "{NONE}"))
    }
}
