package entapp

import (
    "os"

    "github.com/spf13/cobra"
)

// BASE
var EntandoEnvHome = os.Getenv("ENTANDO_ENT_HOME")

// The execution environment
var Params ParamsType

func ParseAppParams(cmd *cobra.Command) {

    tmp := []*string{
        cmd.PersistentFlags().String("home", EntandoEnvHome, "kubeconfig file"),
        cmd.PersistentFlags().String("profile", "", "ent application profile"),
        cmd.PersistentFlags().String("appname", "", "entando app"),
        cmd.PersistentFlags().String("vm", "", "managed vm"),
    }
    tmpB := []*bool{
        cmd.PersistentFlags().Bool("yes", false, "managed vm"),
    }

    cobra.OnInitialize(func() {
        Params = ParamsType{
            TheHomeDir:      *tmp[0],
            TheProfile:      Profile(*tmp[1]),
            TheAppname:      AppName(*tmp[2]),
            TheVmName:       *tmp[3],
            AssumeYesForAll: *tmpB[0],
        }
    })
}

func ParseBaseTestExecutionParams() {
    Params = ParamsType{
        TheProfile:      Profile(os.Getenv("TEST_profile")),
        TheAppname:      AppName(os.Getenv("TEST_appname")),
        TheVmName:       os.Getenv("TEST_vm_name"),
        AssumeYesForAll: os.Getenv("TEST_yes_for_all") == "yes",
    }
}

type AppName string
type Profile string

type ParamsType struct {
    TheHomeDir      string
    TheProfile      Profile
    TheAppname      AppName
    TheVmName       string
    AssumeYesForAll bool
}
