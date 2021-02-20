package sugar

import (
    "testing"
)

func TestFNN(t *testing.T) {
    type args struct {
        args []string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {"Empty", args{[]string{}}, ""},
        {"None", args{[]string{""}}, ""},
        {"Single", args{[]string{"this"}}, "this"},
        {"First 1", args{[]string{"this", ""}}, "this"},
        {"First 2", args{[]string{"this", "Second"}}, "this"},
        {"First 3", args{[]string{"this", "", "Second"}}, "this"},
        {"Second", args{[]string{"", "this", ""}}, "this"},
        {"Last", args{[]string{"", "this", ""}}, "this"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := FNN(tt.args.args...); got != tt.want {
                t.Errorf("FNN() = %v, want %v", got, tt.want)
            }
        })
    }
}
