package entapp

// Macro functions that implement frequent patterns

import (
    "time"

    "entando_go_tools/pkg/util/jnav"

    "entando_go_tools/pkg/util/http"
    "entando_go_tools/pkg/util/validators"
)

// Prepares for the execution of a request against entando
func prepareRemoteAction() (ApiAccessInfo, string) {
    // INGRESS INFO ON KUBE
    ingr := QueryIngressesInfo()
    validators.FATAL_IF(ingr.Main == "", "Unable to determine the main ingress url (s1)")
    var schema string
    //if env.LATEST_APP_SCHEME == "https" {
    //    schema = http.ProbeUrlScheme(ingr.KC, "https", "http", 7*time.Second, "http")
    //} else {
    schema = http.ProbeUrlScheme(ingr.KC, "http", "https", 7*time.Second, "http")
    //}
    // TOKEN ON KEYCLOAK
    aai := ObtainApiAccess(schema + "://" + ingr.KC)
    validators.FATAL_IF(schema == "", "Unable to determine the main ingress url (s2)")
    //config.SaveCfgValue("LATEST_APP_SCHEME", schema)
    return aai, schema + "://" + ingr.ECR
}

func (action ResourceAction) Authenticate() ResourceAction {
    aai, url := prepareRemoteAction()
    action.IngressUrl = url
    action.Token = aai.Token
    return action
}

func (action ResourceAction) Exec() (jnav.Obj, error) {
    return Exec(action)
}

func (action *ResourceAction) Merge(with ResourceAction) *ResourceAction {
    mergeSet := func(this *string, withThis interface{}, exceptFor interface{}) {
        if withThis != exceptFor {
            switch withThis.(type) {
            case string:
                *this = withThis.(string)
            case *string:
                *this = *withThis.(*string)
            }
        }
    }
    mergeSet(&action.Verb, with.Verb, "")
    mergeSet(&action.IngressUrl, with.IngressUrl, "")
    mergeSet(&action.ResourceName, with.ResourceName, "")
    mergeSet(&action.ResourceId, with.ResourceId, "")
    mergeSet(&action.SubResourceName, with.SubResourceName, "")
    mergeSet(&action.Payload, with.Payload, "")
    mergeSet(&action.Token, with.Token, "")

    return action
}
