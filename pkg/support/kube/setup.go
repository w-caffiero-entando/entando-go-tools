package kube

import (
    "strings"

    "entando_go_tools/pkg/util/sugar"

    "entando_go_tools/pkg/util/validators"

    "entando_go_tools/pkg/util/sys"
)

var isCommandAvailable = sys.IsCommandAvailable

type KubeVmInfo struct {
    name      string
    appname   string
    namespace string
}

type KubeStatusType struct {
    command      string
    kubeconfig   string
    kubectlMode  string
    designatedVm KubeVmInfo
}

func shouldSudoKubectl() bool {

    autoSudoDisabled := wasAutoSudoDisabled()
    sudoInCmd := isSudoExplicitlyUsed(string(Params.KubectlCmd))
    sudoAvailable := isCommandAvailable("sudo")

    if ! sudoAvailable {
        return false
    }
    if sudoInCmd {
        return true
    }
    if autoSudoDisabled {
        return false
    }

    // ELSE..

    // ..no sudo if custom command or customized connection
    if sugar.FNN(
        string(Params.KubectlCmd),
        string(Params.Kubectx),
        string(Params.Kubeconfig)) != "" {
        return false
    }

    // ..no sudo on windows
    if sys.OsInfo.IsWin {
        return false
    }

    // ..otherwise automatically sudo the kubectl command
    return true
}

func wasAutoSudoDisabled() bool {
    return !Params.AutoSudo
}

func isSudoExplicitlyUsed(cmd string) bool {
    return strings.HasPrefix(cmd, "sudo ") || strings.HasPrefix(cmd, "sudo\t")
}

func mkBaseArgs() []interface{} {
    // NAMESPACE
    ns := Params.Namespace

    var args []interface{}

    if ns == AllNamespaces {
        args = append(args, "--all-namespaces")
    } else
    if ns != "" {
        args = append(args, "--namespace")
        args = append(args, string(ns))
    } else {
        validators.FATAL("Detected null namespace")
    }

    // CONTEXT
    if Params.Kubectx != "" {
        args = append(args, "--context")
        args = append(args, string(Params.Kubectx))
    }

    return args
}
