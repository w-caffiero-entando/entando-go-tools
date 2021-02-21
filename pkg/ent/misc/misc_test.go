package misc

import (
    "testing"
)

func TestNormalizeBaseEntandoUrl(t *testing.T) {
    type args struct {
        url string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {"Simple", args{"http://example.com/"}, "http://example.com/"},
        {"Simple", args{"http://example.com"}, "http://example.com/"},
        {"Simple", args{"http://example.com\n"}, "http://example.com/"},
        {"Simple", args{"http://example.com\r\n"}, "http://example.com/"},
        {"Simple", args{"http://example.com\n\r"}, "http://example.com/"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := NormalizeBasePath(tt.args.url); got != tt.want {
                t.Errorf("NormalizeBasePath() = %v, want %v", got, tt.want)
            }
        })
    }
}
