package kube

import (
    "reflect"
    "testing"

    "entando_go_tools/pkg/util/sys"

    "github.com/stretchr/testify/assert"
)

var OldIsCommandAvailable = isCommandAvailable

func Test_shouldSudoKubectl(t *testing.T) {
    OldParams := Params
    Params.AutoSudo = true
    assert.True(t, shouldSudoKubectl())
    OldParams = Params
    Params.AutoSudo = true
    isCommandAvailable = func(string) bool { return false }
    assert.False(t, shouldSudoKubectl())
    isCommandAvailable = OldIsCommandAvailable
    Params = OldParams
    Params.AutoSudo = true
    sys.OsInfo.IsWin = true
    assert.False(t, shouldSudoKubectl())
    sys.OsInfo.IsWin = false
    Params = OldParams
    Params.KubectlCmd = KKubectlCommand("kubectl")
    Params.AutoSudo = true
    assert.False(t, shouldSudoKubectl())
    Params = OldParams
    Params.KubectlCmd = KKubectlCommand("sudo kubectl")
    Params.AutoSudo = false
    assert.True(t, shouldSudoKubectl())
    Params.KubectlCmd = KKubectlCommand("kubectl")
    Params.AutoSudo = false
    assert.False(t, shouldSudoKubectl())
    Params = OldParams
    Params.Kubeconfig = KConfig("x")
    Params.KubectlCmd = KKubectlCommand("kubectl")
    Params.AutoSudo = true
    assert.False(t, shouldSudoKubectl())
}

func Test_mkBaseArgs(t *testing.T) {
    tests := []struct {
        name string
        ns   string
        ctx  string
        want []interface{}
    }{
        {"Simple", "xxx", "", []interface{}{"--namespace", "xxx"}},
        {"Simple", "*", "", []interface{}{"--all-namespaces"}},
        {"Simple", "*", "yyy", []interface{}{"--all-namespaces", "--context", "yyy"}},
        {"Simple", "xxx", "yyy", []interface{}{"--namespace", "xxx", "--context", "yyy"}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            OldParams := Params
            Params.Namespace = KNamespace(tt.ns)
            Params.Kubectx = KContext(tt.ctx)
            if got := mkBaseArgs(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("mkBaseArgs() = %v, want %v", got, tt.want)
            }
            Params = OldParams
        })
    }
    Params.Namespace = KNamespace("")
    Params.Kubectx = KContext("yyy")
    assert.Panics(t, func() {
        mkBaseArgs()
    })
}
