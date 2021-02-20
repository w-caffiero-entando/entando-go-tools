package http

import "testing"

func TestConcatURL(t *testing.T) {
    type args struct {
        urlA string
        urlB string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {"Simple", args{"x/", "y"}, "x/y"},
        {"SimpleAlt", args{"x", "/y"}, "x/y"},
        {"TooSimple", args{"x/", ""}, "x/"},
        {"TooSimpleAlt", args{"x", ""}, "x"},
        {"Missing", args{"x", "y"}, "x/y"},
        {"Fusion", args{"/", "/"}, "/"},
        {"Collapse", args{"x/", "/y"}, "x/y"},
        {"CollapseAlt", args{"/", "/y"}, "/y"},
        {"PreserveAbsolute", args{"", "/y"}, "/y"},
        {"PreserveRelative", args{"", "y"}, "y"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := ConcatURL(tt.args.urlA, tt.args.urlB); got != tt.want {
                t.Errorf("ConcatURL() = %v, want %v", got, tt.want)
            }
        })
    }
}
