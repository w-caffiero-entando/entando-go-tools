package http

import (
    "bytes"
    "net/http"
    "time"
)

type H struct {
    Name  string
    Value string
}

type U struct {
    User string
    Pass string
}

type D struct {
    Data string
}

type V struct {
    Verb string
}

type Timeout struct {
    Timeout time.Duration
}

type BaseData struct {
    verb    string
    timeout time.Duration
    data    *bytes.Buffer
}

func parseBasicArgs(args []interface{}, fallbackData BaseData) BaseData {
    res := fallbackData
    for _, a := range args {
        switch a.(type) {
        case V:
            res.verb = a.(V).Verb
        case D:
            res.data = bytes.NewBufferString(a.(D).Data)
        case Timeout:
            res.timeout = a.(Timeout).Timeout
        }
    }
    return res
}

func applyHttpArgs(req *http.Request, args ...interface{}) {
    for _, a := range args {
        switch a.(type) {
        case H:
            t := a.(H)
            req.Header.Add(t.Name, t.Value)
        case U:
            t := a.(U)
            req.SetBasicAuth(t.User, t.Pass)
        }
    }
}
