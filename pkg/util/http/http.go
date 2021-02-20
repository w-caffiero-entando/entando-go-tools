package http

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "entando_go_tools/pkg/util/jnav"
)

func ProbeUrlScheme(baseUrl, schema1, schema2 string, timeout time.Duration, fallback string) string {
    if TestUrl(schema1+"://"+baseUrl, timeout) {
        return schema1
    }
    if TestUrl(schema2+"://"+baseUrl, timeout) {
        return schema2
    }
    return fallback
}

func TestUrl(url string, timeout time.Duration) bool {
    tr := &http.Transport{
        IdleConnTimeout:    timeout,
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Get(url)
    return err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300
}

func Request(url string, args ...interface{}) (string, error) {

    basic := parseBasicArgs(args, BaseData{"GET", 15 * time.Second, nil})

    tr := &http.Transport{
        IdleConnTimeout:    basic.timeout,
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    var req *http.Request
    var err error
    if basic.data == nil {
        req, err = http.NewRequest(basic.verb, url, nil)
    } else {
        req, err = http.NewRequest(basic.verb, url, basic.data)
    }

    applyHttpArgs(req, args...)

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer func() { _ = resp.Body.Close() }()

    if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
        body, err := ioutil.ReadAll(resp.Body)
        return string(body), err
    } else {
        return "", fmt.Errorf("http status error %d", resp.StatusCode)
    }
}

func RequestJson(url string, args ...interface{}) (jnav.Obj, error) {
    res, err := Request(url, args...)
    if err != nil {
        return jnav.Obj{}, err
    }
    var mapRes map[string]interface{}
    err = json.Unmarshal([]byte(res), &mapRes)
    if err != nil {
        return jnav.Obj{}, err
    }
    return jnav.FromMap(mapRes), nil
}
