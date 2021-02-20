package subcmds

import (
    "fmt"
    "time"

    "github.com/spf13/cobra"

    . "entando_go_tools/pkg/util/sugar"

    "entando_go_tools/pkg/ent/entapp"
)

func init() {
    watchCmd := &cobra.Command{
        Use: "watch",
    }
    moduleCmd.AddCommand(watchCmd)
    watchCmd.AddCommand(&cobra.Command{
        Use:   "install",
        Short: "watch the status of a bundle install",
        Long:  `watch the status of a bundle install`,
        Run: func(cmd *cobra.Command, args []string) {
            ecrWatchAction(
                SweetArr(args).At(0).StringOrFail("please provide the bundle id"),
                "install",
                "", "")
        },
    })
    watchCmd.AddCommand(&cobra.Command{
        Use:   "uninstall",
        Short: "watch the status of a bundle uninstall",
        Long:  `watch the status of a bundle uninstall`,
        Run: func(cmd *cobra.Command, args []string) {
            ecrWatchAction(
                SweetArr(args).At(0).StringOrFail("please provide the bundle id"),
                "uninstall",
                "", "")
        },
    })
}

func ecrWatchAction(bundleName string, mode string, ingressUrl string, token string) {
    start := time.Now()
    terminated := false
    for !terminated {
        action := entapp.ResourceAction{
            Verb:            "GET",
            ResourceName:    "components",
            ResourceId:      bundleName,
            SubResourceName: mode,
        }
        if ingressUrl != "" {
            action.IngressUrl = ingressUrl
            action.Token = token
        } else {
            action = action.Authenticate()
        }

        resp, err := action.Exec()

        if err != nil {
            FATAL(err.Error())
        }

        rawStatus, _ := resp.Nav("payload.status").AsString()

        switch rawStatus {
        case "INSTALL_IN_PROGRESS", "INSTALL_CREATED", "UNINSTALL_IN_PROGRESS", "UNINSTALL_CREATED":
        case "INSTALL_COMPLETED":
            terminated = true
        case "UNINSTALL_COMPLETED":
            terminated = true
        case "INSTALL_ROLLBACK", "INSTALL_ERROR":
            rawStatus = "\nERROR, ROLLING BACK."
            terminated = true
        default:
            rawStatus = "Unknown status: \"" + rawStatus + "\""
            terminated = true
        }

        fmt.Printf("%03d | %s                     \r",
            int(time.Since(start)/time.Second),
            rawStatus)

        time.Sleep(3)
    }
}
