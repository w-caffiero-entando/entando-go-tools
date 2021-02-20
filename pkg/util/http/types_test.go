package http

import (
    "bytes"
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func Test_parseBasicArgs(t *testing.T) {

    fallback := BaseData{"XXX", 999 * time.Second, bytes.NewBuffer([]byte("XXXXXXXXXXXXXXXXXX"))}

    assert.Equal(t,
        parseBasicArgs(
            []interface{}{
                V{Verb: "POST"},
                H{Name: "User-Agent", Value: "entando-ent/0.0.0"},
                H{Name: "Accept", Value: "application/json"},
                H{Name: "Content-Type", Value: "application/x-www-form-urlencoded"},
                U{User: "XXX", Pass: "YYY"},
                D{Data: "grant_type=client_credentials"},
                Timeout{Timeout: 15 * time.Second},
            }, fallback,
        ),
        BaseData{
            "POST",
            15 * time.Second,
            bytes.NewBuffer([]byte("grant_type=client_credentials")),
        },
    )
    assert.Equal(t,
        parseBasicArgs(
            []interface{}{
                H{Name: "User-Agent", Value: "entando-ent/0.0.0"},
                H{Name: "Accept", Value: "application/json"},
                H{Name: "Content-Type", Value: "application/x-www-form-urlencoded"},
                U{User: "XXX", Pass: "YYY"},
            }, fallback,
        ), fallback,
    )
}

func Test_applyHttpArgs(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://example.com", nil)

    applyHttpArgs(req,
        V{Verb: "SOMETHING-ELSE"},
        H{Name: "User-Agent", Value: "entando-ent/0.0.0"},
        H{Name: "Accept", Value: "application/json"},
        H{Name: "Content-Type", Value: "application/x-www-form-urlencoded"},
        U{User: "XXX", Pass: "YYY"},
        D{Data: "grant_type=client_credentials"},
        Timeout{Timeout: 15 * time.Second},
    )

    assert.Equal(t, req.Method, "GET")
    assert.Equal(t, req.Header.Get("Accept"), "application/json")
    assert.Equal(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
    assert.Equal(t, req.UserAgent(), "entando-ent/0.0.0")
    username, password, _ := req.BasicAuth()
    assert.Equal(t, username, "XXX")
    assert.Equal(t, password, "YYY")
}
