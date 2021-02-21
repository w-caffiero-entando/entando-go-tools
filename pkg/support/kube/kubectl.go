package kube

import (
    "os/exec"
    "strings"

    "entando_go_tools/pkg/util/sys/spawn"

    "entando_go_tools/pkg/util/sugar"
)

func Kubectl(gocmd **exec.Cmd, args ...interface{}) (spawn.Res, error) {

    cmd := string(Params.KubectlCmd)
    if isSudoExplicitlyUsed(cmd) {
        // It will be eventually re-added by Spawn
        cmd = strings.Replace(cmd, "sudo ", "", 1)
        cmd = strings.Replace(cmd, "sudo\t", "", 1)
    }

    output, err := spawn.Spawn(gocmd,
        sugar.FNN(cmd, "kubectl"),             // The base kubectl command
        spawn.ComposeArgs(args, mkBaseArgs()), // The command line options
        mkEnv(),                               // The environment for the spawned process
        spawn.Options{
            WithSudo:      shouldSudoKubectl(), // The sudo execution flag
            CaptureStdout: true,
        },
    )
    return output, err
}

func mkEnv() spawn.Environ {
    ret := spawn.Environ{}

    if Params.Kubeconfig != DefaultKubeconfig {
        ret["KUBECONFIG"] = string(Params.Kubeconfig)
    }

    return ret
}
