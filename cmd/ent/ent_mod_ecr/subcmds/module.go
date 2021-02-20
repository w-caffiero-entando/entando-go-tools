package subcmds

import (
    "log"
    "os"

    "github.com/spf13/cobra"

    . "entando_go_tools/cmd/ent/all"

    "entando_go_tools/pkg/ent/entapp"
    "entando_go_tools/pkg/support/kube"
)

var moduleCmd = &cobra.Command{
    Use:   "ecr",
    Short: "Helps dealing with the entando component repository (ECR)",
    Long:  `Helps dealing with the entando component repository (ECR)`,
}

func Execute() {
    if err := RootCmd.Execute(); err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
}

func init() {
    kube.ParseKubeParams(RootCmd)
    entapp.ParseAppParams(RootCmd)
    RootCmd.AddCommand(moduleCmd)
}
