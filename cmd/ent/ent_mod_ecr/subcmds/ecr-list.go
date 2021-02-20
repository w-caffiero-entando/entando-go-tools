package subcmds

import (
    "github.com/spf13/cobra"

    "entando_go_tools/pkg/ent/entapp"
    "entando_go_tools/pkg/util/jnav"
)

func init() {
    moduleCmd.AddCommand(&cobra.Command{
        Use:   "list",
        Short: "list the application profiles",
        Long:  `list the application profiles`,
        Run: func(cmd *cobra.Command, args []string) {
            resp, _ := entapp.ResourceAction{
                Verb:         "GET",
                ResourceName: "components",
            }.Authenticate().Exec()

            jnav.PrintTable(resp.Nav("payload").AsArr(),
                "code",
                "lastJob.status",
                "lastJob.componentVersion")
        },
    })
}
