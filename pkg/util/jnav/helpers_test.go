package jnav

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "entando_go_tools/pkg/util/sys"
)

func TestPrintTable(t *testing.T) {

    raw := []interface{}{
        map[string]interface{}{"a": 10, "b": 11, "c": 12},
        map[string]interface{}{"a": 20, "b": 21, "c": 22},
    }

    type args struct {
        arr        Arr
        fieldsPath []string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {
            name: "Simple",
            args: args{
                arr:        FromArr(raw),
                fieldsPath: []string{"a", "b"},
            },
            want: "| 10 | 11 |\n" + "| 20 | 21 |\n",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            output := sys.CaptureGoOutput(true, false, func() {
                PrintTable(tt.args.arr, tt.args.fieldsPath...)
            }).Stdout
            assert.Equal(t, tt.want, output)
        })
    }
}
