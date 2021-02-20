package all

import (
    "github.com/spf13/cobra"

    "entando_go_tools/pkg/ent/profile"
)

var RootCmd = &cobra.Command{
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        profile.PrintCurrentProfileInfo(false)
    },
    Use:   "ent",
    Short: "The entando cli",
    Long:  `The entando cli`,
}
