package http

import "strings"

func ConcatURL(urlA, urlB string) string {
    if urlB == "" {
        return urlA
    }
    if urlA != "" && ! strings.HasSuffix(urlA, "/") && ! strings.HasPrefix(urlB, "/") {
        urlA += "/"
    }
    if strings.HasSuffix(urlA, "/") && strings.HasPrefix(urlB, "/") {
        urlA = strings.Trim(urlA, "/")
    }
    urlA += urlB
    return urlA
}
