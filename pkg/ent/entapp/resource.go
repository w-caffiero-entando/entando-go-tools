package entapp

import (
    "time"

    "entando_go_tools/pkg/util/jnav"

    "entando_go_tools/pkg/util/http"
)

type ResourceAction struct {
    Verb            string
    IngressUrl      string
    ResourceName    string
    ResourceId      string
    SubResourceName string
    Payload         string
    Token           string
}

func Exec(action ResourceAction) (jnav.Obj, error) {
    //-
    var url = action.IngressUrl

    if action.ResourceName != "" {
        url = http.ConcatURL(url, action.ResourceName)
    }
    if action.ResourceId != "" {
        url = http.ConcatURL(url, action.ResourceId)
    }
    if action.SubResourceName != "" {
        url = http.ConcatURL(url, action.SubResourceName)
    }

    jnavRes, err := http.RequestJson(
        url,
        http.V{Verb: action.Verb},
        http.H{Name: "User-Agent", Value: "entando-ent/0.0.0"},
        http.H{Name: "Content-Type", Value: "application/json"},
        http.H{Name: "Accept", Value: "*/*"},
        http.H{Name: "Authorization", Value: "Bearer " + action.Token},
        http.H{Name: "Origin", Value: action.IngressUrl},
        http.Timeout{Timeout: 15 * time.Second},
    )

    if err != nil {
        return jnav.Obj{}, err
    }

    return jnavRes, nil
}
