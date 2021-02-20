package subcmds

import (
    "fmt"

    "github.com/spf13/cobra"

    . "entando_go_tools/pkg/util/sugar"

    "entando_go_tools/pkg/ent/entapp"
)

func init() {
    moduleCmd.AddCommand(&cobra.Command{
        Use:   "uninstall",
        Short: "uninstalls a bundle",
        Long:  `uninstalls a bundle`,
        Run: func(cmd *cobra.Command, args []string) {
            var bundleName = SweetArr(args).At(0).StringOrFail("please provide the bundle id")

            action := entapp.ResourceAction{
                Verb:            "POST",
                ResourceName:    "components",
                ResourceId:      bundleName,
                SubResourceName: "uninstall",
            }

            resp, _ := action.Authenticate().Exec()

            ecrWatchAction(action.ResourceId, "uninstall", action.IngressUrl, action.Token)

            fmt.Printf("%v", resp)
        },
    })
}
