package subcmds

import (
    "net/url"

    "github.com/spf13/cobra"

    "entando_go_tools/pkg/util/validators"

    "entando_go_tools/pkg/util/jnav"

    . "entando_go_tools/pkg/util/sugar"

    . "entando_go_tools/pkg/util/log/def"

    "entando_go_tools/pkg/ent/entapp"
)

func init() {
    cmd := &cobra.Command{
        Use:   "install",
        Short: "installs a registered bundle",
        Long:  `installs a registered bundle`,
        Run: func(cmd *cobra.Command, args []string) {
            validators.Assert_Ext_Ic_Id("Version", InstallParams.Version, "", true, "fatal")
            validators.Assert_Ext_Ic_Id("Conflict strategy", InstallParams.ConflictStrategy, "", true, "fatal")

            var bundleName = SweetArr(args).At(0).StringOrFail("please provide the bundle id")

            postData := &url.Values{}
            setUrlValueNN(postData, "version", InstallParams.Version)
            setUrlValueNN(postData, "conflictStrategy", InstallParams.ConflictStrategy)

            BizLog.Info("Installation of bundle %v started", bundleName)

            action := entapp.ResourceAction{
                Verb:            "POST",
                ResourceName:    "components",
                ResourceId:      bundleName,
                SubResourceName: "install",
                Payload:         postData.Encode(),
            }
            resp, _ := action.Authenticate().Exec()

            jnav.PrintTable(resp.Nav("payload").AsObj(),
                "componentId",
                "componentVersion",
                "progress",
                "startedAt",
                "status")

            ecrWatchAction(action.ResourceId, "install", action.IngressUrl, action.Token)
        },
    }
    ParseParams(cmd)
    moduleCmd.AddCommand(cmd)
}

func ParseParams(cmd *cobra.Command) {
    cmd.PersistentFlags().StringVarP(&InstallParams.Version, "version", "v", "", "version to install")
    cmd.PersistentFlags().StringVar(&InstallParams.ConflictStrategy, "conflict-strategy", "", "component conflict strategy")
    cmd.PersistentFlags().BoolVar(&InstallParams.Watch, "watch", true, "watches the operation until completion")
}

type InstallParamsType struct {
    Version          string
    ConflictStrategy string
    Watch            bool
}

var InstallParams InstallParamsType

func setUrlValueNN(urlValues *url.Values, k string, v string) {
    if v != "" {
        urlValues.Set(k, v)
    }
}
