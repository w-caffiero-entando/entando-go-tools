package subcmds

import (
    "github.com/spf13/cobra"

    "entando_go_tools/pkg/util/sugar"

    "entando_go_tools/pkg/util/jnav"

    "entando_go_tools/pkg/ent/entapp"
)

func init() {
    moduleCmd.AddCommand(&cobra.Command{
        Use:   "usage-status",
        Short: "usage status of the components of a bundle",
        Long:  `usage status of the components of a bundle`,
        Run: func(cmd *cobra.Command, args []string) {
            var bundleName = sugar.SweetArr(args).At(0).StringOrFail("please provide the bundle id")

            resp, _ := entapp.ResourceAction{
                Verb:            "GET",
                ResourceName:    "components",
                ResourceId:      bundleName,
                SubResourceName: "usage",
            }.Authenticate().Exec()

            jnav.PrintTable(resp.Nav("payload").AsArr(), "type", "code", "usage")
        },
    })
}
