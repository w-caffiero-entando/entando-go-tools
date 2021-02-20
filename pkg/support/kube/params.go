package kube

import (
    "os"

    "github.com/spf13/cobra"

    V "entando_go_tools/pkg/util/validators"
)

var Params ParamsType

func ParseKubeParams(cmd *cobra.Command) {

    cmd.PersistentFlags().StringVar((*string)(&Params.Kubeconfig),
        "kubeconfig", string(DefaultKubeconfig), "kubeconfig file")
    cmd.PersistentFlags().StringVar((*string)(&Params.Kubectx),
        "kubectx", string(DefaultKubectx), "context name")
    cmd.PersistentFlags().StringVar((*string)(&Params.KubectlCmd),
        "kubectl", string(DefaultCommand), "kube management command (eg: \"kubectl\", \"oc\", \"k3s kubectl\" ..")
    cmd.PersistentFlags().BoolVar((*bool)(&Params.AutoSudo),
        "kubectl-autosudo", true, "automatic sudo detection for kubectl")

    ParseNamespaceParam(cmd)
}

func ParseNamespaceParam(cmd *cobra.Command) {
    var tmp1 *string
    var tmp2 *bool
    tmp1 = cmd.PersistentFlags().StringP("namespace", "n", "", "uses the given namespace")
    tmp2 = cmd.PersistentFlags().BoolP("all-namespaces", "A", false, "all namespaces")

    cobra.OnInitialize(func() {
        if tmp1 != nil {
            V.Assert_Ext_Ic_Id("namespace", *tmp1, "", false, "fatal")
            Params.Namespace = KNamespace(*tmp1)
        } else
        if tmp2 != nil && *tmp2 == true {
            Params.Namespace = AllNamespaces
        }
    })
}

func ParseBaseTestExecutionParams() {
    Params = ParamsType{
        KubectlCmd: KKubectlCommand(os.Getenv("TEST_kubectl_cmd")),
        Kubeconfig: KConfig(os.Getenv("TEST_kubeconfig")),
        Kubectx:    KContext(os.Getenv("TEST_kubectx")),
        Namespace:  KNamespace(os.Getenv("TEST_namespace")),
        AutoSudo:   os.Getenv("TEST_autosudo") == "true",
    }
}

type ParamsType struct {
    Namespace  KNamespace
    Kubeconfig KConfig
    Kubectx    KContext
    KubectlCmd KKubectlCommand
    AutoSudo   bool
}
